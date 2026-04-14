package zdb

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const (
	DefaultConfigKey              = "database"
	DefaultServiceName            = "db"
	DefaultManagerServiceName     = "gorm"
	defaultPingTimeoutSeconds     = 3
	defaultSlowThresholdMillis    = 200
	defaultMaxIdleConns           = 10
	defaultMaxOpenConns           = 100
	defaultConnMaxLifetimeSeconds = 3600
	defaultConnMaxIdleTimeSeconds = 600
	defaultGORMLogLevel           = "warn"
)

// Settings describes the GORM component configuration.
//
// Recommended config shape:
//
//	[database]
//	default = "main"
//	prepare_stmt = true
//	skip_default_transaction = true
//
//	[database.logger]
//	enabled = true
//	level = "warn"
//	slow_threshold_millis = 200
//
//	[[database.connections]]
//	name = "main"
//	driver = "mysql"
//	host = "127.0.0.1"
//	port = 3306
//	username = "app"
//	password = "secret"
//	database = "app"
//	params = "charset=utf8mb4&parseTime=True&loc=Local"
type Settings struct {
	Default                  string             `mapstructure:"default"`
	PrepareStmt              *bool              `mapstructure:"prepare_stmt"`
	SkipDefaultTransaction   *bool              `mapstructure:"skip_default_transaction"`
	DisableNestedTransaction *bool              `mapstructure:"disable_nested_transaction"`
	PingTimeout              int                `mapstructure:"ping_timeout"`
	Logger                   LogConfig          `mapstructure:"logger"`
	Connections              []ConnectionConfig `mapstructure:"connections"`
}

// LogConfig controls the GORM SQL logger.
type LogConfig struct {
	Enabled                   *bool  `mapstructure:"enabled"`
	Level                     string `mapstructure:"level"`
	SlowThresholdMillis       int    `mapstructure:"slow_threshold_millis"`
	IgnoreRecordNotFoundError *bool  `mapstructure:"ignore_record_not_found_error"`
}

