package ginadapter

import (
	"context"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"sync/atomic"
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
	logger zen.Logger
}

var _ zen.Engine = (*GinEngine)(nil)

const (
	requestIDHeader = "X-Request-ID"
	traceIDHeader   = "X-Trace-ID"
	spanIDHeader    = "X-Span-ID"
	requestIDKey    = "request_id"
	traceIDKey      = "trace_id"
	spanIDKey       = "span_id"
)

var requestSequence atomic.Uint64

type discardLogger struct{}

func (discardLogger) Debug(string, ...any) {}
func (discardLogger) Info(string, ...any)  {}
func (discardLogger) Warn(string, ...any)  {}
func (discardLogger) Error(string, ...any) {}
func (discardLogger) Fatal(string, ...any) {}
func (discardLogger) With(...any) zen.Logger {
	return discardLogger{}
}

// NewEngine creates a GinEngine with structured request logging and panic recovery.
func NewEngine(logger zen.Logger) *GinEngine {
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	if logger == nil {
		logger = discardLogger{}
	}
	engine := &GinEngine{
		engine: g,
		logger: logger,
	}
	g.Use(engine.requestContextMiddleware(), engine.accessLogMiddleware(), engine.recoveryMiddleware())
	return engine
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

func (e *GinEngine) requestContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		requestID := strings.TrimSpace(c.GetHeader(requestIDHeader))
		if requestID == "" {
			requestID = nextRequestID()
		}
		ctx = context.WithValue(ctx, requestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Set(requestIDKey, requestID)
		c.Header(requestIDHeader, requestID)

		if traceID := strings.TrimSpace(c.GetHeader(traceIDHeader)); traceID != "" {
			ctx = context.WithValue(c.Request.Context(), traceIDKey, traceID)
			c.Request = c.Request.WithContext(ctx)
			c.Set(traceIDKey, traceID)
		}
		if spanID := strings.TrimSpace(c.GetHeader(spanIDHeader)); spanID != "" {
			ctx = context.WithValue(c.Request.Context(), spanIDKey, spanID)
			c.Request = c.Request.WithContext(ctx)
			c.Set(spanIDKey, spanID)
		}

		c.Next()
	}
}

func (e *GinEngine) accessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		fields := []any{
			requestIDKey, requestIDFromContext(c),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"route", route,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"resp_bytes", c.Writer.Size(),
		}
		if len(c.Errors) > 0 {
			fields = append(fields, "errors", c.Errors.String())
		}

		switch status := c.Writer.Status(); {
		case status >= http.StatusInternalServerError:
			e.logger.Error("http request completed", fields...)
		case status >= http.StatusBadRequest:
			e.logger.Warn("http request completed", fields...)
		default:
			e.logger.Info("http request completed", fields...)
		}
	}
}

func (e *GinEngine) recoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		e.logger.Error(
			"http panic recovered",
			requestIDKey, requestIDFromContext(c),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP(),
			"panic", recovered,
			"stack", string(debug.Stack()),
		)
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

func requestIDFromContext(c *gin.Context) string {
	if value, exists := c.Get(requestIDKey); exists {
		if requestID, ok := value.(string); ok && requestID != "" {
			return requestID
		}
	}
	return ""
}

func nextRequestID() string {
	now := strconv.FormatInt(time.Now().UnixNano(), 36)
	seq := strconv.FormatUint(requestSequence.Add(1), 36)
	return now + "-" + seq
}
