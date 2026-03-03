package zen

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const defaultBanner = `
  ______  _____ _   _ 
 |___  / | ____| \ | |
    / /  |  _| |  \| |
   / /   | |___| |\  |
  /____|_|_____|_| \_|  v%s
  High-performance modular framework
`

const version = "1.0.0"

// App is the core application container.
// It manages components lifecycle, HTTP routing, configuration, and DI.
type App struct {
	name        string
	banner      *string
	stopTimeout time.Duration

	engine     Engine
	cfg        Config
	logger     Logger
	ctr        *Container
	components map[string]Component
	order      []Component

	// hooks
	onStart []func() error
	onStop  []func(context.Context) error
}

// New creates a new App with the given options.
//
//	app := zen.New(
//	    zen.Name("myapp"),
//	    zen.WithConfig(cfg),
//	    zen.WithLogger(logger),
//	    zen.WithEngine(eng),
//	    zen.StopTimeout(15*time.Second),
//	)
func New(opts ...Option) *App {
	a := &App{
		name:        "zen-app",
		stopTimeout: 15 * time.Second,
		ctr:         newContainer(),
		components:  make(map[string]Component),
	}
	for _, o := range opts {
		o(a)
	}
	a.init()
	return a
}

// init bootstraps config, logger fallbacks and registers services in the container.
func (a *App) init() {
	// 1. Config — use provided or fallback to empty
	if a.cfg == nil {
		a.cfg = &emptyConfig{}
	}

	// 2. Logger — use provided or fallback to std logger
	if a.logger == nil {
		a.logger = newStdLogger()
	}

	// 3. Engine must be provided via WithEngine
	if a.engine == nil {
		panic("zen: Engine is required — use zen.WithEngine() to provide one")
	}

	// 4. Register self in container
	a.ctr.Provide("app", a)
	a.ctr.Provide("config", a.cfg)
	a.ctr.Provide("logger", a.logger)
}

// ---------- Component Registration ----------

// Use registers one or more components.
// Components are initialized and started in dependency order when Run() is called.
//
//	app.Use(
//	    logger.New(),
//	    db.New(),
//	    cache.New(),
//	)
func (a *App) Use(components ...Component) *App {
	for _, c := range components {
		a.components[c.Name()] = c
	}
	return a
}

// Register is an alias for Use (backward compatibility).
// Deprecated: Use app.Use() instead.
func (a *App) Register(components ...Component) *App {
	return a.Use(components...)
}

// ---------- Routing (delegates to Engine) ----------

// GET registers a handler for GET requests.
func (a *App) GET(path string, h ...Handler) { a.engine.GET(path, h...) }

// POST registers a handler for POST requests.
func (a *App) POST(path string, h ...Handler) { a.engine.POST(path, h...) }

// PUT registers a handler for PUT requests.
func (a *App) PUT(path string, h ...Handler) { a.engine.PUT(path, h...) }

// DELETE registers a handler for DELETE requests.
func (a *App) DELETE(path string, h ...Handler) { a.engine.DELETE(path, h...) }

// PATCH registers a handler for PATCH requests.
func (a *App) PATCH(path string, h ...Handler) { a.engine.PATCH(path, h...) }

// Group creates a new route group with the given prefix and optional middleware.
func (a *App) Group(prefix string, mw ...Handler) RouterGroup {
	return a.engine.Group(prefix, mw...)
}

// Middleware adds global HTTP middleware.
func (a *App) Middleware(mw ...Handler) {
	a.engine.Use(mw...)
}

// ---------- Service Container ----------

// Provide registers a named service in the DI container.
func (a *App) Provide(name string, svc any) { a.ctr.Provide(name, svc) }

// Resolve retrieves a named service from the DI container.
func (a *App) Resolve(name string) (any, bool) { return a.ctr.Resolve(name) }

// MustResolve retrieves a named service and panics if not found.
func (a *App) MustResolve(name string) any { return a.ctr.MustResolve(name) }

// ---------- Accessors ----------

// Config returns the configuration interface.
func (a *App) Config() Config { return a.cfg }

// Logger returns the logger interface.
func (a *App) Logger() Logger { return a.logger }

