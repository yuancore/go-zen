package zen

import "sync"

// Container is a lightweight, thread-safe service locator.
type Container struct {
	mu  sync.RWMutex
	svc map[string]any
}

func newContainer() *Container {
	return &Container{svc: make(map[string]any)}
}

// Provide registers a named service.
func (c *Container) Provide(name string, v any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.svc[name] = v
}

// Resolve retrieves a named service.
func (c *Container) Resolve(name string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.svc[name]
	return v, ok
}
