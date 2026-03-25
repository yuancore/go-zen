package serve

import (
	"fmt"
	"strings"

	zdb "github.com/yuancore/go-zen/adapter/db/gorm"
	"github.com/yuancore/go-zen/zen"
)

func newDatabaseComponent(cfg zen.Config, baseDir string) (*zdb.Component, error) {
	if cfg == nil || (!cfg.IsSet("database") && !cfg.IsSet("connections")) {
		return zdb.New(), nil
	}

	settings := zdb.DefaultSettings()
	switch {
	case cfg.IsSet("database"):
		if err := cfg.Unmarshal("database", &settings); err != nil {
			return nil, fmt.Errorf("unmarshal database config: %w", err)
		}
	case cfg.IsSet("connections"):
		if err := cfg.Unmarshal("", &settings); err != nil {
			return nil, fmt.Errorf("unmarshal database config: %w", err)
		}
	}

	normalizeSQLiteConnections(&settings, baseDir)
	return zdb.New(zdb.WithSettings(settings)), nil
}

func normalizeSQLiteConnections(settings *zdb.Settings, baseDir string) {
	if settings == nil {
		return
	}

	for index := range settings.Connections {
		conn := &settings.Connections[index]
		if !isSQLite(conn.Driver, conn.Type) {
			continue
		}
		if conn.Database != "" {
			conn.Database = resolveSQLitePath(baseDir, conn.Database)
		}
		if conn.DSN != "" {
			conn.DSN = resolveSQLitePath(baseDir, conn.DSN)
		}
	}
}

func isSQLite(driver, fallback string) bool {
	driver = strings.ToLower(strings.TrimSpace(driver))
	if driver == "" {
		driver = strings.ToLower(strings.TrimSpace(fallback))
	}
	return driver == "sqlite" || driver == "sqlite3"
}

func resolveSQLitePath(baseDir, raw string) string {
	raw = strings.TrimSpace(raw)
	switch {
	case raw == "", raw == ":memory:", strings.Contains(raw, "mode=memory"):
		return raw
	case strings.Contains(raw, "://"):
		return raw
	default:
		return resolveLocalPath(baseDir, raw)
	}
}
