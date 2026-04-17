package zredis

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/yuancore/go-zen/zen"
)

// Option customizes the Redis component.
type Option func(*Component)

// Component manages Redis connections and injects them into the zen container.
//
// Usage:
//
//	app.Use(zredis.New())                   // reads [[redis]] from config
//	app.Use(zredis.New(
//	    zredis.WithConfigKey("my_redis"),   // custom config key
//	))
type Component struct {
	zen.BaseComponent

	configKey          string
	serviceName        string
	managerServiceName string
	instances          []InstanceConfig // explicit override; nil = read from config

	manager   *Manager
	closeOnce sync.Once
}

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

// WithConfigKey overrides the config section key (default: "redis").
func WithConfigKey(key string) Option {
	return func(c *Component) {
		if trimmed := strings.TrimSpace(key); trimmed != "" {
			c.configKey = trimmed
		}
	}
}

// WithInstances injects explicit instance configurations.
// Useful for tests or fully code-driven setups.
//
//	app.Use(zredis.New(
//	    zredis.WithInstances(zredis.InstanceConfig{
//	        Name:    "default",
//	        Address: "127.0.0.1:6379",
//	    }),
//	))
func WithInstances(instances ...InstanceConfig) Option {
	return func(c *Component) {
		c.instances = append(c.instances, instances...)
	}
}

// WithServiceName overrides the container key for the default zen.Cache (default: "cache").
func WithServiceName(name string) Option {
	return func(c *Component) {
		if trimmed := strings.TrimSpace(name); trimmed != "" {
			c.serviceName = trimmed
		}
	}
}

// WithManagerServiceName overrides the container key for the *Manager (default: "redis_manager").
func WithManagerServiceName(name string) Option {
	return func(c *Component) {
		if trimmed := strings.TrimSpace(name); trimmed != "" {
			c.managerServiceName = trimmed
		}
	}
}

// ---------- Component lifecycle ----------

// Init opens Redis connections and registers them in the app container.
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
		name := strings.TrimSpace(inst.Name)
		if name == "" {
			return fmt.Errorf("redis: every instance must have a non-empty name")
		}
		client := newClient(inst)

		// Startup connectivity check — mirrors the GORM ping behaviour.
		// Fails fast so the process exits instead of silently running without cache.
		if timeout := inst.effectivePingTimeout(); timeout > 0 {
			pingCtx, cancel := context.WithTimeout(context.Background(), timeout)
			pingErr := client.Ping(pingCtx)
			cancel()
			if pingErr != nil {
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

	// Register manager and individual clients in the container.
	app.Provide(c.managerServiceName, mgr)
	for _, name := range mgr.Names() {
		cl, _ := mgr.Get(name)
		app.Provide(namedService(c.serviceName, name), cl)
	}
	// Register the default instance under the plain service name.
	app.Provide(c.serviceName, mgr.MustDefault())

	logger.Info("redis initialized",
		"default", mgr.DefaultName(),
		"instances", mgr.Names(),
	)
	return nil
}

// Start is a no-op; connections are ready after Init.
func (c *Component) Start() error { return nil }

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

// Resolve returns the default zen.Cache registered by this component.
func Resolve(app *zen.App) (zen.Cache, bool) {
	return zen.ResolveAs[zen.Cache](app, DefaultServiceName)
}

// MustResolve returns the default zen.Cache or panics.
func MustResolve(app *zen.App) zen.Cache {
	c, ok := Resolve(app)
	if !ok {
		panic("redis: default cache service not found")
	}
	return c
}

// ResolveNamed returns the named zen.Cache registered by this component.
func ResolveNamed(app *zen.App, name string) (zen.Cache, bool) {
	return zen.ResolveAs[zen.Cache](app, NamedService(name))
}

// MustResolveNamed returns the named zen.Cache or panics.
func MustResolveNamed(app *zen.App, name string) zen.Cache {
	c, ok := ResolveNamed(app, name)
	if !ok {
		panic("redis: named cache service not found: " + name)
	}
	return c
}

// ResolveManager returns the *Manager registered by this component.
func ResolveManager(app *zen.App) (*Manager, bool) {
	return zen.ResolveAs[*Manager](app, DefaultManagerServiceName)
}

// NamedService returns the container key for a named Redis instance
// (e.g. "cache.session").
func NamedService(name string) string {
	return namedService(DefaultServiceName, name)
}

func namedService(prefix, name string) string {
	return strings.TrimSpace(prefix) + "." + strings.TrimSpace(name)
}
