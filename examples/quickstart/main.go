package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	viperAdapter "github.com/yuancore/go-zen/adapter/config/viper"
	ginAdapter "github.com/yuancore/go-zen/adapter/http/gin"
	zapAdapter "github.com/yuancore/go-zen/adapter/logger/zap"
	"github.com/yuancore/go-zen/zen"
	"go.uber.org/zap"
)

func main() {
	configPath := flag.String("config", "./examples/quickstart/config/config.toml", "Configuration file path")
	flag.Parse()

	cfg := viperAdapter.NewConfig(*configPath)
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

	// Register components here when needed.
	// app.Register(dbModule)

	app.GET("/ping", func(c zen.Context) {
		fmt.Println(app.Config().GetString("logger.encoding"))
		logger.Error("Received ping request", zap.String("version", app.Config().GetString("logger.encoding")))
		c.JSON(200, map[string]string{"pong": "ok"})
	})

	if err := app.Run(":8080"); err != nil {
		log.Fatalf("run failed: %v", err)
	}
}
