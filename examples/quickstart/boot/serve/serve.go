package serve

import (
	"context"
	"flag"
	"strings"
	"time"

	zgin "github.com/yuancore/go-zen/adapter/http/gin"
	zlog "github.com/yuancore/go-zen/adapter/logger/zap"
	"github.com/yuancore/go-zen/examples/quickstart/app/http/api"
	"github.com/yuancore/go-zen/examples/quickstart/app/service"
	"github.com/yuancore/go-zen/examples/quickstart/boot/router"
	"github.com/yuancore/go-zen/zen"
)

// RunCLI 命令行入口
func RunCLI() error {
	configPath := flag.String("config", "config/config.toml", "configuration file path")
	flag.Parse()
	return Run(*configPath)
}

// Run 启动应用
func Run(configPath string) error {
	app, err := Build(configPath)
	if err != nil {
		return err
	}
	return app.Run(serverAddress(app.Config()))
}

// Build 构建应用实例
func Build(configPath string) (*zen.App, error) {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	logger, err := zlog.NewLoggerFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	app := zen.New(
		zen.Name(appName(cfg)),
		zen.WithConfig(cfg),
		zen.WithLogger(logger),
		zen.WithEngine(zgin.NewEngine(logger)),
		zen.StopTimeout(10*time.Second),
	)

	// 初始化数据库（使用 connections 数组中的第一个连接，可根据需要修改）
	db, err := newDB(cfg)
	if err != nil {
		return nil, err
	}
	app.Use(db)

	// 依赖组装

	userService := service.NewUserService(userDAO, logger)
	systemAPI := api.NewSystemAPI(app, logger)
	userAPI := api.NewUserAPI(userService, logger)

	app.OnStart(func() error {
		return userService.Migrate(context.Background())
	})

	if err := router.Register(app, systemAPI, userAPI); err != nil {
		return nil, err
	}

	return app, nil
}

func serverAddress(cfg zen.Config) string {
	addr := strings.TrimSpace(cfg.GetString("system.address"))
	if addr == "" {
		addr = "8080"
	}
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}
	return addr
}

func appName(cfg zen.Config) string {
	if name := strings.TrimSpace(cfg.GetString("system.app_name")); name != "" {
		return name
	}
	return "go-zen-quickstart"
}
