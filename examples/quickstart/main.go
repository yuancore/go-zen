package main

import (
	"log"
	"time"

	viperAdapter "github.com/yuancore/go-zen/adapter/config/viper"
	ginAdapter "github.com/yuancore/go-zen/adapter/http/gin"
	zapAdapter "github.com/yuancore/go-zen/adapter/logger/zap"
	"github.com/yuancore/go-zen/zen"
)

func main() {
	cfg := viperAdapter.NewConfig("./config.yaml")
	logger := zapAdapter.NewLogger()
	eng := ginAdapter.NewEngine(logger)

	app := zen.New(
		zen.WithConfig(cfg),
		zen.WithLogger(logger),
		zen.WithEngine(eng),
		zen.WithStopTimeout(10*time.Second),
	)

	// 注册模块（示例模块略）
	// app.Register(dbModule)

	app.GET("/ping", func(c zen.Context) {
		c.JSON(200, map[string]string{"pong": "ok"})
	})

	if err := app.Run(":8080"); err != nil {
		log.Fatalf("run failed: %v", err)
	}
}
