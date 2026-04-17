package zredis

import (
	"errors"
	"fmt"
	"sync"

	"github.com/yuancore/go-zen/zen"
)

// ErrNil is returned when a requested key does not exist in the cache.
// Users should check for it with errors.Is:
//
//	val, err := cache.Get(ctx, "key")
//	if errors.Is(err, zredis.ErrNil) { ... }
var ErrNil = errors.New("redis: nil")

// Manager holds all opened Redis cache clients for this application instance.
// It mirrors the design of the GORM Manager so that multi-instance patterns
// are consistent across the framework.
type Manager struct {
	mu          sync.RWMutex
	defaultName string
	order       []string
	clients     map[string]zen.Cache
}

func newManager(defaultName string) *Manager {
	return &Manager{
		defaultName: defaultName,
		clients:     make(map[string]zen.Cache),
	}
}

func (m *Manager) register(name string, c zen.Cache) error {
	if c == nil {
		return fmt.Errorf("redis: nil client for %q", name)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.clients[name]; exists {
		return fmt.Errorf("redis: client %q already registered", name)
	}

	m.clients[name] = c
	m.order = append(m.order, name)
	return nil
}

// DefaultName returns the configured default client name.
func (m *Manager) DefaultName() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.defaultName
}

// Default returns the default zen.Cache and whether it was found.
func (m *Manager) Default() (zen.Cache, bool) {
	return m.Get(m.DefaultName())
}

// MustDefault returns the default zen.Cache or panics.
func (m *Manager) MustDefault() zen.Cache {
	c, ok := m.Default()
	if !ok {
		panic("redis: default client not found: " + m.DefaultName())
	}
	return c
}

// Get returns the named zen.Cache and whether it was found.
func (m *Manager) Get(name string) (zen.Cache, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.clients[name]
	return c, ok
}

// MustGet returns the named zen.Cache or panics.
func (m *Manager) MustGet(name string) zen.Cache {
	c, ok := m.Get(name)
	if !ok {
		panic("redis: client not found: " + name)
	}
	return c
}

// Names returns all registered client names in registration order.
func (m *Manager) Names() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, len(m.order))
	copy(names, m.order)
	return names
}

// Close closes all Redis clients.
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for _, name := range m.order {
		if c := m.clients[name]; c != nil {
			if err := c.Close(); err != nil {
				errs = append(errs, fmt.Errorf("redis: close %q: %w", name, err))
			}
		}
	}
	return errors.Join(errs...)
}
