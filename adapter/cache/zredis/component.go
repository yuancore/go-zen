package zredis

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/yuancore/go-zen/zen"
)

// Option 用于自定义 Redis 组件行为。
// Option customizes the Redis component.
type Option func(*Component)

// Component 管理 Redis 连接并将其注入 zen 容器。
//
// 使用示例 / Usage:
//
//	app.Use(zredis.New())                   // 从配置读取 [[redis]] / reads [[redis]] from config
//	app.Use(zredis.New(
//	    zredis.WithConfigKey("my_redis"),   // 自定义配置键 / custom config key
//	))
type Component struct {
	zen.BaseComponent

	configKey          string
	serviceName        string
	managerServiceName string
	instances          []InstanceConfig // 显式覆盖；nil 表示从配置读取 / explicit override; nil = read from config

	manager   *Manager
	closeOnce sync.Once
}

// New 创建 Redis 组件，支持可选配置项。
// 未传入任何选项时，自动读取 [[redis]] 配置节。
// New creates a Redis component with optional customizations.
// When no options are given it reads the [[redis]] config section.
func New(opts ...Option) *Component {
	c := &Component{
		BaseComponent: zen.BaseComponent{
			ComponentName: "redis",
		},
		configKey:          DefaultConfigKey,
		serviceName:        DefaultServiceName,
		managerServiceName: DefaultManagerServiceName,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithConfigKey 覆盖配置节键名（默认："redis"）。
// WithConfigKey overrides the config section key (default: "redis").
func WithConfigKey(key string) Option {
	return func(c *Component) {
		if trimmed := strings.TrimSpace(key); trimmed != "" {
			c.configKey = trimmed
		}
	}
}

// WithInstances 注入显式实例配置，适用于测试或完全代码驱动的场景。
//
//	app.Use(zredis.New(
//	    zredis.WithInstances(zredis.InstanceConfig{
//	        Name:    "default",
//	        Address: "127.0.0.1:6379",
//	    }),
//	))
//
// WithInstances injects explicit instance configurations.
// Useful for tests or fully code-driven setups.
func WithInstances(instances ...InstanceConfig) Option {
	return func(c *Component) {
		c.instances = append(c.instances, instances...)
	}
}

// WithServiceName 覆盖默认 zen.Cache 的容器键（默认："cache"）。
// WithServiceName overrides the container key for the default zen.Cache (default: "cache").
func WithServiceName(name string) Option {
	return func(c *Component) {
		if trimmed := strings.TrimSpace(name); trimmed != "" {
			c.serviceName = trimmed
		}
	}
}

// WithManagerServiceName 覆盖 *Manager 的容器键（默认："redis_manager"）。
// WithManagerServiceName overrides the container key for the *Manager (default: "redis_manager").
func WithManagerServiceName(name string) Option {
	return func(c *Component) {
		if trimmed := strings.TrimSpace(name); trimmed != "" {
			c.managerServiceName = trimmed
		}
	}
}

// ---------- Component lifecycle ----------

// Init 打开 Redis 连接并将其注册到应用容器，同时自动注册 context 注入中间件。
// 配置错误或 Ping 失败均会返回 error，应用启动时 fail-fast 退出。
// 注册的服务键：
//   - "<serviceName>"              → 默认 zen.Cache  （如 "cache"）
//   - "<serviceName>.<name>"       → 具名 zen.Cache  （如 "cache.redis"）
//   - "<managerServiceName>"       → *Manager        （如 "redis_manager"）
//
// Init opens Redis connections, registers them in the app container, and
// auto-registers the context-injection middleware via an OnStart hook.
// Config errors or a failed ping cause Init to return an error (fail-fast).
// Registered services:
//   - "<serviceName>"              → default zen.Cache  (e.g. "cache")
//   - "<serviceName>.<name>"       → named zen.Cache    (e.g. "cache.redis")
//   - "<managerServiceName>"       → *Manager           (e.g. "redis_manager")
func (c *Component) Init(app *zen.App) error {
	logger := app.Logger().With("component", c.Name())

	instances, defaultName, err := c.resolveInstances(app.Config())
	if err != nil {
		return err
	}
	if len(instances) == 0 {
		logger.Info("redis skipped: no instances configured")
		return nil
	}
	if defaultName == "" {
		defaultName = instances[0].Name
	}

	mgr := newManager(defaultName)
	for _, inst := range instances {
		// 校验必填字段：配置字段名写错时 mapstructure 会静默赋零值，此处提前拦截。
		// Validate required fields early; misspelled config keys cause silent zero values.
		if err := inst.validate(); err != nil {
			_ = mgr.Close()
			return err
		}

		client := newClient(inst)
		name := inst.Name // validate() already ensures non-empty / validate() 已保证非空

		// 启动连通性检查 — 与 GORM ping 行为保持一致。
		// 快速失败，避免应用在无缓存状态下静默运行。
		// Startup connectivity check — mirrors the GORM ping behaviour.
		// Fails fast so the process exits instead of silently running without cache.
		if timeout := inst.effectivePingTimeout(); timeout > 0 {
			pingCtx, cancel := context.WithTimeout(context.Background(), timeout)
			pingErr := client.Ping(pingCtx)
			cancel()
			if pingErr != nil {
				// 该客户端尚未加入 mgr，需手动关闭。
				// Close this client before returning — it is not yet in mgr so mgr.Close() won't reach it.
				_ = client.Close()
				_ = mgr.Close()
				return fmt.Errorf("redis: startup ping failed for %q (%s): %w — check address/password in config", name, inst.effectiveAddress(), pingErr)
			}
			logger.Info("redis: ping OK", "name", name, "addr", inst.effectiveAddress())
		}

		if err := mgr.register(name, client); err != nil {
			_ = client.Close()
			_ = mgr.Close()
			return err
		}
	}
	c.manager = mgr

	// 向容器注册 manager 和各客户端实例。
	// Register manager and individual clients in the container.
	app.Provide(c.managerServiceName, mgr)
	for _, name := range mgr.Names() {
		cl, _ := mgr.Get(name)
		app.Provide(namedService(c.serviceName, name), cl)
	}
	// 将默认实例注册到不带名称的服务键下。
	// Register the default instance under the plain service name.
	app.Provide(c.serviceName, mgr.MustDefault())

	// 自动注册 context 注入中间件，无需调用方手动添加。
	// Auto-register the context-injection middleware so callers need not wire it manually.
	app.OnStart(func() error {
		app.Middleware(InjectMiddleware(app))
		return nil
	})

	logger.Info("redis initialized",
		"default", mgr.DefaultName(),
		"instances", mgr.Names(),
	)
	return nil
}

// Start 是空操作；连接在 Init 时已就绪。
// Start is a no-op; connections are ready after Init.
func (c *Component) Start() error { return nil }

// Stop 关闭所有 Redis 连接。
// Stop closes all Redis connections.
func (c *Component) Stop(_ context.Context) error {
	var err error
	c.closeOnce.Do(func() {
		if c.manager != nil {
			if closeErr := c.manager.Close(); closeErr != nil {
				err = closeErr
			}
		}
	})
	return err
}

func (c *Component) resolveInstances(cfg zen.Config) ([]InstanceConfig, string, error) {
	if len(c.instances) > 0 {
		defaultName := ""
		if len(c.instances) > 0 {
			defaultName = c.instances[0].Name
		}
		return c.instances, defaultName, nil
	}
	return loadSettings(cfg, c.configKey)
}

// ---------- Package-level resolution helpers ----------

// Resolve 返回该组件注册的默认 zen.Cache。
// Resolve returns the default zen.Cache registered by this component.
func Resolve(app *zen.App) (zen.Cache, bool) {
	return zen.ResolveAs[zen.Cache](app, DefaultServiceName)
}

// MustResolve 返回默认 zen.Cache，未注册时 panic。
// MustResolve returns the default zen.Cache or panics.
func MustResolve(app *zen.App) zen.Cache {
	c, ok := Resolve(app)
	if !ok {
		panic("redis: default cache service not found")
	}
	return c
}

// ResolveNamed 返回该组件注册的具名 zen.Cache。
// ResolveNamed returns the named zen.Cache registered by this component.
func ResolveNamed(app *zen.App, name string) (zen.Cache, bool) {
	return zen.ResolveAs[zen.Cache](app, NamedService(name))
}

// MustResolveNamed 返回具名 zen.Cache，未注册时 panic。
// MustResolveNamed returns the named zen.Cache or panics.
func MustResolveNamed(app *zen.App, name string) zen.Cache {
	c, ok := ResolveNamed(app, name)
	if !ok {
		panic("redis: named cache service not found: " + name)
	}
	return c
}

// ResolveManager 返回该组件注册的 *Manager。
// ResolveManager returns the *Manager registered by this component.
func ResolveManager(app *zen.App) (*Manager, bool) {
	return zen.ResolveAs[*Manager](app, DefaultManagerServiceName)
}

// NamedService 返回具名 Redis 实例的容器键（如 "cache.session"）。
// NamedService returns the container key for a named Redis instance (e.g. "cache.session").
func NamedService(name string) string {
	return namedService(DefaultServiceName, name)
}

func namedService(prefix, name string) string {
	return strings.TrimSpace(prefix) + "." + strings.TrimSpace(name)
}
