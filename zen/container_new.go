package zen

import "sync"

// Container is a lightweight, thread-safe service locator / DI container.
// Components register their services here during Init(), and other
// components can Resolve them.
type Container struct {
	mu  sync.RWMutex
	svc map[string]any
}

func newContainer() *Container {
	return &Container{svc: make(map[string]any)}
}

// Provide registers a named service in the container.
func (c *Container) Provide(name string, v any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.svc[name] = v
}

// Resolve retrieves a named service. Returns (nil, false) if not found.
func (c *Container) Resolve(name string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.svc[name]
	return v, ok
}

// MustResolve retrieves a named service and panics if not found.
func (c *Container) MustResolve(name string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.svc[name]
	if !ok {
		panic("zen: service not found: " + name)
	}
	return v
}

// Has checks whether a named service is registered.
func (c *Container) Has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.svc[name]
	return ok
}

// ResolveAs is a generic helper to resolve and type-assert a service.
func ResolveAs[T any](app *App, name string) (T, bool) {
	v, ok := app.Resolve(name)
	if !ok {
		var zero T
		return zero, false
	}
	t, ok := v.(T)
	return t, ok
}
