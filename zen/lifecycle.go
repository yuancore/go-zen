package zen

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	engine      Engine
	modules     map[string]Module
	cfg         Config
	logger      Logger
	ctr         *Container
	stopTimeout time.Duration
	order       []Module
}

func New(opts ...Option) *App {
	a := &App{
		modules:     make(map[string]Module),
		ctr:         newContainer(),
		cfg:         &emptyConfig{},
		logger:      newStdLogger(),
		stopTimeout: 15 * time.Second,
	}
	for _, o := range opts {
		o(a)
	}
	return a
}
func (a *App) Register(mods ...Module) *App {
	for _, m := range mods {
		a.modules[m.Name()] = m
	}
	return a
}
func (a *App) GET(p string, h ...Handler)                     { a.engine.GET(p, h...) }
func (a *App) POST(p string, h ...Handler)                    { a.engine.POST(p, h...) }
func (a *App) PUT(p string, h ...Handler)                     { a.engine.PUT(p, h...) }
func (a *App) DELETE(p string, h ...Handler)                  { a.engine.DELETE(p, h...) }
func (a *App) PATCH(p string, h ...Handler)                   { a.engine.PATCH(p, h...) }
func (a *App) Group(prefix string, mw ...Handler) RouterGroup { return a.engine.Group(prefix, mw...) }
func (a *App) Use(mw ...Handler)                              { a.engine.Use(mw...) }
func (a *App) Provide(name string, svc any)                   { a.ctr.Provide(name, svc) }
func (a *App) Resolve(name string) (any, bool)                { return a.ctr.Resolve(name) }
func (a *App) GetEngine() Engine                              { return a.engine }
func (a *App) GetConfig() Config                              { return a.cfg }
func (a *App) GetLogger() Logger                              { return a.logger }
func (a *App) Run(addr string) error {
	if a.engine == nil {
		return fmt.Errorf("zen: no engine; use zen.WithEngine()")
	}
	order, err := topoSort(a.modules)
	if err != nil {
		return fmt.Errorf("zen: %w", err)
	}
	a.order = order
	a.logger.Info("zen: modules resolved", "order", moduleNames(order))
	for _, m := range order {
		a.logger.Info("zen: init", "module", m.Name())
		if err := m.Init(a); err != nil {
			return fmt.Errorf("zen: module %q init: %w", m.Name(), err)
		}
	}
	for _, m := range order {
		a.logger.Info("zen: start", "module", m.Name())
		if err := m.Start(); err != nil {
			return fmt.Errorf("zen: module %q start: %w", m.Name(), err)
		}
	}
	engineErr := make(chan error, 1)
	go func() {
		a.logger.Info("zen: listening", "addr", addr)
		engineErr <- a.engine.Start(addr)
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-quit:
		a.logger.Info("zen: signal", "signal", sig.String())
	case err := <-engineErr:
		if err != nil && err != http.ErrServerClosed {
			_ = a.shutdown()
			return fmt.Errorf("zen: engine: %w", err)
		}
	}
	return a.shutdown()
}
func (a *App) Stop() error { return a.shutdown() }
func (a *App) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), a.stopTimeout)
	defer cancel()
	var first error
	for i := len(a.order) - 1; i >= 0; i-- {
		m := a.order[i]
		a.logger.Info("zen: stop", "module", m.Name())
		if err := m.Stop(ctx); err != nil {
			a.logger.Error("zen: module stop", "name", m.Name(), "err", err)
			if first == nil {
				first = err
			}
		}
	}
	if a.engine != nil {
		a.logger.Info("zen: stopping engine")
		if err := a.engine.Stop(ctx); err != nil {
			if first == nil {
				first = err
			}
		}
	}
	a.logger.Info("zen: shutdown complete")
	return first
}
func moduleNames(mods []Module) []string {
	out := make([]string, len(mods))
	for i, m := range mods {
		out[i] = m.Name()
	}
	return out
}

type stdLogger struct {
	l      *log.Logger
	fields []any
}

func newStdLogger() *stdLogger {
	return &stdLogger{l: log.New(os.Stdout, "[zen] ", log.LstdFlags)}
}
func (s *stdLogger) Debug(msg string, kv ...any) { s.out("DBG", msg, kv) }
func (s *stdLogger) Info(msg string, kv ...any)  { s.out("INF", msg, kv) }
func (s *stdLogger) Warn(msg string, kv ...any)  { s.out("WRN", msg, kv) }
func (s *stdLogger) Error(msg string, kv ...any) { s.out("ERR", msg, kv) }
func (s *stdLogger) Fatal(msg string, kv ...any) { s.out("FTL", msg, kv); os.Exit(1) }
func (s *stdLogger) With(kv ...any) Logger {
	return &stdLogger{l: s.l, fields: append(append([]any{}, s.fields...), kv...)}
}
func (s *stdLogger) out(lvl, msg string, kv []any) {
	all := append(append([]any{}, s.fields...), kv...)
	if len(all) == 0 {
		s.l.Printf("%s %s", lvl, msg)
	} else {
		s.l.Printf("%s %s %v", lvl, msg, all)
	}
}

type emptyConfig struct{}

func (emptyConfig) GetString(string) string        { return "" }
func (emptyConfig) GetInt(string) int              { return 0 }
func (emptyConfig) GetBool(string) bool            { return false }
func (emptyConfig) GetStringSlice(string) []string { return nil }
func (e emptyConfig) Sub(string) Config            { return e }
func (emptyConfig) Unmarshal(string, any) error    { return nil }
