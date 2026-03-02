package ginadapter

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuancore/go-zen/zen"
)

// GinEngine implements zen.Engine backed by gin-gonic/gin.
type GinEngine struct {
	engine       *gin.Engine
	server       *http.Server
	healthChecks map[string]func() error
}

var _ zen.Engine = (*GinEngine)(nil)

// NewEngine creates a GinEngine with recovery and logger middleware.
func NewEngine(logger zen.Logger) *GinEngine {
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(gin.Recovery(), gin.Logger())
	return &GinEngine{
		engine:       g,
		healthChecks: make(map[string]func() error),
	}
}

// --- RouterGroup (top-level) ---

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

func (e *GinEngine) AddHealthCheck(name string, check func() error) {
	e.healthChecks[name] = check
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
