package zredis

import (
	"errors"
	"fmt"
	"sync"

	"github.com/yuancore/go-zen/zen"
)

// ErrNil 表示请求的 key 在缓存中不存在，调用方应使用 errors.Is 判断：
// ErrNil is returned when a requested key does not exist in the cache.
// Use errors.Is to check for it:
//
//	val, err := cache.Get(ctx, "key")
//	if errors.Is(err, zredis.ErrNil) { ... }
var ErrNil = errors.New("redis: nil")

// Manager 持有当前应用的所有 Redis 客户端实例，设计与 GORM Manager 保持一致。
// Manager holds all opened Redis cache clients for this application instance.
// Its design mirrors the GORM Manager for consistent multi-instance patterns.
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

// DefaultName 返回默认客户端名称。
// DefaultName returns the configured default client name.
func (m *Manager) DefaultName() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.defaultName
}

// Default 返回默认 zen.Cache 及是否找到。
// 使用单次读锁，避免原来 DefaultName()+Get() 的两次加锁开销。
// Default returns the default zen.Cache and whether it was found.
// A single read lock is used to avoid the double-lock overhead of DefaultName()+Get().
func (m *Manager) Default() (zen.Cache, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.clients[m.defaultName]
	return c, ok && c != nil
}

// MustDefault 返回默认 zen.Cache，找不到时 panic。
// MustDefault returns the default zen.Cache or panics if not found.
func (m *Manager) MustDefault() zen.Cache {
	c, ok := m.Default()
	if !ok {
		panic("redis: default client not found: " + m.DefaultName())
	}
	return c
}

// Get 返回指定名称的 zen.Cache 及是否找到。
// Get returns the named zen.Cache and whether it was found.
func (m *Manager) Get(name string) (zen.Cache, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.clients[name]
	return c, ok
}

// MustGet 返回指定名称的 zen.Cache，找不到时 panic。
// MustGet returns the named zen.Cache or panics if not found.
func (m *Manager) MustGet(name string) zen.Cache {
	c, ok := m.Get(name)
	if !ok {
		panic("redis: client not found: " + name)
	}
	return c
}

// Names 按注册顺序返回所有客户端名称。
// Names returns all registered client names in registration order.
func (m *Manager) Names() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, len(m.order))
	copy(names, m.order)
	return names
}

// Close 关闭所有 Redis 客户端连接。
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