// ConnectionConfig describes a single database connection.
//
// Legacy antgo-style flat fields remain supported:
//
//	log = true
//	level = 4
//	max_idle_conns = 50
//	max_open_conns = 200
//	conn_max_lifetime = 1800
//	conn_max_idle_time = 300
type ConnectionConfig struct {
	Name       string `mapstructure:"name"`
	Driver     string `mapstructure:"driver"`
	Type       string `mapstructure:"type"`
	DSN        string `mapstructure:"dsn"`
	Host       string `mapstructure:"host"`
	Hostname   string `mapstructure:"hostname"`
	Port       int    `mapstructure:"port"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Database   string `mapstructure:"database"`
	Parameters string `mapstructure:"params"`

	PrepareStmt              *bool `mapstructure:"prepare_stmt"`
	SkipDefaultTransaction   *bool `mapstructure:"skip_default_transaction"`
	DisableNestedTransaction *bool `mapstructure:"disable_nested_transaction"`
	PingTimeout              int   `mapstructure:"ping_timeout"`

	LogEnabled                *bool  `mapstructure:"log"`
	LogLevel                  string `mapstructure:"level"`
	SlowThresholdMillis       int    `mapstructure:"slow_threshold_millis"`
	IgnoreRecordNotFoundError *bool  `mapstructure:"ignore_record_not_found_error"`

	MaxIdleConns    int `mapstructure:"max_idle_conns"`
	MaxOpenConns    int `mapstructure:"max_open_conns"`
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime int `mapstructure:"conn_max_idle_time"`

	sqlDB                             *sql.DB
	connPool                          gorm.ConnPool
	effectivePrepareStmt              bool
	effectiveSkipDefaultTransaction   bool
	effectiveDisableNestedTransaction bool
	effectivePingTimeout              time.Duration
	effectiveLogger                   resolvedLogConfig
}

type resolvedLogConfig struct {
	enabled                   bool
	level                     gormlogger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

type legacyRootSettings struct {
	Default                  string             `mapstructure:"default"`
	PrepareStmt              *bool              `mapstructure:"prepare_stmt"`
	SkipDefaultTransaction   *bool              `mapstructure:"skip_default_transaction"`
	DisableNestedTransaction *bool              `mapstructure:"disable_nested_transaction"`
	PingTimeout              int                `mapstructure:"ping_timeout"`
	Connections              []ConnectionConfig `mapstructure:"connections"`
}

// DefaultSettings returns production-safe defaults.
func DefaultSettings() Settings {
	return Settings{
		PrepareStmt:              boolPtr(true),
		SkipDefaultTransaction:   boolPtr(true),
		DisableNestedTransaction: boolPtr(false),
		PingTimeout:              defaultPingTimeoutSeconds,
		Logger: LogConfig{
			Enabled:                   boolPtr(true),
			Level:                     defaultGORMLogLevel,
			SlowThresholdMillis:       defaultSlowThresholdMillis,
			IgnoreRecordNotFoundError: boolPtr(true),
		},
	}
}

func loadSettings(cfg ConfigProvider, key string) (Settings, error) {
	settings := DefaultSettings()
	if cfg == nil {
		return settings, nil
	}

	switch {
	case strings.TrimSpace(key) == "":
		if err := cfg.Unmarshal("", &settings); err != nil {
			return Settings{}, fmt.Errorf("unmarshal gorm config: %w", err)
		}
	case cfg.IsSet(key):
		if err := cfg.Unmarshal(key, &settings); err != nil {
			return Settings{}, fmt.Errorf("unmarshal gorm config %q: %w", key, err)
		}
	case cfg.IsSet("connections"):
		legacy := legacyRootSettings{
			PrepareStmt:              settings.PrepareStmt,
			SkipDefaultTransaction:   settings.SkipDefaultTransaction,
			DisableNestedTransaction: settings.DisableNestedTransaction,
			PingTimeout:              settings.PingTimeout,
		}
		if err := cfg.Unmarshal("", &legacy); err != nil {
			return Settings{}, fmt.Errorf("unmarshal legacy gorm config: %w", err)
		}
		settings.Default = legacy.Default
		settings.PrepareStmt = legacy.PrepareStmt
		settings.SkipDefaultTransaction = legacy.SkipDefaultTransaction
		settings.DisableNestedTransaction = legacy.DisableNestedTransaction
		settings.PingTimeout = legacy.PingTimeout
		settings.Connections = legacy.Connections
	}

	return settings.normalize()
}

func (s Settings) normalize() (Settings, error) {
	if s.PrepareStmt == nil {
		s.PrepareStmt = boolPtr(true)
	}
	if s.SkipDefaultTransaction == nil {
		s.SkipDefaultTransaction = boolPtr(true)
	}
	if s.DisableNestedTransaction == nil {
		s.DisableNestedTransaction = boolPtr(false)
	}
	if s.PingTimeout <= 0 {
		s.PingTimeout = defaultPingTimeoutSeconds
	}
	if s.Logger.Enabled == nil {
		s.Logger.Enabled = boolPtr(true)
	}
	if strings.TrimSpace(s.Logger.Level) == "" {
		s.Logger.Level = defaultGORMLogLevel
	}
	if s.Logger.SlowThresholdMillis <= 0 {
		s.Logger.SlowThresholdMillis = defaultSlowThresholdMillis
	}
	if s.Logger.IgnoreRecordNotFoundError == nil {
		s.Logger.IgnoreRecordNotFoundError = boolPtr(true)
	}
	if len(s.Connections) == 0 {
		return s, nil
	}

	seen := make(map[string]struct{}, len(s.Connections))
	for i := range s.Connections {
		conn, err := s.Connections[i].normalize(s)
		if err != nil {
			return Settings{}, err
		}
		if _, exists := seen[conn.Name]; exists {
			return Settings{}, fmt.Errorf("gorm: duplicated connection name %q", conn.Name)
		}
		seen[conn.Name] = struct{}{}
		s.Connections[i] = conn
	}

	if strings.TrimSpace(s.Default) == "" {
		if _, ok := seen["default"]; ok {
			s.Default = "default"
		} else {
			s.Default = s.Connections[0].Name
		}
	}
	if _, ok := seen[s.Default]; !ok {
		return Settings{}, fmt.Errorf("gorm: default connection %q not found", s.Default)
	}

	return s, nil
}

func (c ConnectionConfig) normalize(settings Settings) (ConnectionConfig, error) {
	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return ConnectionConfig{}, errors.New("gorm: connection name is required")
	}

	driver := firstNonEmpty(c.Driver, c.Type)
	c.Driver = normalizeDriver(driver)
	if c.Driver == "" {
		return ConnectionConfig{}, fmt.Errorf("gorm: connection %q has unsupported driver %q", c.Name, driver)
	}

	if c.Host == "" {
		c.Host = c.Hostname
	}
	c.Database = strings.TrimSpace(c.Database)
	c.Host = strings.TrimSpace(c.Host)
	c.Username = strings.TrimSpace(c.Username)
	c.Parameters = strings.TrimSpace(c.Parameters)

	if c.DSN == "" && c.sqlDB == nil && c.connPool == nil {
		switch c.Driver {
		case "sqlite":
			if c.Database == "" {
				return ConnectionConfig{}, fmt.Errorf("gorm: sqlite connection %q requires database or dsn", c.Name)
			}
		default:
			if c.Host == "" || c.Database == "" || c.Username == "" {
				return ConnectionConfig{}, fmt.Errorf("gorm: connection %q requires host, username and database when dsn is empty", c.Name)
			}
		}
	}

	if c.MaxIdleConns <= 0 {
		c.MaxIdleConns = defaultMaxIdleConns
	}
	if c.MaxOpenConns <= 0 {
		c.MaxOpenConns = defaultMaxOpenConns
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = defaultConnMaxLifetimeSeconds
	}
	if c.ConnMaxIdleTime == 0 {
		c.ConnMaxIdleTime = defaultConnMaxIdleTimeSeconds
	}
	if c.PingTimeout <= 0 {
		c.PingTimeout = settings.PingTimeout
	}

	c.effectivePrepareStmt = resolveBool(c.PrepareStmt, settings.PrepareStmt, true)
	c.effectiveSkipDefaultTransaction = resolveBool(c.SkipDefaultTransaction, settings.SkipDefaultTransaction, true)
	c.effectiveDisableNestedTransaction = resolveBool(c.DisableNestedTransaction, settings.DisableNestedTransaction, false)
	c.effectivePingTimeout = secondsToDuration(c.PingTimeout, defaultPingTimeoutSeconds)

	logLevel, err := parseLogLevel(firstNonEmpty(c.LogLevel, settings.Logger.Level))
	if err != nil {
		return ConnectionConfig{}, fmt.Errorf("gorm: connection %q parse log level: %w", c.Name, err)
	}
	logEnabled := resolveBool(c.LogEnabled, settings.Logger.Enabled, true)
	if !logEnabled {
		logLevel = gormlogger.Silent
	}
	c.effectiveLogger = resolvedLogConfig{
		enabled:                   logEnabled,
		level:                     logLevel,
		slowThreshold:             time.Duration(resolvePositiveInt(c.SlowThresholdMillis, settings.Logger.SlowThresholdMillis, defaultSlowThresholdMillis)) * time.Millisecond,
		ignoreRecordNotFoundError: resolveBool(c.IgnoreRecordNotFoundError, settings.Logger.IgnoreRecordNotFoundError, true),
	}

	return c, nil
}

func parseLogLevel(raw string) (gormlogger.LogLevel, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "warn", "warning", "3":
		return gormlogger.Warn, nil
	case "silent", "off", "0", "1":
		return gormlogger.Silent, nil
	case "error", "2":
		return gormlogger.Error, nil
	case "info", "4":
		return gormlogger.Info, nil
	default:
		return gormlogger.Warn, fmt.Errorf("unsupported gorm log level %q", raw)
	}
}

func normalizeDriver(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "mysql":
		return "mysql"
	case "pgsql", "postgres", "postgresql":
		return "postgres"
	case "sqlserver", "sqlsrv", "mssql":
		return "sqlserver"
	case "clickhouse":
		return "clickhouse"
	case "sqlite", "sqlite3":
		return "sqlite"
	default:
		return ""
	}
}

func resolveBool(current, fallback *bool, defaultValue bool) bool {
	switch {
	case current != nil:
		return *current
	case fallback != nil:
		return *fallback
	default:
		return defaultValue
	}
}

func resolvePositiveInt(current, fallback, defaultValue int) int {
	switch {
	case current > 0:
		return current
	case fallback > 0:
		return fallback
	default:
		return defaultValue
	}
}

func secondsToDuration(seconds, defaultSeconds int) time.Duration {
	if seconds < 0 {
		return 0
	}
	if seconds == 0 {
		seconds = defaultSeconds
	}
	return time.Duration(seconds) * time.Second
}

func boolPtr(v bool) *bool { return &v }

func (s Settings) enabled() bool {
	return len(s.Connections) > 0
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