// Engine returns the underlying Engine for advanced use.
func (a *App) Engine() Engine { return a.engine }

// Name returns the application name.
func (a *App) AppName() string { return a.name }

// ---------- Lifecycle Hooks ----------

// OnStart registers a hook that runs after all components start.
func (a *App) OnStart(fn func() error) {
	a.onStart = append(a.onStart, fn)
}

// OnStop registers a hook that runs before components stop.
func (a *App) OnStop(fn func(context.Context) error) {
	a.onStop = append(a.onStop, fn)
}

// ---------- Backward Compat ----------

// GetConfig is deprecated. Use Config() instead.
func (a *App) GetConfig() Config { return a.cfg }

// GetLogger is deprecated. Use Logger() instead.
func (a *App) GetLogger() Logger { return a.logger }

// ---------- Run & Shutdown ----------

// Run starts the application lifecycle:
//  1. Print banner
//  2. Resolve component dependency order (topological sort)
//  3. Init all components in order
//  4. Start all components in order
//  5. Start HTTP server via Engine
//  6. Wait for signal (SIGINT/SIGTERM)
//  7. Graceful shutdown (stop components in reverse order, then stop server)
func (a *App) Run(addr string) error {
	// Banner
	a.printBanner()

	// Topo sort components
	order, err := topoSort(a.components)
	if err != nil {
		return fmt.Errorf("zen: %w", err)
	}
	a.order = order

	if len(order) > 0 {
		a.logger.Info("zen: components resolved", "order", componentNames(order))
	}

	// Init phase
	for _, c := range order {
		a.logger.Info("zen: init", "component", c.Name())
		if err := c.Init(a); err != nil {
			return fmt.Errorf("zen: component %q init: %w", c.Name(), err)
		}
	}

	// Start phase
	for _, c := range order {
		a.logger.Info("zen: start", "component", c.Name())
		if err := c.Start(); err != nil {
			return fmt.Errorf("zen: component %q start: %w", c.Name(), err)
		}
	}

	// OnStart hooks
	for _, fn := range a.onStart {
		if err := fn(); err != nil {
			return fmt.Errorf("zen: onStart hook: %w", err)
		}
	}

	// Start HTTP server via Engine
	engineErr := make(chan error, 1)
	go func() {
		a.logger.Info("zen: listening", "name", a.name, "addr", addr)
		engineErr <- a.engine.Start(addr)
	}()

	// Wait for signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		a.logger.Info("zen: received signal", "signal", sig.String())
	case err := <-engineErr:
		if err != nil {
			_ = a.shutdown()
			return fmt.Errorf("zen: server: %w", err)
		}
	}

	return a.shutdown()
}

// Stop triggers a graceful shutdown programmatically.
func (a *App) Stop() error {
	return a.shutdown()
}

func (a *App) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), a.stopTimeout)
	defer cancel()

	var first error

	// OnStop hooks (LIFO)
	for i := len(a.onStop) - 1; i >= 0; i-- {
		if err := a.onStop[i](ctx); err != nil {
			a.logger.Error("zen: onStop hook error", "err", err)
			if first == nil {
				first = err
			}
		}
	}

	// Stop components in reverse order
	for i := len(a.order) - 1; i >= 0; i-- {
		c := a.order[i]
		a.logger.Info("zen: stop", "component", c.Name())
		if err := c.Stop(ctx); err != nil {
			a.logger.Error("zen: component stop error", "name", c.Name(), "err", err)
			if first == nil {
				first = err
			}
		}
	}

	// Stop HTTP server via Engine
	if a.engine != nil {
		a.logger.Info("zen: stopping HTTP server")
		if err := a.engine.Stop(ctx); err != nil {
			if first == nil {
				first = err
			}
		}
	}

	a.logger.Info("zen: shutdown complete")
	return first
}

func (a *App) printBanner() {
	if a.banner != nil {
		if *a.banner == "" {
			return // disabled
		}
		fmt.Print(*a.banner)
		return
	}
	fmt.Printf(defaultBanner, version)
	fmt.Printf("  Go: %s | PID: %d\n\n", runtime.Version(), os.Getpid())
}
