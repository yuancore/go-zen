package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	viperAdapter "github.com/yuancore/go-zen/adapter/config/viper"
	gormAdapter "github.com/yuancore/go-zen/adapter/db/gorm"
	ginAdapter "github.com/yuancore/go-zen/adapter/http/gin"
	zapAdapter "github.com/yuancore/go-zen/adapter/logger/zap"
	"github.com/yuancore/go-zen/zen"
	"go.uber.org/zap"
)

const defaultConfigPath = "./examples/quickstart/config/config.toml"

func main() {
	configPath := flag.String("config", defaultConfigPath, "Configuration file path")
	flag.Parse()

	configFile, err := resolveConfigPath(*configPath)
	if err != nil {
		log.Fatalf("resolve config path failed: %v", err)
	}
	fmt.Println(configFile)

	cfg := viperAdapter.NewConfig()
	if err := cfg.Load(configFile); err != nil {
		log.Fatalf("load config failed: %v", err)
	}
	normalizeLoggerPaths(cfg, configBaseDir(configFile))

	logger, err := zapAdapter.NewLoggerFromConfig(cfg)
	if err != nil {
		log.Fatalf("init logger failed: %v", err)
	}
	defer func() { _ = logger.Sync() }()

	eng := ginAdapter.NewEngine(logger)

	app := zen.New(
		zen.WithConfig(cfg),
		zen.WithLogger(logger),
		zen.WithEngine(eng),
		zen.WithStopTimeout(10*time.Second),
	)

	app.Use(
		gormAdapter.New(),
	)

	app.OnStart(func() error {
		manager, ok := gormAdapter.ResolveManager(app)
		if !ok {
			logger.Info("database not configured, skip init")
			return nil
		}
		logger.Info(
			"database connections initialized",
			zap.String("default", manager.DefaultName()),
			zap.Strings("connections", manager.Names()),
		)
		return nil
	})

	app.GET("/ping", func(c zen.Context) {
		if db, ok := gormAdapter.Resolve(app); ok {
			sqlDB, err := db.DB()
			if err != nil {
				logger.Error("resolve sql db failed", zap.Error(err))
				c.JSON(500, map[string]string{"error": err.Error()})
				return
			}
			if err := sqlDB.PingContext(c.Request().Context()); err != nil {
				logger.Error("ping default database failed", zap.Error(err))
				c.JSON(500, map[string]string{"error": err.Error()})
				return
			}
		}

		logger.Info("Received ping request", zap.String("version", app.Config().GetString("system.app_name")))
		c.JSON(200, map[string]string{"pong": "ok"})
	})

	app.GET("/db", func(c zen.Context) {
		gormAdapter.DefaultSettings()
		c.JSON(200, map[string]string{"pong": "ok"})
	})
	const (
		DBMain   = "mysql1"
		DBReport = "mysql_report"
	)

	app.GET("/users", func(c zen.Context) {
		db := gormAdapter.MustResolveNamed(app, DBMain)
		var list = make([]map[string]interface{}, 0)
		if err := db.WithContext(c.Request().Context()).
			Limit(20).
			Table("sys_admin_users").Find(&list).Error; err != nil {
			c.JSON(500, map[string]any{"error": err.Error()})
			return
		}
		c.JSON(200, map[string]any{
			"list": list,
		})
	})

	app.GET("/dbs", func(c zen.Context) {
		manager, ok := gormAdapter.ResolveManager(app)
		if !ok {
			c.JSON(200, map[string]any{
				"default":     "",
				"connections": []string{},
			})
			return
		}
		c.JSON(200, map[string]any{
			"default":     manager.DefaultName(),
			"connections": manager.Names(),
		})
	})

	if err := app.Run(":" + app.Config().GetString("system.address")); err != nil {
		log.Fatalf("run failed: %v", err)
	}
}

func normalizeLoggerPaths(cfg *viperAdapter.ViperConfig, baseDir string) {
	if cfg == nil || strings.TrimSpace(baseDir) == "" {
		return
	}

	normalizePath := func(value string) string {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" || filepath.IsAbs(trimmed) || trimmed == "stdout" || trimmed == "stderr" {
			return trimmed
		}
		return filepath.Join(baseDir, trimmed)
	}

	if path := cfg.GetString("log.path"); path != "" {
		cfg.Set("log.path", normalizePath(path))
	}
	if paths := cfg.GetStringSlice("log.output_paths"); len(paths) > 0 {
		normalized := make([]string, 0, len(paths))
		for _, path := range paths {
			normalized = append(normalized, normalizePath(path))
		}
		cfg.Set("log.output_paths", normalized)
	}
	if paths := cfg.GetStringSlice("log.error_output_paths"); len(paths) > 0 {
		normalized := make([]string, 0, len(paths))
		for _, path := range paths {
			normalized = append(normalized, normalizePath(path))
		}
		cfg.Set("log.error_output_paths", normalized)
	}
}

func configBaseDir(configPath string) string {
	if strings.TrimSpace(configPath) == "" {
		return ""
	}

	dir := filepath.Dir(configPath)
	if strings.EqualFold(filepath.Base(dir), "config") {
		return filepath.Dir(dir)
	}
	return dir
}

func resolveConfigPath(raw string) (string, error) {
	candidates := []string{raw}
	if raw == defaultConfigPath {
		candidates = append(candidates, "./config/config.toml")
	}

	for _, candidate := range candidates {
		abs, err := filepath.Abs(candidate)
		if err != nil {
			return "", err
		}
		if _, err := os.Stat(abs); err == nil {
			return abs, nil
		}
	}

	abs, err := filepath.Abs(raw)
	if err != nil {
		return "", err
	}
	return "", fmt.Errorf("config file not found: %s", abs)
}
