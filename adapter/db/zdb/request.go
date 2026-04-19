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

// Scope \u89e3\u6790\u5f53\u524d HTTP \u8bf7\u6c42\u7684\u6570\u636e\u5e93\u53e5\u67c4\u3002
// Scope resolves request-aware database handles for the current HTTP handler.
type Scope struct {
	app *zen.App
	ctx zen.Context
}

// DB \u8fd4\u56de\u7ed1\u5b9a\u5f53\u524d\u8bf7\u6c42 context \u7684\u9ed8\u8ba4\u6570\u636e\u5e93\u3002
// DB returns the default database bound to the current request context.
func DB(app *zen.App, ctx zen.Context) *gorm.DB {
	return For(app, ctx).Default()
}

// Connection \u8fd4\u56de\u7ed1\u5b9a\u5f53\u524d\u8bf7\u6c42 context \u7684\u5177\u540d\u6570\u636e\u5e93\u3002
// Connection returns a named database bound to the current request context.
func Connection(app *zen.App, ctx zen.Context, name string) *gorm.DB {
	return For(app, ctx).Connection(name)
}

// Table \u8fd4\u56de\u9ed8\u8ba4\u6570\u636e\u5e93\u4e0a\u7ed1\u5b9a\u5f53\u524d\u8bf7\u6c42 context \u7684\u8868\u7ea7\u67e5\u8be2\u3002
// Table returns a table-scoped query on the default database with the current
// HTTP request context attached.
func Table(app *zen.App, ctx zen.Context, name string) *gorm.DB {
	return For(app, ctx).Table(name)
}

// TableOn \u8fd4\u56de\u5177\u540d\u6570\u636e\u5e93\u4e0a\u7ed1\u5b9a\u5f53\u524d\u8bf7\u6c42 context \u7684\u8868\u7ea7\u67e5\u8be2\u3002
// TableOn returns a table-scoped query on the named database with the current
// HTTP request context attached.
func TableOn(app *zen.App, ctx zen.Context, connectionName, tableName string) *gorm.DB {
	return For(app, ctx).TableOn(connectionName, tableName)
}

// For \u8fd4\u56de\u81ea\u52a8\u5c06\u5f53\u524d HTTP \u8bf7\u6c42 context \u7ed1\u5b9a\u5230\u6240\u6709\u89e3\u6790\u8fde\u63a5\u7684 Scope\u3002
// For returns a scope that automatically binds the current HTTP request context
// to resolved GORM connections.
func For(app *zen.App, ctx zen.Context) Scope {
	return Scope{app: app, ctx: ctx}
}

// Default \u8fd4\u56de\u7ed1\u5b9a\u5f53\u524d\u8bf7\u6c42 context \u7684\u9ed8\u8ba4\u6570\u636e\u5e93\uff0c\u672a\u6ce8\u518c\u65f6 panic\u3002
// Default returns the default database bound to the current request context.
// It panics if the database service is not registered.
func (s Scope) Default() *gorm.DB {
	return WithRequest(MustResolve(s.app), s.ctx)
}

// DB \u662f Default \u7684\u522b\u540d\u3002
// DB is an alias of Default.
func (s Scope) DB() *gorm.DB { return s.Default() }

// Named \u8fd4\u56de\u7ed1\u5b9a\u5f53\u524d\u8bf7\u6c42 context \u7684\u5177\u540d\u6570\u636e\u5e93\uff0c\u672a\u6ce8\u518c\u65f6 panic\u3002
// Named returns a named database bound to the current request context.
// It panics if the named database service is not registered.
func (s Scope) Named(name string) *gorm.DB {
	return WithRequest(MustResolveNamed(s.app, name), s.ctx)
}

// Connection \u662f Named \u7684 Laravel \u98ce\u683c\u522b\u540d\u3002
// Connection is a Laravel-style alias of Named.
func (s Scope) Connection(name string) *gorm.DB { return s.Named(name) }

// Table \u8fd4\u56de\u9ed8\u8ba4\u6570\u636e\u5e93\u4e0a\u7684\u8868\u7ea7\u67e5\u8be2\u3002
// Table returns a table-scoped query on the default database.
func (s Scope) Table(name string) *gorm.DB { return s.Default().Table(name) }

// TableOn \u8fd4\u56de\u5177\u540d\u6570\u636e\u5e93\u4e0a\u7684\u8868\u7ea7\u67e5\u8be2\u3002
// TableOn returns a table-scoped query on the named database.
func (s Scope) TableOn(connectionName, tableName string) *gorm.DB {
	return s.Connection(connectionName).Table(tableName)
}

// Resolve \u8fd4\u56de\u7ed1\u5b9a\u5f53\u524d\u8bf7\u6c42 context \u7684\u9ed8\u8ba4\u6570\u636e\u5e93\u3002
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

