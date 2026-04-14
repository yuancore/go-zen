package zdb

import (
	"errors"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

// Manager holds all opened GORM connections for this app instance.
type Manager struct {
	mu          sync.RWMutex
	defaultName string
	order       []string
	dbs         map[string]*gorm.DB
}

func newManager(defaultName string) *Manager {
	return &Manager{
		defaultName: defaultName,
		dbs:         make(map[string]*gorm.DB),
	}
}

func (m *Manager) register(name string, db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("gorm: nil database for %q", name)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.dbs[name]; exists {
		return fmt.Errorf("gorm: database %q already registered", name)
	}

	m.dbs[name] = db
	m.order = append(m.order, name)
	return nil
}

// DefaultName returns the configured default connection name.
func (m *Manager) DefaultName() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.defaultName
}

// Default returns the default *gorm.DB.
func (m *Manager) Default() (*gorm.DB, bool) {
	return m.Get(m.DefaultName())
}

// MustDefault returns the default *gorm.DB or panics.
func (m *Manager) MustDefault() *gorm.DB {
	db, ok := m.Default()
	if !ok {
		panic("gorm: default database not found")
	}
	return db
}

// Get returns a named *gorm.DB.
func (m *Manager) Get(name string) (*gorm.DB, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	db, ok := m.dbs[name]
	return db, ok
}

// MustGet returns a named *gorm.DB or panics.
func (m *Manager) MustGet(name string) *gorm.DB {
	db, ok := m.Get(name)
	if !ok {
		panic("gorm: database not found: " + name)
	}
	return db
}

// Names returns registered connection names in registration order.
func (m *Manager) Names() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, len(m.order))
	copy(names, m.order)
	return names
}

// Close closes all SQL connections in reverse registration order.
func (m *Manager) Close() error {
	m.mu.RLock()
	names := make([]string, len(m.order))
	copy(names, m.order)
	dbs := make(map[string]*gorm.DB, len(m.dbs))
	for name, db := range m.dbs {
		dbs[name] = db
	}
	m.mu.RUnlock()

	var errs []error
	for i := len(names) - 1; i >= 0; i-- {
		name := names[i]
		db := dbs[name]
		if db == nil {
			continue
		}
		sqlDB, err := db.DB()
		if err != nil {
			errs = append(errs, fmt.Errorf("gorm: get sql.DB for %q: %w", name, err))
			continue
		}
		if err := sqlDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("gorm: close %q: %w", name, err))
		}
	}

	return errors.Join(errs...)
}
