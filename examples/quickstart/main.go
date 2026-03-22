package main

import (
	"flag"
	"log"
	"time"

	zconfig "github.com/yuancore/go-zen/adapter/config/viper"
	zdb "github.com/yuancore/go-zen/adapter/db/gorm"
	zgin "github.com/yuancore/go-zen/adapter/http/gin"
	zlog "github.com/yuancore/go-zen/adapter/logger/zap"
	"github.com/yuancore/go-zen/zen"
	"go.uber.org/zap"
)

const defaultConfigPath = "./examples/quickstart/config/config.toml"

func main() {
	configPath := flag.String("config", defaultConfigPath, "Configuration file path")
	flag.Parse()

	//configFile, err := resolveConfigPath(*configPath)
	//if err != nil {
	//	log.Fatalf("resolve config path failed: %v", err)
	//}
	//fmt.Println(configFile)

	config := zconfig.NewConfig()
	if err := config.Load(*configPath); err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	logger, err := zlog.NewLoggerFromConfig(config)
	if err != nil {
		log.Fatalf("init logger failed: %v", err)
	}
	defer func() { _ = logger.Sync() }()

	eng := zgin.NewEngine(logger)

	app := zen.New(
		zen.WithConfig(config),
		zen.WithLogger(logger),
		zen.WithEngine(eng),
		zen.WithStopTimeout(10*time.Second),
	)

	app.Use(
		zdb.New(),
	)

	app.OnStart(func() error {
		manager, ok := zdb.ResolveManager(app)
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
		if db, ok := zdb.Resolve(app); ok {
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
		zdb.DefaultSettings()
		c.JSON(200, map[string]string{"pong": "ok"})
	})
	const (
		DBMain   = "mysql1"
		DBReport = "mysql_report"
	)

	app.GET("/users", func(c zen.Context) {
		var list = make([]map[string]interface{}, 0)
		if err := zdb.Connection(app, c, DBMain).
			Table("sys_admin_users").
			Limit(20).
			Find(&list).Error; err != nil {
			c.JSON(500, map[string]any{"error": err.Error()})
			return
		}
		c.JSON(200, map[string]any{
			"list": list,
		})
	})

	app.GET("/dbs", func(c zen.Context) {
		manager, ok := zdb.ResolveManager(app)
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
