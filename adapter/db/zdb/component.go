package zdb

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	mysqlcfg "github.com/go-sql-driver/mysql"
	"github.com/yuancore/go-zen/zen"
	"gorm.io/driver/clickhouse"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// ConfigProvider is the minimal config surface required by this component.
type ConfigProvider interface {
	IsSet(key string) bool
	Unmarshal(key string, v any) error
}

// Option customizes the GORM component.
type Option func(*Component)

// Component manages GORM connections and injects them into the zen container.
type Component struct {
	zen.BaseComponent

	configKey          string
	serviceName        string
	managerServiceName string
	settings           *Settings

	manager   *Manager
	logger    zen.Logger
	closeOnce sync.Once
}

// New creates a production-oriented GORM component.
func New(opts ...Option) *Component {
	c := &Component{
		BaseComponent: zen.BaseComponent{
			ComponentName: "gorm",
		},
		configKey:          DefaultConfigKey,
		serviceName:        DefaultServiceName,
		managerServiceName: DefaultManagerServiceName,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithConfigKey overrides the config key used to load settings.
func WithConfigKey(key string) Option {
	return func(c *Component) {
		c.configKey = strings.TrimSpace(key)
	}
}

// WithSettings injects explicit settings, useful for tests or fully code-driven apps.
func WithSettings(settings Settings) Option {
	return func(c *Component) {
		cloned := settings
		cloned.Connections = append([]ConnectionConfig(nil), settings.Connections...)
		c.settings = &cloned
	}
}

// WithServiceName overrides the service name used for the default *gorm.DB.
func WithServiceName(name string) Option {
	return func(c *Component) {
		if trimmed := strings.TrimSpace(name); trimmed != "" {
			c.serviceName = trimmed
		}
	}
}

// WithManagerServiceName overrides the service name used for the manager registry.
func WithManagerServiceName(name string) Option {
	return func(c *Component) {
		if trimmed := strings.TrimSpace(name); trimmed != "" {
			c.managerServiceName = trimmed
		}
	}
}

// Init opens connections and injects them into the app container.
func (c *Component) Init(app *zen.App) error {
	c.logger = app.Logger().With("component", c.Name())

	settings, err := c.resolveSettings(app.Config())
	if err != nil {
		return err
	}
	if !settings.enabled() {
		c.logger.Info("gorm skipped")
		return nil
	}

	manager, err := openManager(settings, app.Logger())
	if err != nil {
		return err
	}
	c.manager = manager

	app.Provide(c.managerServiceName, manager)
	for _, name := range manager.Names() {
		db := manager.MustGet(name)
		app.Provide(namedService(c.serviceName, name), db)
	}
	app.Provide(c.serviceName, manager.MustDefault())

	c.logger.Info("gorm initialized", "default", manager.DefaultName(), "connections", manager.Names())
	return nil
}

// Start is a no-op because GORM connections are ready after Init.
func (c *Component) Start() error { return nil }

// Stop closes all opened SQL connections.
func (c *Component) Stop(_ context.Context) error {
	var err error
	c.closeOnce.Do(func() {
		if c.manager != nil {
			err = c.manager.Close()
			if err != nil {
				c.logger.Error("gorm shutdown failed", "err", err)
				return
			}
			c.logger.Info("gorm shutdown complete")
		}
	})
	return err
}

// Resolve returns the default *gorm.DB registered under the default service key.
func Resolve(app *zen.App) (*gorm.DB, bool) {
	return zen.ResolveAs[*gorm.DB](app, DefaultServiceName)
}

// MustResolve returns the default *gorm.DB and panics if missing.
func MustResolve(app *zen.App) *gorm.DB {
	db, ok := Resolve(app)
	if !ok {
		panic("gorm: default database service not found")
	}
	return db
}

// ResolveNamed returns a named *gorm.DB registered by this component.
func ResolveNamed(app *zen.App, name string) (*gorm.DB, bool) {
	return zen.ResolveAs[*gorm.DB](app, NamedService(name))
}

// MustResolveNamed returns a named *gorm.DB or panics.
func MustResolveNamed(app *zen.App, name string) *gorm.DB {
	db, ok := ResolveNamed(app, name)
	if !ok {
		panic("gorm: named database service not found: " + name)
	}
	return db
}

// ResolveManager returns the connection manager registered by this component.
func ResolveManager(app *zen.App) (*Manager, bool) {
	return zen.ResolveAs[*Manager](app, DefaultManagerServiceName)
}

// MustResolveManager returns the connection manager or panics.
func MustResolveManager(app *zen.App) *Manager {
	manager, ok := ResolveManager(app)
	if !ok {
		panic("gorm: manager service not found")
	}
	return manager
}

// NamedService returns the container key used for a named database.
func NamedService(name string) string {
	return namedService(DefaultServiceName, name)
}

func namedService(prefix, name string) string {
	return strings.TrimSpace(prefix) + "." + strings.TrimSpace(name)
}

func (c *Component) resolveSettings(cfg zen.Config) (Settings, error) {
	if c.settings != nil {
		return c.settings.normalize()
	}
	return loadSettings(cfg, c.configKey)
}

func openManager(settings Settings, logger zen.Logger) (*Manager, error) {
	manager := newManager(settings.Default)
	for _, conn := range settings.Connections {
		db, err := openConnection(conn, logger)
		if err != nil {
			_ = manager.Close()
			return nil, fmt.Errorf("gorm: open connection %q: %w", conn.Name, err)
		}
		if err := manager.register(conn.Name, db); err != nil {
			_ = manager.Close()
			return nil, err
		}
	}
	return manager, nil
}

func openConnection(conn ConnectionConfig, logger zen.Logger) (*gorm.DB, error) {
	dsn, err := buildDSN(conn)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialector(conn, dsn), &gorm.Config{
		SkipDefaultTransaction:   conn.effectiveSkipDefaultTransaction,
		DisableNestedTransaction: conn.effectiveDisableNestedTransaction,
		DisableAutomaticPing:     true,
		TranslateError:           true,
		PrepareStmt:              conn.effectivePrepareStmt,
		Logger: newSQLLogger(
			logger.With("db", conn.Name, "driver", conn.Driver),
			conn.effectiveLogger,
		),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(conn.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conn.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(secondsToDuration(conn.ConnMaxLifetime, defaultConnMaxLifetimeSeconds))
	sqlDB.SetConnMaxIdleTime(secondsToDuration(conn.ConnMaxIdleTime, defaultConnMaxIdleTimeSeconds))

	if conn.effectivePingTimeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), conn.effectivePingTimeout)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			_ = sqlDB.Close()
			return nil, err
		}
	}

	return db, nil
}

