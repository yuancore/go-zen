package zdb

import (
	"context"

	"github.com/yuancore/go-zen/zen"
	"gorm.io/gorm"
)

// ---------- Context key types ----------

type dbContextKey struct{}
type namedDBContextKey struct{ name string }

// ---------- HTTP-request-aware helpers (zen.Context) ----------

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
	return Scope{app: app, ctx: ctx}
}

// Default returns the default database bound to the current request context.
// It panics if the database service is not registered.
func (s Scope) Default() *gorm.DB {
	return WithRequest(MustResolve(s.app), s.ctx)
}

// DB is an alias of Default.
func (s Scope) DB() *gorm.DB { return s.Default() }

// Named returns a named database bound to the current request context.
// It panics if the named database service is not registered.
func (s Scope) Named(name string) *gorm.DB {
	return WithRequest(MustResolveNamed(s.app, name), s.ctx)
}

// Connection is a Laravel-style alias of Named.
func (s Scope) Connection(name string) *gorm.DB { return s.Named(name) }

// Table returns a table-scoped query on the default database.
func (s Scope) Table(name string) *gorm.DB { return s.Default().Table(name) }

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

// ---------- Standard context.Context helpers ----------
// Use these in service / DAO layers that receive a plain context.Context
// (background jobs, gRPC handlers, CLI commands, etc.)

// DBCtx returns the default *gorm.DB from context with ctx attached.
// It extracts the DB injected by InjectMiddleware — no *App reference needed.
// This is the recommended way to access the database in service and DAO layers.
//
//	func NewOrdersDao(ctx context.Context) *OrdersDao {
//	    return &OrdersDao{db: zdb.DBCtx(ctx)}
//	}
func DBCtx(ctx context.Context) *gorm.DB {
	db, ok := FromContext(ctx)
	if !ok {
		panic("zdb: no database in context; ensure zdb.InjectMiddleware is registered before this handler")
	}
	return db.WithContext(ctx)
}

// ConnCtx returns a named *gorm.DB from context with ctx attached.
// Use when you need a specific connection (e.g. a read replica or secondary DB).
//
//	db := zdb.ConnCtx(ctx, "replica")
func ConnCtx(ctx context.Context, name string) *gorm.DB {
	db, ok := NamedFromContext(ctx, name)
	if !ok {
		panic("zdb: named database " + name + " not in context; ensure zdb.InjectMiddleware is registered")
	}
	return db.WithContext(ctx)
}

// DBFrom returns the default *gorm.DB with ctx attached.
// Falls back to the app-registered default DB when no DB is found in ctx.
// Prefer DBCtx for handler/service code; use DBFrom only when *App is available
// but injection middleware is not (e.g. CLI commands, one-off scripts).
func DBFrom(app *zen.App, ctx context.Context) *gorm.DB {
	if db, ok := FromContext(ctx); ok {
		return db.WithContext(ctx)
	}
	return MustResolve(app).WithContext(ctx)
}

// ConnFrom returns a named *gorm.DB with ctx attached, falling back to app.
func ConnFrom(app *zen.App, ctx context.Context, name string) *gorm.DB {
	if db, ok := NamedFromContext(ctx, name); ok {
		return db.WithContext(ctx)
	}
	return MustResolveNamed(app, name).WithContext(ctx)
}

// ---------- context.Context injection / extraction ----------

// InjectDB stores the default *gorm.DB in a context so downstream code can
// retrieve it without holding an *App reference.
//
//	ctx = zdb.InjectDB(ctx, db)
func InjectDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbContextKey{}, db)
}

// InjectNamedDB stores a named *gorm.DB in a context.
//
//	ctx = zdb.InjectNamedDB(ctx, "orders_db", db)
func InjectNamedDB(ctx context.Context, name string, db *gorm.DB) context.Context {
	return context.WithValue(ctx, namedDBContextKey{name}, db)
}

// FromContext retrieves the default *gorm.DB previously stored via InjectDB.
func FromContext(ctx context.Context) (*gorm.DB, bool) {
	db, ok := ctx.Value(dbContextKey{}).(*gorm.DB)
	return db, ok && db != nil
}

// NamedFromContext retrieves a named *gorm.DB previously stored via InjectNamedDB.
func NamedFromContext(ctx context.Context, name string) (*gorm.DB, bool) {
	db, ok := ctx.Value(namedDBContextKey{name}).(*gorm.DB)
	return db, ok && db != nil
}

// ---------- Middleware helper ----------

// InjectMiddleware returns a zen.Handler middleware that pre-injects all registered
// database connections into every request's context.
// This lets service/DAO code call zdb.DBCtx(app, ctx.Request().Context()) or
// zdb.FromContext(ctx.Request().Context()) without needing *zen.App references.
//
//	app.Middleware(zdb.InjectMiddleware(app))
func InjectMiddleware(app *zen.App) zen.Handler {
	return func(c zen.Context) {
		manager, ok := ResolveManager(app)
		if !ok {
			c.Next()
			return
		}

		req := c.Request()
		ctx := req.Context()

		// inject default DB
		if def := manager.MustDefault(); def != nil {
			ctx = InjectDB(ctx, def)
		}

		// inject each named DB
		for _, name := range manager.Names() {
			db := manager.MustGet(name)
			ctx = InjectNamedDB(ctx, name, db)
		}

		// update request with enriched context
		*req = *req.WithContext(ctx)
		c.Next()
	}
}
