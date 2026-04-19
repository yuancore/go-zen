package serve

import (
	"flag"
	"strings"
	"time"

	"github.com/yuancore/go-zen/adapter/cache/zredis"
	"github.com/yuancore/go-zen/adapter/db/zdb"
	"github.com/yuancore/go-zen/adapter/http/zgin"
	"github.com/yuancore/go-zen/adapter/logger/zlog"
	"github.com/yuancore/go-zen/examples/quickstart/boot/router"
	"github.com/yuancore/go-zen/zen"
)

// Run 是 CLI 入口点，解析 -config 标志后启动应用。
// Run is the CLI entrypoint. It parses the -config flag and starts the application.
func Run() error {
	configDir := flag.String("config", "./config/config_dev.toml", "Configuration file path")
	flag.Parse()
	return Start(*configDir)
}

// Start 从指定路径加载配置并启动应用。
// Start loads config from the given path and starts the application.
func Start(configDir string) error {
	cfg, err := loadConfig(configDir)
	if err != nil {
		return err
	}
	app, err := newApp(cfg)
	if err != nil {
		return err
	}
	return app.Run(listenAddr(cfg))
}

// newApp is the shared constructor used by all entry points.
func newApp(cfg zen.Config) (*zen.App, error) {
	app := zen.New(
		zen.WithName(name(cfg)),
		zen.WithConfig(cfg),
		zlog.Factory(),              // 从 [log] 配置节构建 ZapLogger / builds ZapLogger from [log] config section
		zgin.FactoryFromConfig(cfg), // 构建 GinEngine，访问日志由 [http_log] 驱动 / builds GinEngine; access log driven by [http_log]
		zen.WithStopTimeout(10*time.Second),
	)

	// zdb.New() 自动读取 [[connections]] 并打开所有数据库连接，同时注册 context 注入中间件。
	// zdb.New() auto-reads [[connections]], opens all DB connections, and registers the inject middleware.
	app.Use(zdb.New())

	// zredis.New() 自动读取 [[redis]] 并打开所有 Redis 连接，同时注册 context 注入中间件。
	// 默认客户端：zredis.MustResolve(app)；具名客户端：zredis.MustResolveNamed(app, "session")
	// zredis.New() auto-reads [[redis]], opens all Redis connections, and registers the inject middleware.
	// Default client: zredis.MustResolve(app)  Named client: zredis.MustResolveNamed(app, "session")
	app.Use(zredis.New())

	// 所有 DB 和 Redis 的 context 注入中间件已由各组件自动注册，此处仅需设置路由。
	// DB and Redis context-injection middlewares are auto-registered by each component above.
	// Only route setup is needed here.
	app.OnStart(func() error {
		return router.Setup(app)
	})

	return app, nil
}

// listenAddr 从配置中读取监听地址，确保带 ":" 前缀。
// listenAddr reads the listen address from config, ensuring it has a ":" prefix.
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

// name 从配置中读取应用名称，默认为 "go-zen-quickstart"。
// name reads the application name from config, defaulting to "go-zen-quickstart".
func name(cfg zen.Config) string {
	if name := strings.TrimSpace(cfg.GetString("system.app_name")); name != "" {
		return name
	}
	return "go-zen-quickstart"
}
