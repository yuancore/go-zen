package ginadapter

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuancore/go-zen/zen"
)

// ---------- ginContext ----------

// ginContext wraps gin.Context to implement zen.Context.
type ginContext struct {
	c *gin.Context
}

func newContext(c *gin.Context) zen.Context {
	return &ginContext{c: c}
}

func (g *ginContext) Param(key string) string                { return g.c.Param(key) }
func (g *ginContext) Query(key string) string                { return g.c.Query(key) }
func (g *ginContext) DefaultQuery(key, def string) string    { return g.c.DefaultQuery(key, def) }
func (g *ginContext) Header(key string) string               { return g.c.GetHeader(key) }
func (g *ginContext) SetHeader(key, value string)            { g.c.Header(key, value) }
func (g *ginContext) BindJSON(v any) error                   { return g.c.ShouldBindJSON(v) }
func (g *ginContext) BindQuery(v any) error                  { return g.c.ShouldBindQuery(v) }
func (g *ginContext) ShouldBind(v any) error                 { return g.c.ShouldBind(v) }
func (g *ginContext) JSON(code int, v any)                   { g.c.JSON(code, v) }
func (g *ginContext) String(code int, f string, vals ...any) { g.c.String(code, f, vals...) }
func (g *ginContext) Data(code int, ct string, data []byte)  { g.c.Data(code, ct, data) }
func (g *ginContext) Status(code int)                        { g.c.Status(code) }
func (g *ginContext) Redirect(code int, location string)     { g.c.Redirect(code, location) }
func (g *ginContext) Request() *http.Request                 { return g.c.Request }
func (g *ginContext) ResponseWriter() http.ResponseWriter    { return g.c.Writer }
func (g *ginContext) ClientIP() string                       { return g.c.ClientIP() }
func (g *ginContext) FullPath() string                       { return g.c.FullPath() }
func (g *ginContext) Set(key string, value any)              { g.c.Set(key, value) }
func (g *ginContext) Get(key string) (any, bool)             { return g.c.Get(key) }
func (g *ginContext) MustGet(key string) any                 { return g.c.MustGet(key) }
func (g *ginContext) Next()                                  { g.c.Next() }
func (g *ginContext) Abort()                                 { g.c.Abort() }
func (g *ginContext) AbortWithStatusJSON(code int, v any)    { g.c.AbortWithStatusJSON(code, v) }
func (g *ginContext) IsAborted() bool                        { return g.c.IsAborted() }

// Raw returns the underlying *gin.Context for advanced use.
func (g *ginContext) Raw() *gin.Context { return g.c }

// ---------- wrapHandler helpers ----------

func wrapHandler(h zen.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(newContext(c))
	}
}

func wrapHandlers(hs []zen.Handler) []gin.HandlerFunc {
	out := make([]gin.HandlerFunc, len(hs))
	for i, h := range hs {
		out[i] = wrapHandler(h)
	}
	return out
}

// ---------- GinEngine ----------

// GinEngine implements zen.Engine backed by gin-gonic/gin.
type GinEngine struct {
	engine *gin.Engine
	server *http.Server
}

var _ zen.Engine = (*GinEngine)(nil)

// NewEngine creates a GinEngine with recovery and logger middleware.
// The logger parameter is accepted for future use and API consistency.
func NewEngine(_ zen.Logger) *GinEngine {
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(gin.Recovery(), gin.Logger())
	return &GinEngine{
		engine: g,
	}
}

// --- Routing (implements zen.Engine) ---

func (e *GinEngine) GET(p string, h ...zen.Handler)    { e.engine.GET(p, wrapHandlers(h)...) }
func (e *GinEngine) POST(p string, h ...zen.Handler)   { e.engine.POST(p, wrapHandlers(h)...) }
func (e *GinEngine) PUT(p string, h ...zen.Handler)    { e.engine.PUT(p, wrapHandlers(h)...) }
func (e *GinEngine) DELETE(p string, h ...zen.Handler) { e.engine.DELETE(p, wrapHandlers(h)...) }
func (e *GinEngine) PATCH(p string, h ...zen.Handler)  { e.engine.PATCH(p, wrapHandlers(h)...) }
func (e *GinEngine) Use(mw ...zen.Handler)             { e.engine.Use(wrapHandlers(mw)...) }

func (e *GinEngine) Group(prefix string, mw ...zen.Handler) zen.RouterGroup {
	g := e.engine.Group(prefix, wrapHandlers(mw)...)
	return &ginRouterGroup{group: g}
}

// --- Engine lifecycle ---

func (e *GinEngine) Start(addr string) error {
	e.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	e.server = &http.Server{Addr: addr, Handler: e.engine}
	return e.server.ListenAndServe()
}

func (e *GinEngine) Stop(ctx context.Context) error {
	if e.server == nil {
		return nil
	}
	shutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return e.server.Shutdown(shutCtx)
}

// Raw returns the underlying *gin.Engine for advanced use.
func (e *GinEngine) Raw() *gin.Engine { return e.engine }

// --- ginRouterGroup ---

type ginRouterGroup struct {
	group *gin.RouterGroup
}

func (g *ginRouterGroup) GET(p string, h ...zen.Handler)    { g.group.GET(p, wrapHandlers(h)...) }
func (g *ginRouterGroup) POST(p string, h ...zen.Handler)   { g.group.POST(p, wrapHandlers(h)...) }
func (g *ginRouterGroup) PUT(p string, h ...zen.Handler)    { g.group.PUT(p, wrapHandlers(h)...) }
func (g *ginRouterGroup) DELETE(p string, h ...zen.Handler) { g.group.DELETE(p, wrapHandlers(h)...) }
func (g *ginRouterGroup) PATCH(p string, h ...zen.Handler)  { g.group.PATCH(p, wrapHandlers(h)...) }
func (g *ginRouterGroup) Use(mw ...zen.Handler)             { g.group.Use(wrapHandlers(mw)...) }

func (g *ginRouterGroup) Group(prefix string, mw ...zen.Handler) zen.RouterGroup {
	sub := g.group.Group(prefix, wrapHandlers(mw)...)
	return &ginRouterGroup{group: sub}
}
