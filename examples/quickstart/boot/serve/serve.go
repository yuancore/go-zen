package serve

import (
	"flag"
	"strings"
	"time"

	gormadapter "github.com/yuancore/go-zen/adapter/db/gorm"
	zgin "github.com/yuancore/go-zen/adapter/http/gin"
	zlog "github.com/yuancore/go-zen/adapter/logger/zap"
	"github.com/yuancore/go-zen/examples/quickstart/boot/router"
	"github.com/yuancore/go-zen/zen"
)

// RunCLI is the CLI entrypoint.
func RunCLI() error {
	configPath := flag.String("config", "config/config.toml", "configuration file path")
	flag.Parse()
	return Run(*configPath)
}

// Run loads config and starts the application.
func Run(configPath string) error {
	app, err := Build(configPath)
	if err != nil {
		return err
	}
	return app.Run(serverAddress(app.Config()))
}

// Build constructs the application without starting it.
// Component Init (including DB connections) happens inside app.Run().
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

	// Register database component: reads [database] / [[connections]] from config.
	// Init is deferred until app.Run() kicks off the component lifecycle.
	app.Use(gormadapter.New())

	// Register routes after all components have been initialized so the DB is
	// already available in the container via gormadapter.DB(app, ctx).
	app.OnStart(func() error {
		return router.Register(app)
	})

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
