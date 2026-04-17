package zen

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
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
	order      []Component   // flattened init order, used for reverse shutdown
	levels     [][]Component // grouped by depth for parallel init/start

	// lazy factories — called in init() when the concrete impl is not directly provided
	loggerFactory func(Config) (Logger, error)
	engineFactory func(Logger) Engine

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

	// 2. Logger — use provided, or build lazily from factory, or fallback to std logger
	if a.logger == nil && a.loggerFactory != nil {
		l, err := a.loggerFactory(a.cfg)
		if err != nil {
			panic("zen: logger factory: " + err.Error())
		}
		a.logger = l
	}
	if a.logger == nil {
		a.logger = newStdLogger()
	}

	// 3. Engine — use provided, or build lazily from factory
	if a.engine == nil && a.engineFactory != nil {
		a.engine = a.engineFactory(a.logger)
	}
	if a.engine == nil {
		panic("zen: Engine is required — use zen.WithEngine() or zen.WithEngineFactory() to provide one")
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

// HEAD registers a handler for HEAD requests.
func (a *App) HEAD(path string, h ...Handler) { a.engine.HEAD(path, h...) }

// OPTIONS registers a handler for OPTIONS requests.
func (a *App) OPTIONS(path string, h ...Handler) { a.engine.OPTIONS(path, h...) }

// Any registers a handler for all HTTP methods.
func (a *App) Any(path string, h ...Handler) { a.engine.Any(path, h...) }

// Handle registers a handler for the specified HTTP method.
func (a *App) Handle(method, path string, h ...Handler) { a.engine.Handle(method, path, h...) }

// Group creates a new route group with the given prefix and optional middleware.
func (a *App) Group(prefix string, mw ...Handler) RouterGroup {
	return a.engine.Group(prefix, mw...)
}

// Middleware adds global HTTP middleware.
func (a *App) Middleware(mw ...Handler) {
	a.engine.Use(mw...)
}

// StaticFile serves a single file at the given path.
func (a *App) StaticFile(relativePath, filepath string) {
	a.engine.StaticFile(relativePath, filepath)
}

// Static serves files from the given root directory.
func (a *App) Static(relativePath, root string) {
	a.engine.Static(relativePath, root)
}

// StaticFS serves files from the given file system.
func (a *App) StaticFS(relativePath string, fs http.FileSystem) {
	a.engine.StaticFS(relativePath, fs)
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

// RegisterCache registers a custom Cache implementation as the default cache.
// Call this to replace the built-in go-redis adapter with your own backend
// (e.g. Memcached, in-memory, mock):
//
//	app.RegisterCache(myMemcachedCache)
//
// RegisterCache is equivalent to app.Provide(zen.DefaultCacheServiceName, cache)
// but is self-documenting and type-safe.
func (a *App) RegisterCache(cache Cache) {
	a.ctr.Provide(DefaultCacheServiceName, cache)
}

// Cache returns the default cache registered in the container.
// Returns (nil, false) if no cache has been registered yet.
func (a *App) Cache() (Cache, bool) {
	return ResolveAs[Cache](a, DefaultCacheServiceName)
}

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
//  2. Resolve component dependency order (topological sort by levels)
//  3. Init all components (parallel within each level)
//  4. Start all components (parallel within each level)
//  5. Start HTTP server via Engine
//  6. Wait for signal (SIGINT/SIGTERM)
//  7. Graceful shutdown (stop components in reverse order, then stop server)
func (a *App) Run(addr string) error {
	// Banner
	a.printBanner()

	// Topo sort into levels for parallel init/start
	levels, err := topoLevels(a.components)
	if err != nil {
		return fmt.Errorf("zen: %w", err)
	}
	a.levels = levels

	// Flatten levels into a.order for reverse-order shutdown
	a.order = make([]Component, 0, len(a.components))
	for _, level := range levels {
		a.order = append(a.order, level...)
	}

	if len(a.order) > 0 {
		a.logger.Info("zen: components resolved", "order", componentNames(a.order))
	}

	// Init phase — parallel within each level.
	// On failure, stop any components that were already initialised so
	// connections/goroutines are not leaked (all Stop() impls guard with nil
	// checks, so calling them on partially-initialised components is safe).
	for _, level := range levels {
		if err := a.initLevel(level); err != nil {
			stopCtx, stopCancel := context.WithTimeout(context.Background(), a.stopTimeout)
			_ = a.stopComponents(stopCtx, a.order)
			stopCancel()
			return err
		}
	}

	// Start phase — parallel within each level.
	for _, level := range levels {
		if err := a.startLevel(level); err != nil {
			stopCtx, stopCancel := context.WithTimeout(context.Background(), a.stopTimeout)
			_ = a.stopComponents(stopCtx, a.order)
			stopCancel()
			return err
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

// stopComponents stops components in reverse registration order.
// It is safe to call even when components are partially initialised —
// every Stop() implementation guards access with a nil-check on its manager.
func (a *App) stopComponents(ctx context.Context, components []Component) error {
	var first error
	for i := len(components) - 1; i >= 0; i-- {
		c := components[i]
		if err := c.Stop(ctx); err != nil {
			a.logger.Error("zen: component stop error", "name", c.Name(), "err", err)
			if first == nil {
				first = err
			}
		}
	}
	return first
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

// ---------- Parallel helpers ----------

// initLevel initializes all components in a level concurrently.
// Components at the same level have no inter-dependencies, so parallel execution is safe.
func (a *App) initLevel(level []Component) error {
	if len(level) == 1 {
		c := level[0]
		a.logger.Info("zen: init", "component", c.Name())
		if err := c.Init(a); err != nil {
			return fmt.Errorf("zen: component %q init: %w", c.Name(), err)
		}
		return nil
	}
	a.logger.Info("zen: init level", "components", componentNames(level))
	errs := make([]error, len(level))
	var wg sync.WaitGroup
	for i, c := range level {
		wg.Add(1)
		go func(i int, c Component) {
			defer wg.Done()
			errs[i] = c.Init(a)
		}(i, c)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			return fmt.Errorf("zen: component %q init: %w", level[i].Name(), err)
		}
	}
	return nil
}

// startLevel starts all components in a level concurrently.
func (a *App) startLevel(level []Component) error {
	if len(level) == 1 {
		c := level[0]
		a.logger.Info("zen: start", "component", c.Name())
		if err := c.Start(); err != nil {
			return fmt.Errorf("zen: component %q start: %w", c.Name(), err)
		}
		return nil
	}
	a.logger.Info("zen: start level", "components", componentNames(level))
	errs := make([]error, len(level))
	var wg sync.WaitGroup
	for i, c := range level {
		wg.Add(1)
		go func(i int, c Component) {
			defer wg.Done()
			errs[i] = c.Start()
		}(i, c)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			return fmt.Errorf("zen: component %q start: %w", level[i].Name(), err)
		}
	}
	return nil
}
