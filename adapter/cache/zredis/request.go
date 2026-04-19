package zredis

import (
	"context"

	"github.com/yuancore/go-zen/zen"
)

// ---------- Context key types ----------

type cacheContextKey struct{}
type namedCacheContextKey struct{ name string }

// ---------- Context injection / extraction ----------

// InjectCache 将默认 zen.Cache 存入 context。
// InjectCache stores the default zen.Cache in a context.
func InjectCache(ctx context.Context, c zen.Cache) context.Context {
	return context.WithValue(ctx, cacheContextKey{}, c)
}

// InjectNamedCache 将具名 zen.Cache 存入 context。
// InjectNamedCache stores a named zen.Cache in a context.
func InjectNamedCache(ctx context.Context, name string, c zen.Cache) context.Context {
	return context.WithValue(ctx, namedCacheContextKey{name}, c)
}

// CacheFromContext 从 context 中取出默认 zen.Cache。
// CacheFromContext retrieves the default zen.Cache from context.
func CacheFromContext(ctx context.Context) (zen.Cache, bool) {
	c, ok := ctx.Value(cacheContextKey{}).(zen.Cache)
	return c, ok && c != nil
}

// NamedCacheFromContext 从 context 中取出具名 zen.Cache。
// NamedCacheFromContext retrieves a named zen.Cache from context.
func NamedCacheFromContext(ctx context.Context, name string) (zen.Cache, bool) {
	c, ok := ctx.Value(namedCacheContextKey{name}).(zen.Cache)
	return c, ok && c != nil
}

// ---------- Middleware ----------

// namedEntry 缓存单个具名缓存实例，用于中间件预注入。
// namedEntry caches a single named cache instance for middleware pre-injection.
type namedEntry struct {
	name  string
	cache zen.Cache
}

// InjectMiddleware 返回一个 zen.Handler，在每次请求前将所有已注册的 Redis 客户端注入到 context。
// Manager 和客户端列表在 handler 创建时一次性解析，不在请求热路径上重复查找，适合高并发场景。
// 若 manager 未注册（如未配置 [[redis]]），记录警告并返回透传 handler。
//
// 通常无需手动调用，zredis.New() 会在组件 Init 阶段自动注册。
// 仅在需要自定义注册顺序时手动使用：
//
//	app.Middleware(zredis.InjectMiddleware(app))
//
// InjectMiddleware returns a zen.Handler that pre-injects all registered Redis clients
// into every request's context. The manager and client list are resolved once at handler
// creation time — not on every request — making it safe for high-concurrency workloads.
// If the manager is not registered (e.g. no [[redis]] config), a warning is logged and
// a pass-through handler is returned.
//
// Normally you do not need to call this directly; zredis.New() auto-registers it during Init.
// Use it manually only when you need custom middleware ordering:
//
//	app.Middleware(zredis.InjectMiddleware(app))
func InjectMiddleware(app *zen.App) zen.Handler {
	mgr, ok := ResolveManager(app)
	if !ok {
		app.Logger().Warn("zredis: manager not found; cache injection skipped — is [[redis]] configured?")
		return func(c zen.Context) { c.Next() }
	}

	// 初始化时缓存 default 和所有 named caches，热路径上无锁、无容器查找。
	// Capture default and named caches once; the hot path is lock-free and allocation-free.
	defCache := mgr.MustDefault()
	names := mgr.Names()
	entries := make([]namedEntry, 0, len(names))
	for _, name := range names {
		if cl, ok := mgr.Get(name); ok {
			entries = append(entries, namedEntry{name, cl})
		}
	}

	return func(c zen.Context) {
		req := c.Request()
		ctx := req.Context()

		ctx = InjectCache(ctx, defCache)
		for _, e := range entries {
			ctx = InjectNamedCache(ctx, e.name, e.cache)
		}

		*req = *req.WithContext(ctx)
		c.Next()
	}
}

// ---------- Context-only accessors (no *zen.App needed) ----------

// CacheCtx 从 context 中取出默认 zen.Cache，未找到时 panic。
// 确保 InjectMiddleware 已注册后方可调用。
//
//	func NewOrdersDao(ctx context.Context) *OrdersDao {
//	    cache := zredis.CacheCtx(ctx)
//	    return &OrdersDao{db: zdb.DBCtx(ctx), cache: cache}
//	}
//
// CacheCtx returns the default zen.Cache from context.
// Panics if no cache is found — ensure InjectMiddleware is registered.
func CacheCtx(ctx context.Context) zen.Cache {
	c, ok := CacheFromContext(ctx)
	if !ok {
		panic("zredis: no cache in context; ensure zredis.InjectMiddleware is registered")
	}
	return c
}

// CacheCtxNamed 从 context 中取出具名 zen.Cache，未找到时 panic。
//
//	session := zredis.CacheCtxNamed(ctx, "session")
//
// CacheCtxNamed returns a named zen.Cache from context.
// Panics if the named cache is not found.
func CacheCtxNamed(ctx context.Context, name string) zen.Cache {
	c, ok := NamedCacheFromContext(ctx, name)
	if !ok {
		panic("zredis: named cache " + name + " not in context; ensure zredis.InjectMiddleware is registered")
	}
	return c
}

// ---------- App-scope accessors (kept for backward compat / CLI use) ----------

// Scope 提供通过 *App 引用访问 Redis 客户端的便捷方式。
// Scope provides access to Redis cache clients from an *App reference.
type Scope struct {
	app *zen.App
}

// For 返回绑定到指定应用的 Scope。
// For returns a Scope bound to the given application.
func For(app *zen.App) Scope {
	return Scope{app: app}
}

// Default 返回默认 zen.Cache，未注册时 panic。
// Default returns the default zen.Cache. Panics if not registered.
func (s Scope) Default() zen.Cache { return MustResolve(s.app) }

// Named 返回指定名称的缓存客户端，未注册时 panic。
// Named returns the cache client for the given instance name. Panics if not registered.
func (s Scope) Named(name string) zen.Cache { return MustResolveNamed(s.app, name) }

// CacheFromApp 从应用容器中返回默认缓存。
// 在 handler/service 代码中优先使用 CacheCtx；此函数仅用于 CLI 或后台任务。
// CacheFromApp returns the default cache from the app container.
// Prefer CacheCtx in handler/service code; use this only in CLI or background tasks.
func CacheFromApp(app *zen.App) zen.Cache {
	return MustResolve(app)
}

// CacheFromAppNamed 从应用容器中返回具名缓存。
// CacheFromAppNamed returns a named cache from the app container.
func CacheFromAppNamed(app *zen.App, name string) zen.Cache {
	return MustResolveNamed(app, name)
}
