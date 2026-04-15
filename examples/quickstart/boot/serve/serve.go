package serve

import (
	"flag"
	"strings"
	"time"

	"github.com/yuancore/go-zen/adapter/db/zdb"
	"github.com/yuancore/go-zen/adapter/http/zgin"
	"github.com/yuancore/go-zen/adapter/logger/zlog"
	"github.com/yuancore/go-zen/examples/quickstart/boot/router"
	"github.com/yuancore/go-zen/zen"
)

// Run is the CLI entrypoint.
// Flags:
//
//	-config-dir ./config       directory containing config.toml + env overlays (default)
//	-config     ./config/x.toml  [legacy] explicit single-file path (overrides -config-dir)
func Run() error {
	configDir := flag.String("config-dir", "./config", "directory containing config.toml and env overlays")
	legacyPath := flag.String("config", "", "[legacy] explicit config file path (overrides -config-dir)")
	flag.Parse()

	if *legacyPath != "" {
		return StartFromFile(*legacyPath)
	}
	return Start(*configDir)
}

// Start loads env-based config from dir and starts the application.
// APP_ENV selects the overlay: config.toml + config_{APP_ENV}.toml.
func Start(configDir string) error {
	cfg, err := loadEnvConfig(configDir)
	if err != nil {
		return err
	}
	app, err := newApp(cfg)
	if err != nil {
		return err
	}
	return app.Run(listenAddr(cfg))
}

// StartFromFile loads a single explicit config file (legacy / CI usage).
func StartFromFile(configPath string) error {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return err
	}
	app, err := newApp(cfg)
	if err != nil {
		return err
	}
	return app.Run(listenAddr(cfg))
}

// NewApp constructs the application without starting it (used by tests).
func NewApp(configPath string) (*zen.App, error) {
	cfg, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return newApp(cfg)
}

// newApp is the shared constructor used by all entry points.
func newApp(cfg zen.Config) (*zen.App, error) {
	app := zen.New(
		zen.WithName(name(cfg)),
		zen.WithConfig(cfg),
		zlog.Factory(), // builds ZapLogger from [log] config section
		zgin.Factory(), // builds GinEngine using the constructed logger
		zen.WithStopTimeout(10*time.Second),
	)

	// zdb.New() auto-reads [[connections]] from config and opens all DB connections.
	app.Use(zdb.New())

	// Pre-inject all DB connections into each request's context so service/DAO
	// code can call zdb.DBCtx(app, ctx) or zdb.FromContext(ctx) without *App.
	app.OnStart(func() error {
		app.Middleware(zdb.InjectMiddleware(app))
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
