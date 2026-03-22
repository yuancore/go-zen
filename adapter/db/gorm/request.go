package gormadapter

import (
	"github.com/yuancore/go-zen/zen"
	"gorm.io/gorm"
)

// Scope resolves request-aware database handles for the current HTTP handler.
type Scope struct {
	app *zen.App
	ctx zen.Context
}

// DB returns the default database bound to the current request context.
func DB(app *zen.App, ctx zen.Context) *gorm.DB {
	return For(app, ctx).Default()
}

// Connection returns a named database bound to the current request context.
func Connection(app *zen.App, ctx zen.Context, name string) *gorm.DB {
	return For(app, ctx).Connection(name)
}

// Table returns a table-scoped query on the default database with the current
// HTTP request context attached.
func Table(app *zen.App, ctx zen.Context, name string) *gorm.DB {
	return For(app, ctx).Table(name)
}

// TableOn returns a table-scoped query on the named database with the current
// HTTP request context attached.
func TableOn(app *zen.App, ctx zen.Context, connectionName, tableName string) *gorm.DB {
	return For(app, ctx).TableOn(connectionName, tableName)
}

// For returns a scope that automatically binds the current HTTP request context
// to resolved GORM connections.
func For(app *zen.App, ctx zen.Context) Scope {
	return Scope{
		app: app,
		ctx: ctx,
	}
}

// Default returns the default database bound to the current request context.
// It panics if the database service is not registered.
func (s Scope) Default() *gorm.DB {
	return WithRequest(MustResolve(s.app), s.ctx)
}

// DB is an alias of Default.
func (s Scope) DB() *gorm.DB {
	return s.Default()
}

// Named returns a named database bound to the current request context.
// It panics if the named database service is not registered.
func (s Scope) Named(name string) *gorm.DB {
	return WithRequest(MustResolveNamed(s.app, name), s.ctx)
}

// Connection is a Laravel-style alias of Named.
func (s Scope) Connection(name string) *gorm.DB {
	return s.Named(name)
}

// Table returns a table-scoped query on the default database.
func (s Scope) Table(name string) *gorm.DB {
	return s.Default().Table(name)
}

// TableOn returns a table-scoped query on the named database.
func (s Scope) TableOn(connectionName, tableName string) *gorm.DB {
	return s.Connection(connectionName).Table(tableName)
}

// Resolve returns the default database bound to the current request context.
func (s Scope) Resolve() (*gorm.DB, bool) {
	if s.app == nil {
		return nil, false
	}
	db, ok := Resolve(s.app)
	if !ok {
		return nil, false
	}
	return WithRequest(db, s.ctx), true
}

// ResolveNamed returns a named database bound to the current request context.
func (s Scope) ResolveNamed(name string) (*gorm.DB, bool) {
	if s.app == nil {
		return nil, false
	}
	db, ok := ResolveNamed(s.app, name)
	if !ok {
		return nil, false
	}
	return WithRequest(db, s.ctx), true
}

// WithRequest clones db with the current HTTP request context attached.
func WithRequest(db *gorm.DB, ctx zen.Context) *gorm.DB {
	if db == nil || ctx == nil || ctx.Request() == nil {
		return db
	}
	return db.WithContext(ctx.Request().Context())
}