// ResolveNamed \u8fd4\u56de\u7ed1\u5b9a\u5f53\u524d\u8bf7\u6c42 context \u7684\u5177\u540d\u6570\u636e\u5e93\u3002
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

// WithRequest \u5c06\u5f53\u524d HTTP \u8bf7\u6c42 context \u9644\u5230 db \u5e76\u8fd4\u56de\u514b\u9686\u3002
// WithRequest clones db with the current HTTP request context attached.
func WithRequest(db *gorm.DB, ctx zen.Context) *gorm.DB {
	if db == nil || ctx == nil || ctx.Request() == nil {
		return db
	}
	return db.WithContext(ctx.Request().Context())
}

// ---------- Standard context.Context helpers ----------
// \u9002\u7528\u4e8e service / DAO \u5c42\u63a5\u6536\u666e\u901a context.Context \u7684\u573a\u666f
// \uff08\u540e\u53f0\u4efb\u52a1\u3001gRPC handler\u3001CLI \u547d\u4ee4\u7b49\uff09\u3002
// Use these in service / DAO layers that receive a plain context.Context
// (background jobs, gRPC handlers, CLI commands, etc.)

// DBCtx \u4ece context \u4e2d\u8fd4\u56de\u9ed8\u8ba4 *gorm.DB\uff0c\u672a\u627e\u5230\u65f6 panic\u3002
// \u8fd9\u662f\u5728 service \u548c DAO \u5c42\u8bbf\u95ee\u6570\u636e\u5e93\u7684\u63a8\u8350\u65b9\u5f0f\u3002
//
//	func NewOrdersDao(ctx context.Context) *OrdersDao {
//	    return &OrdersDao{db: zdb.DBCtx(ctx)}
//	}
//
// DBCtx returns the default *gorm.DB from context with ctx attached.
// This is the recommended way to access the database in service and DAO layers.
func DBCtx(ctx context.Context) *gorm.DB {
	db, ok := FromContext(ctx)
	if !ok {
		panic("zdb: no database in context; ensure zdb.InjectMiddleware is registered before this handler")
	}
	return db.WithContext(ctx)
}

// ConnCtx \u4ece context \u4e2d\u8fd4\u56de\u5177\u540d *gorm.DB\uff0c\u672a\u627e\u5230\u65f6 panic\u3002
//
//	db := zdb.ConnCtx(ctx, "replica")
//
// ConnCtx returns a named *gorm.DB from context with ctx attached.
func ConnCtx(ctx context.Context, name string) *gorm.DB {
	db, ok := NamedFromContext(ctx, name)
	if !ok {
		panic("zdb: named database " + name + " not in context; ensure zdb.InjectMiddleware is registered")
	}
	return db.WithContext(ctx)
}

// DBFrom \u8fd4\u56de\u5e26\u6709 ctx \u7684\u9ed8\u8ba4 *gorm.DB\uff0ccontext \u4e2d\u6ca1\u6709\u65f6\u56de\u9000\u5230 app \u5bb9\u5668\u3002
// \u5728 handler/service \u4ee3\u7801\u4e2d\u4f18\u5148\u4f7f\u7528 DBCtx\uff1b\u4ec5\u5728 CLI \u6216\u4e00\u6b21\u6027\u811a\u672c\u4e2d\u4f7f\u7528\u6b64\u51fd\u6570\u3002
// DBFrom returns the default *gorm.DB with ctx attached.
// Prefer DBCtx in handler/service code; use DBFrom only in CLI or one-off scripts.
func DBFrom(app *zen.App, ctx context.Context) *gorm.DB {
	if db, ok := FromContext(ctx); ok {
		return db.WithContext(ctx)
	}
	return MustResolve(app).WithContext(ctx)
}

// ConnFrom \u8fd4\u56de\u5e26\u6709 ctx \u7684\u5177\u540d *gorm.DB\uff0c\u56de\u9000\u5230 app \u5bb9\u5668\u3002
// ConnFrom returns a named *gorm.DB with ctx attached, falling back to app.
func ConnFrom(app *zen.App, ctx context.Context, name string) *gorm.DB {
	if db, ok := NamedFromContext(ctx, name); ok {
		return db.WithContext(ctx)
	}
	return MustResolveNamed(app, name).WithContext(ctx)
}

// ---------- context.Context injection / extraction ----------

// InjectDB \u5c06\u9ed8\u8ba4 *gorm.DB \u5b58\u5165 context\u3002
//
//	ctx = zdb.InjectDB(ctx, db)
//
// InjectDB stores the default *gorm.DB in a context so downstream code can
// retrieve it without holding an *App reference.
func InjectDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbContextKey{}, db)
}

