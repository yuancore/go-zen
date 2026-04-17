package zredis

import (
	"context"

	"github.com/yuancore/go-zen/zen"
)

// ---------- Context key types ----------

type cacheContextKey struct{}
type namedCacheContextKey struct{ name string }

// ---------- Context injection / extraction ----------

// InjectCache stores the default zen.Cache in a context.
func InjectCache(ctx context.Context, c zen.Cache) context.Context {
	return context.WithValue(ctx, cacheContextKey{}, c)
}

// InjectNamedCache stores a named zen.Cache in a context.
func InjectNamedCache(ctx context.Context, name string, c zen.Cache) context.Context {
	return context.WithValue(ctx, namedCacheContextKey{name}, c)
}

// CacheFromContext retrieves the default zen.Cache from context.
func CacheFromContext(ctx context.Context) (zen.Cache, bool) {
	c, ok := ctx.Value(cacheContextKey{}).(zen.Cache)
	return c, ok && c != nil
}

// NamedCacheFromContext retrieves a named zen.Cache from context.
func NamedCacheFromContext(ctx context.Context, name string) (zen.Cache, bool) {
	c, ok := ctx.Value(namedCacheContextKey{name}).(zen.Cache)
	return c, ok && c != nil
}

// ---------- Middleware ----------

// InjectMiddleware returns a zen.Handler that pre-injects all registered
// Redis clients into every request's context. Register it alongside
// zdb.InjectMiddleware in your app bootstrap:
//
//	app.OnStart(func() error {
//	    app.Middleware(zdb.InjectMiddleware(app), zredis.InjectMiddleware(app))
//	    return router.Setup(app)
//	})
func InjectMiddleware(app *zen.App) zen.Handler {
	return func(c zen.Context) {
		mgr, ok := ResolveManager(app)
		if !ok {
			c.Next()
			return
		}

		req := c.Request()
		ctx := req.Context()

		// inject default cache
		if def := mgr.MustDefault(); def != nil {
			ctx = InjectCache(ctx, def)
		}

		// inject each named cache
		for _, name := range mgr.Names() {
			cl, _ := mgr.Get(name)
			ctx = InjectNamedCache(ctx, name, cl)
		}

		*req = *req.WithContext(ctx)
		c.Next()
	}
}

// ---------- Context-only accessors (no *zen.App needed) ----------

// CacheCtx returns the default zen.Cache from context.
// Panics if no cache is found — ensure InjectMiddleware is registered.
//
//	func NewOrdersDao(ctx context.Context) *OrdersDao {
//	    cache := zredis.CacheCtx(ctx)
//	    return &OrdersDao{db: zdb.DBCtx(ctx), cache: cache}
//	}
func CacheCtx(ctx context.Context) zen.Cache {
	c, ok := CacheFromContext(ctx)
	if !ok {
		panic("zredis: no cache in context; ensure zredis.InjectMiddleware is registered")
	}
	return c
}

// CacheCtxNamed returns a named zen.Cache from context.
// Panics if the named cache is not found.
//
//	session := zredis.CacheCtxNamed(ctx, "session")
func CacheCtxNamed(ctx context.Context, name string) zen.Cache {
	c, ok := NamedCacheFromContext(ctx, name)
	if !ok {
		panic("zredis: named cache " + name + " not in context; ensure zredis.InjectMiddleware is registered")
	}
	return c
}

// ---------- App-scope accessors (kept for backward compat / CLI use) ----------

// Scope provides access to Redis cache clients from an *App reference.
type Scope struct {
	app *zen.App
}

// For returns a Scope bound to the given application.
func For(app *zen.App) Scope {
	return Scope{app: app}
}

// Default returns the default zen.Cache. Panics if not registered.
func (s Scope) Default() zen.Cache { return MustResolve(s.app) }

// Named returns the cache client for the given instance name. Panics if not registered.
func (s Scope) Named(name string) zen.Cache { return MustResolveNamed(s.app, name) }

// CacheFromApp returns the default cache from the app container.
// Prefer CacheCtx in handler/service code. Use this only in CLI or
// background tasks where InjectMiddleware is not involved.
func CacheFromApp(app *zen.App) zen.Cache {
	return MustResolve(app)
}

// CacheFromAppNamed returns a named cache from the app container.
func CacheFromAppNamed(app *zen.App, name string) zen.Cache {
	return MustResolveNamed(app, name)
}
