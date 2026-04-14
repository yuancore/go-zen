package serve

import (
	"flag"
	"strings"
	"time"

	"github.com/yuancore/go-zen/adapter/http/zgin"
	"github.com/yuancore/go-zen/adapter/logger/zlog"
	"github.com/yuancore/go-zen/examples/quickstart/boot/router"
	"github.com/yuancore/go-zen/zen"
)

// Run is the CLI entrypoint.
func Run() error {
	configPath := flag.String("config", "./config/config_dev.toml", "configuration file path")
	flag.Parse()
	return Start(*configPath)
}

// Start loads config and starts the application.
func Start(configPath string) error {
	app, err := NewApp(configPath)
	if err != nil {
		return err
	}
	return app.Run(listenAddr(app.Config()))
}

// NewApp constructs the application without starting it.
// Component Init (including DB connections) happens inside app.Run().
func NewApp(configPath string) (*zen.App, error) {
	cfg, err := loadConfig(configPath)

	if err != nil {
		return nil, err
	}

	logger, err := zlog.New(cfg)

	if err != nil {
		return nil, err
	}

	app := zen.New(
		zen.WithName(name(cfg)),
		zen.WithConfig(cfg),
		zen.WithLogger(logger),
		zen.WithEngine(zgin.NewEngine(logger)),
		zen.WithStopTimeout(10*time.Second),
	)

	// Register database component: reads [database] / [[connections]] from config.
	// Init is deferred until app.Run() kicks off the component lifecycle.

	if dbComponent, err := loadZdb(cfg); err != nil {
		return nil, err
	} else if dbComponent != nil {
		app.Use(dbComponent)
	}

	// Register routes after all components have been initialized so the DB is
	// already available in the container via gormadapter.DB(app, ctx).
	app.OnStart(func() error {
		return router.Setup(app)
	})

	return app, nil
}

func listenAddr(cfg zen.Config) string {
	addr := strings.TrimSpace(cfg.GetString("system.listen"))
	if addr == "" {
		addr = "8080"
	}
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}
	return addr
}

func name(cfg zen.Config) string {
	if name := strings.TrimSpace(cfg.GetString("system.app_name")); name != "" {
		return name
	}
	return "go-zen-quickstart"
}