// InjectNamedDB \u5c06\u5177\u540d *gorm.DB \u5b58\u5165 context\u3002
//
//	ctx = zdb.InjectNamedDB(ctx, "orders_db", db)
//
// InjectNamedDB stores a named *gorm.DB in a context.
func InjectNamedDB(ctx context.Context, name string, db *gorm.DB) context.Context {
	return context.WithValue(ctx, namedDBContextKey{name}, db)
}

// FromContext \u53d6\u51fa\u7531 InjectDB \u5b58\u5165\u7684\u9ed8\u8ba4 *gorm.DB\u3002
// FromContext retrieves the default *gorm.DB previously stored via InjectDB.
func FromContext(ctx context.Context) (*gorm.DB, bool) {
	db, ok := ctx.Value(dbContextKey{}).(*gorm.DB)
	return db, ok && db != nil
}

// NamedFromContext \u53d6\u51fa\u7531 InjectNamedDB \u5b58\u5165\u7684\u5177\u540d *gorm.DB\u3002
// NamedFromContext retrieves a named *gorm.DB previously stored via InjectNamedDB.
func NamedFromContext(ctx context.Context, name string) (*gorm.DB, bool) {
	db, ok := ctx.Value(namedDBContextKey{name}).(*gorm.DB)
	return db, ok && db != nil
}

// ---------- Middleware helper ----------

// namedDBEntry \u7f13\u5b58\u5355\u4e2a\u5177\u540d\u6570\u636e\u5e93\u5b9e\u4f8b\uff0c\u7528\u4e8e\u4e2d\u95f4\u4ef6\u9884\u6ce8\u5165\u3002
// namedDBEntry caches a single named DB instance for middleware pre-injection.
type namedDBEntry struct {
	name string
	db   *gorm.DB
}

// InjectMiddleware \u8fd4\u56de\u4e00\u4e2a zen.Handler\uff0c\u5c06\u6240\u6709\u5df2\u6ce8\u518c\u7684\u6570\u636e\u5e93\u8fde\u63a5\u9884\u6ce8\u5165\u5230\u6bcf\u6b21\u8bf7\u6c42\u7684 context\u3002
// Manager \u548c\u8fde\u63a5\u5217\u8868\u5728 handler \u521b\u5efa\u65f6\u4e00\u6b21\u6027\u89e3\u6790\uff0c\u4e0d\u5728\u8bf7\u6c42\u70ed\u8def\u5f84\u4e0a\u91cd\u590d\u67e5\u627e\u3002
// \u82e5 manager \u672a\u6ce8\u518c\uff0c\u8bb0\u5f55\u8b66\u544a\u5e76\u8fd4\u56de\u900f\u4f20 handler\u3002
//
// \u901a\u5e38\u65e0\u9700\u624b\u52a8\u8c03\u7528\uff0czdb.New() \u4f1a\u5728\u7ec4\u4ef6 Init \u9636\u6bb5\u81ea\u52a8\u6ce8\u518c\u3002
// \u4ec5\u5728\u9700\u8981\u81ea\u5b9a\u4e49\u6ce8\u518c\u987a\u5e8f\u65f6\u624b\u52a8\u4f7f\u7528\uff1a
//
//	app.Middleware(zdb.InjectMiddleware(app))
//
// InjectMiddleware returns a zen.Handler middleware that pre-injects all registered
// database connections into every request's context. The manager and connection list
// are resolved once at handler creation time — not on every request.
// If the manager is not registered, a warning is logged and a pass-through handler is returned.
//
// Normally you do not need to call this directly; zdb.New() auto-registers it during Init.
func InjectMiddleware(app *zen.App) zen.Handler {
	manager, ok := ResolveManager(app)
	if !ok {
		app.Logger().Warn("zdb: manager not found; DB injection skipped — are [[connections]] configured?")
		return func(c zen.Context) { c.Next() }
	}

	// \u521d\u59cb\u5316\u65f6\u7f13\u5b58 default \u548c\u6240\u6709 named DB\uff0c\u70ed\u8def\u5f84\u65e0\u9501\u3001\u65e0\u5bb9\u5668\u67e5\u627e\u3002
	// Capture default and named DBs once; the hot path is lock-free and allocation-free.
	defDB := manager.MustDefault()
	names := manager.Names()
	entries := make([]namedDBEntry, 0, len(names))
	for _, name := range names {
		entries = append(entries, namedDBEntry{name, manager.MustGet(name)})
	}

	return func(c zen.Context) {
		req := c.Request()
		ctx := req.Context()

		ctx = InjectDB(ctx, defDB)
		for _, e := range entries {
			ctx = InjectNamedDB(ctx, e.name, e.db)
		}

		*req = *req.WithContext(ctx)
		c.Next()
	}
}