func dialector(conn ConnectionConfig, dsn string) gorm.Dialector {
	connPool := conn.connPool
	if connPool == nil && conn.sqlDB != nil {
		connPool = conn.sqlDB
	}

	switch conn.Driver {
	case "mysql":
		precision := 3
		return gormmysql.New(gormmysql.Config{
			DSN:                       dsn,
			Conn:                      connPool,
			DefaultStringSize:         256,
			DefaultDatetimePrecision:  &precision,
			DisableDatetimePrecision:  false,
			DontSupportRenameIndex:    false,
			DontSupportRenameColumn:   false,
			SkipInitializeWithVersion: connPool != nil,
		})
	case "postgres":
		return postgres.New(postgres.Config{
			DSN:                  dsn,
			Conn:                 connPool,
			PreferSimpleProtocol: true,
		})
	case "sqlserver":
		return sqlserver.New(sqlserver.Config{
			DSN:  dsn,
			Conn: connPool,
		})
	case "clickhouse":
		return clickhouse.New(clickhouse.Config{
			DSN:                       dsn,
			Conn:                      connPool,
			SkipInitializeWithVersion: connPool != nil,
		})
	case "sqlite":
		return sqlite.New(sqlite.Config{
			DSN:  dsn,
			Conn: connPool,
		})
	default:
		panic("gorm: unsupported driver: " + conn.Driver)
	}
}

func buildDSN(conn ConnectionConfig) (string, error) {
	if strings.TrimSpace(conn.DSN) != "" {
		return conn.DSN, nil
	}

	switch conn.Driver {
	case "mysql":
		cfg := mysqlcfg.NewConfig()
		cfg.User = conn.Username
		cfg.Passwd = conn.Password
		cfg.Net = "tcp"
		cfg.Addr = netAddress(conn.Host, conn.Port)
		cfg.DBName = conn.Database
		cfg.Params = make(map[string]string)
		if err := applyQueryParams(cfg.Params, conn.Parameters); err != nil {
			return "", fmt.Errorf("mysql params: %w", err)
		}
		return cfg.FormatDSN(), nil
	case "postgres":
		values, err := parseQueryParams(conn.Parameters)
		if err != nil {
			return "", fmt.Errorf("postgres params: %w", err)
		}
		u := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword(conn.Username, conn.Password),
			Host:     netAddress(conn.Host, conn.Port),
			Path:     conn.Database,
			RawQuery: values.Encode(),
		}
		return u.String(), nil
	case "sqlserver":
		values, err := parseQueryParams(conn.Parameters)
		if err != nil {
			return "", fmt.Errorf("sqlserver params: %w", err)
		}
		if conn.Database != "" && values.Get("database") == "" {
			values.Set("database", conn.Database)
		}
		u := url.URL{
			Scheme:   "sqlserver",
			User:     url.UserPassword(conn.Username, conn.Password),
			Host:     netAddress(conn.Host, conn.Port),
			RawQuery: values.Encode(),
		}
		return u.String(), nil
	case "clickhouse":
		values, err := parseQueryParams(conn.Parameters)
		if err != nil {
			return "", fmt.Errorf("clickhouse params: %w", err)
		}
		u := url.URL{
			Scheme:   "clickhouse",
			User:     url.UserPassword(conn.Username, conn.Password),
			Host:     netAddress(conn.Host, conn.Port),
			Path:     conn.Database,
			RawQuery: values.Encode(),
		}
		return u.String(), nil
	case "sqlite":
		return conn.Database, nil
	default:
		return "", fmt.Errorf("unsupported driver %q", conn.Driver)
	}
}

func parseQueryParams(raw string) (url.Values, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return url.Values{}, nil
	}
	values, err := url.ParseQuery(raw)
	if err != nil {
		return nil, err
	}
	return values, nil
}

func applyQueryParams(target map[string]string, raw string) error {
	values, err := parseQueryParams(raw)
	if err != nil {
		return err
	}
	for key, vals := range values {
		if len(vals) == 0 {
			continue
		}
		target[key] = vals[len(vals)-1]
	}
	return nil
}

func netAddress(host string, port int) string {
	if port <= 0 {
		return host
	}
	return host + ":" + strconv.Itoa(port)
}
