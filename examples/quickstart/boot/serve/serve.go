package serve

import (
	"flag"
	"os"
	"time"

	zgin "github.com/yuancore/go-zen/adapter/http/gin"
	zlog "github.com/yuancore/go-zen/adapter/logger/zap"
	"github.com/yuancore/go-zen/examples/quickstart/boot/router"
	"github.com/yuancore/go-zen/zen"
)

func RunCLI(defaultConfigPath string) error {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	configPath := flagSet.String("config", defaultConfigPath, "Configuration file path")
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return err
	}
	return Run(*configPath)
}

func Run(configPath string) error {
	app, err := Build(configPath)
	if err != nil {
		return err
	}

	runErr := app.Run(serverAddress(app.Config()))
	syncLogger(app.Logger())
	return runErr
}

func Build(configPath string) (*zen.App, error) {
	config, baseDir, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	logger, err := zlog.NewLoggerFromConfig(config)
	if err != nil {
		return nil, err
	}

	app := zen.New(
		zen.Name(appName(config)),
		zen.WithConfig(config),
		zen.WithLogger(logger),
		zen.WithEngine(zgin.NewEngine(logger)),
		zen.StopTimeout(10*time.Second),
	)

	dbComponent, err := newDatabaseComponent(config, baseDir)
	if err != nil {
		return nil, err
	}
	app.Use(dbComponent)

	services := newServices(app)
	registerStartupHooks(app, services.userService)

	if err := router.Register(app, services.systemAPI, services.userAPI); err != nil {
		return nil, err
	}

	return app, nil
}

func syncLogger(logger zen.Logger) {
	if synced, ok := logger.(interface{ Sync() error }); ok {
		_ = synced.Sync()
	}
}
