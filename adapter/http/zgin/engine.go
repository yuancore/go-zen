package zgin

import (
	"context"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// --- Path / Query params ---
func (g *ginContext) Param(key string) string             { return g.c.Param(key) }
func (g *ginContext) Query(key string) string             { return g.c.Query(key) }
func (g *ginContext) DefaultQuery(key, def string) string { return g.c.DefaultQuery(key, def) }
func (g *ginContext) GetQuery(key string) (string, bool)  { return g.c.GetQuery(key) }
func (g *ginContext) QueryArray(key string) []string      { return g.c.QueryArray(key) }
func (g *ginContext) QueryMap(key string) map[string]string {
	return g.c.QueryMap(key)
}

// --- Request headers ---
func (g *ginContext) Header(key string) string    { return g.c.GetHeader(key) }
func (g *ginContext) SetHeader(key, value string) { g.c.Header(key, value) }
func (g *ginContext) ContentType() string         { return g.c.ContentType() }
func (g *ginContext) GetRawData() ([]byte, error) { return g.c.GetRawData() }
func (g *ginContext) IsWebsocket() bool           { return g.c.IsWebsocket() }

// --- Bind (writes 400 on error) ---
func (g *ginContext) Bind(v any) error       { return g.c.Bind(v) }
func (g *ginContext) BindJSON(v any) error   { return g.c.BindJSON(v) }
func (g *ginContext) BindQuery(v any) error  { return g.c.BindQuery(v) }
func (g *ginContext) BindUri(v any) error    { return g.c.BindUri(v) }
func (g *ginContext) BindHeader(v any) error { return g.c.BindHeader(v) }

// --- ShouldBind (no auto response on error) ---
func (g *ginContext) ShouldBind(v any) error       { return g.c.ShouldBind(v) }
func (g *ginContext) ShouldBindJSON(v any) error   { return g.c.ShouldBindJSON(v) }
func (g *ginContext) ShouldBindQuery(v any) error  { return g.c.ShouldBindQuery(v) }
func (g *ginContext) ShouldBindUri(v any) error    { return g.c.ShouldBindUri(v) }
func (g *ginContext) ShouldBindHeader(v any) error { return g.c.ShouldBindHeader(v) }

// --- Form ---
func (g *ginContext) PostForm(key string) string               { return g.c.PostForm(key) }
func (g *ginContext) DefaultPostForm(key, def string) string   { return g.c.DefaultPostForm(key, def) }
func (g *ginContext) GetPostForm(key string) (string, bool)    { return g.c.GetPostForm(key) }
func (g *ginContext) PostFormArray(key string) []string        { return g.c.PostFormArray(key) }
func (g *ginContext) PostFormMap(key string) map[string]string { return g.c.PostFormMap(key) }

// --- Multipart / File upload ---
func (g *ginContext) FormFile(name string) (*multipart.FileHeader, error) {
	return g.c.FormFile(name)
}
func (g *ginContext) MultipartForm() (*multipart.Form, error) { return g.c.MultipartForm() }
func (g *ginContext) SaveUploadedFile(file *multipart.FileHeader, dst string, perm ...fs.FileMode) error {
	return g.c.SaveUploadedFile(file, dst, perm...)
}

// --- Cookie ---
func (g *ginContext) Cookie(name string) (string, error) { return g.c.Cookie(name) }
func (g *ginContext) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	g.c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}
func (g *ginContext) SetSameSite(samesite http.SameSite) { g.c.SetSameSite(samesite) }

// --- Response: JSON variants ---
func (g *ginContext) JSON(code int, v any)           { g.c.JSON(code, v) }
func (g *ginContext) SecureJSON(code int, v any)     { g.c.SecureJSON(code, v) }
func (g *ginContext) JSONP(code int, obj any)        { g.c.JSONP(code, obj) }
func (g *ginContext) IndentedJSON(code int, obj any) { g.c.IndentedJSON(code, obj) }
func (g *ginContext) AsciiJSON(code int, obj any)    { g.c.AsciiJSON(code, obj) }
func (g *ginContext) PureJSON(code int, obj any)     { g.c.PureJSON(code, obj) }

// --- Response: other ---
func (g *ginContext) XML(code int, obj any)                  { g.c.XML(code, obj) }
func (g *ginContext) YAML(code int, obj any)                 { g.c.YAML(code, obj) }
func (g *ginContext) String(code int, f string, vals ...any) { g.c.String(code, f, vals...) }
func (g *ginContext) Data(code int, ct string, data []byte)  { g.c.Data(code, ct, data) }
func (g *ginContext) Status(code int)                        { g.c.Status(code) }
func (g *ginContext) Redirect(code int, location string)     { g.c.Redirect(code, location) }

// --- Underlying request / response ---
func (g *ginContext) Request() *http.Request              { return g.c.Request }
func (g *ginContext) ResponseWriter() http.ResponseWriter { return g.c.Writer }
func (g *ginContext) ClientIP() string                    { return g.c.ClientIP() }
func (g *ginContext) FullPath() string                    { return g.c.FullPath() }

// --- Context store ---
func (g *ginContext) Set(key string, value any)              { g.c.Set(key, value) }
func (g *ginContext) Get(key string) (any, bool)             { return g.c.Get(key) }
func (g *ginContext) MustGet(key string) any                 { return g.c.MustGet(key) }
func (g *ginContext) GetString(key string) string            { return g.c.GetString(key) }
func (g *ginContext) GetBool(key string) bool                { return g.c.GetBool(key) }
func (g *ginContext) GetInt(key string) int                  { return g.c.GetInt(key) }
func (g *ginContext) GetInt64(key string) int64              { return g.c.GetInt64(key) }
func (g *ginContext) GetFloat64(key string) float64          { return g.c.GetFloat64(key) }
func (g *ginContext) GetStringSlice(key string) []string     { return g.c.GetStringSlice(key) }
func (g *ginContext) GetStringMap(key string) map[string]any { return g.c.GetStringMap(key) }
func (g *ginContext) GetStringMapString(key string) map[string]string {
	return g.c.GetStringMapString(key)
}

// --- Flow control ---
func (g *ginContext) Next()                               { g.c.Next() }
func (g *ginContext) Abort()                              { g.c.Abort() }
func (g *ginContext) AbortWithStatus(code int)            { g.c.AbortWithStatus(code) }
func (g *ginContext) AbortWithStatusJSON(code int, v any) { g.c.AbortWithStatusJSON(code, v) }
func (g *ginContext) IsAborted() bool                     { return g.c.IsAborted() }

// Raw returns the underlying *gin.Context for advanced use.
func (g *ginContext) Raw() *gin.Context { return g.c }

// --- context.Context (allows zen.Context to be passed to context-aware functions) ---
func (g *ginContext) Deadline() (time.Time, bool) { return g.c.Deadline() }
func (g *ginContext) Done() <-chan struct{}       { return g.c.Done() }
func (g *ginContext) Err() error                  { return g.c.Err() }
func (g *ginContext) Value(key any) any           { return g.c.Value(key) }

// --- URL params ---
func (g *ginContext) AddParam(key, value string) { g.c.AddParam(key, value) }

// --- Query string ---
func (g *ginContext) GetQueryArray(key string) ([]string, bool)        { return g.c.GetQueryArray(key) }
func (g *ginContext) GetQueryMap(key string) (map[string]string, bool) { return g.c.GetQueryMap(key) }

// --- Request info ---
func (g *ginContext) RemoteIP() string       { return g.c.RemoteIP() }
func (g *ginContext) HandlerName() string    { return g.c.HandlerName() }
func (g *ginContext) HandlerNames() []string { return g.c.HandlerNames() }

// --- Bind ---
func (g *ginContext) BindXML(v any) error   { return g.c.BindXML(v) }
func (g *ginContext) BindYAML(v any) error  { return g.c.BindYAML(v) }
func (g *ginContext) BindTOML(v any) error  { return g.c.BindTOML(v) }
func (g *ginContext) BindPlain(v any) error { return g.c.BindPlain(v) }

// --- ShouldBind ---
func (g *ginContext) ShouldBindXML(v any) error           { return g.c.ShouldBindXML(v) }
func (g *ginContext) ShouldBindYAML(v any) error          { return g.c.ShouldBindYAML(v) }
func (g *ginContext) ShouldBindTOML(v any) error          { return g.c.ShouldBindTOML(v) }
func (g *ginContext) ShouldBindPlain(v any) error         { return g.c.ShouldBindPlain(v) }
func (g *ginContext) ShouldBindBodyWithJSON(v any) error  { return g.c.ShouldBindBodyWithJSON(v) }
func (g *ginContext) ShouldBindBodyWithXML(v any) error   { return g.c.ShouldBindBodyWithXML(v) }
func (g *ginContext) ShouldBindBodyWithYAML(v any) error  { return g.c.ShouldBindBodyWithYAML(v) }
func (g *ginContext) ShouldBindBodyWithTOML(v any) error  { return g.c.ShouldBindBodyWithTOML(v) }
func (g *ginContext) ShouldBindBodyWithPlain(v any) error { return g.c.ShouldBindBodyWithPlain(v) }

// --- Form ---
func (g *ginContext) GetPostFormArray(key string) ([]string, bool) {
	return g.c.GetPostFormArray(key)
}
func (g *ginContext) GetPostFormMap(key string) (map[string]string, bool) {
	return g.c.GetPostFormMap(key)
}

// --- Cookie ---
func (g *ginContext) SetCookieData(cookie *http.Cookie) { g.c.SetCookieData(cookie) }

// --- Response: additional formats ---
func (g *ginContext) HTML(code int, name string, obj any) { g.c.HTML(code, name, obj) }
func (g *ginContext) TOML(code int, obj any)              { g.c.TOML(code, obj) }
func (g *ginContext) ProtoBuf(code int, obj any)          { g.c.ProtoBuf(code, obj) }
func (g *ginContext) BSON(code int, obj any)              { g.c.BSON(code, obj) }
func (g *ginContext) DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string) {
	g.c.DataFromReader(code, contentLength, contentType, reader, extraHeaders)
}
func (g *ginContext) File(filepath string)                           { g.c.File(filepath) }
func (g *ginContext) FileFromFS(filepath string, fs http.FileSystem) { g.c.FileFromFS(filepath, fs) }
func (g *ginContext) FileAttachment(filepath, filename string) {
	g.c.FileAttachment(filepath, filename)
}
func (g *ginContext) SSEvent(name string, message any)        { g.c.SSEvent(name, message) }
func (g *ginContext) Stream(step func(w io.Writer) bool) bool { return g.c.Stream(step) }
func (g *ginContext) NegotiateFormat(offered ...string) string {
	return g.c.NegotiateFormat(offered...)
}
func (g *ginContext) SetAccepted(formats ...string) { g.c.SetAccepted(formats...) }

// --- Context store: additional types ---
func (g *ginContext) Delete(key string)                    { g.c.Delete(key) }
func (g *ginContext) GetInt8(key string) int8              { return g.c.GetInt8(key) }
func (g *ginContext) GetInt16(key string) int16            { return g.c.GetInt16(key) }
func (g *ginContext) GetInt32(key string) int32            { return g.c.GetInt32(key) }
func (g *ginContext) GetUint(key string) uint              { return g.c.GetUint(key) }
func (g *ginContext) GetUint8(key string) uint8            { return g.c.GetUint8(key) }
func (g *ginContext) GetUint16(key string) uint16          { return g.c.GetUint16(key) }
func (g *ginContext) GetUint32(key string) uint32          { return g.c.GetUint32(key) }
func (g *ginContext) GetUint64(key string) uint64          { return g.c.GetUint64(key) }
func (g *ginContext) GetFloat32(key string) float32        { return g.c.GetFloat32(key) }
func (g *ginContext) GetTime(key string) time.Time         { return g.c.GetTime(key) }
func (g *ginContext) GetDuration(key string) time.Duration { return g.c.GetDuration(key) }
func (g *ginContext) GetStringMapStringSlice(key string) map[string][]string {
	return g.c.GetStringMapStringSlice(key)
}

// --- Error tracking ---
func (g *ginContext) AddError(err error)                 { _ = g.c.Error(err) }
func (g *ginContext) AbortWithError(code int, err error) { _ = g.c.AbortWithError(code, err) }

// --- Abort variants ---
func (g *ginContext) AbortWithStatusPureJSON(code int, obj any) {
	g.c.AbortWithStatusPureJSON(code, obj)
}

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

// GinContext returns the underlying *gin.Context from a zen.Context.
// Use this for gin-specific features not exposed by the zen.Context interface:
// Negotiate, Render, ShouldBindWith, MustBindWith, Handler(), Copy(), etc.
//
//	gc := zgin.GinContext(c)
//	gc.Negotiate(200, gin.Negotiate{Offered: []string{"application/json"}, Data: resp})
func GinContext(c zen.Context) *gin.Context {
	if gc, ok := c.(*ginContext); ok {
		return gc.c
	}
	return nil
}

// ---------- GinEngine ----------

// GinEngine implements zen.Engine backed by gin-gonic/gin.
type GinEngine struct {
	engine        *gin.Engine
	server        *http.Server
	logger        zen.Logger
	serverOptions ServerOptions
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

type discardLogger struct{}

func (discardLogger) Debug(string, ...any) {}
func (discardLogger) Info(string, ...any)  {}
func (discardLogger) Warn(string, ...any)  {}
func (discardLogger) Error(string, ...any) {}
func (discardLogger) Fatal(string, ...any) {}
func (discardLogger) With(...any) zen.Logger {
	return discardLogger{}
}

type ServerOptions struct {
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}

func DefaultServerOptions() ServerOptions {
	return ServerOptions{
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}

// NewEngine creates a GinEngine with structured request logging and panic recovery.
func NewEngine(logger zen.Logger) *GinEngine {
	return NewEngineWithOptions(logger, DefaultServerOptions())
}

// NewEngineWithOptions creates a GinEngine with explicit HTTP server timeouts.
func NewEngineWithOptions(logger zen.Logger, options ServerOptions) *GinEngine {
	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	if logger == nil {
		logger = discardLogger{}
	}
	_ = g.SetTrustedProxies(nil)
	engine := &GinEngine{
		engine:        g,
		logger:        logger,
		serverOptions: normalizeServerOptions(options),
	}
	g.Use(engine.requestContextMiddleware(), engine.accessLogMiddleware(), NewRecovery(engine.logger))
	return engine
}

// --- Routing (implements zen.Engine) ---

func (e *GinEngine) GET(p string, h ...zen.Handler)     { e.engine.GET(p, wrapHandlers(h)...) }
func (e *GinEngine) POST(p string, h ...zen.Handler)    { e.engine.POST(p, wrapHandlers(h)...) }
func (e *GinEngine) PUT(p string, h ...zen.Handler)     { e.engine.PUT(p, wrapHandlers(h)...) }
func (e *GinEngine) DELETE(p string, h ...zen.Handler)  { e.engine.DELETE(p, wrapHandlers(h)...) }
func (e *GinEngine) PATCH(p string, h ...zen.Handler)   { e.engine.PATCH(p, wrapHandlers(h)...) }
func (e *GinEngine) HEAD(p string, h ...zen.Handler)    { e.engine.HEAD(p, wrapHandlers(h)...) }
func (e *GinEngine) OPTIONS(p string, h ...zen.Handler) { e.engine.OPTIONS(p, wrapHandlers(h)...) }
func (e *GinEngine) Use(mw ...zen.Handler)              { e.engine.Use(wrapHandlers(mw)...) }

func (e *GinEngine) Any(p string, h ...zen.Handler) {
	e.engine.Any(p, wrapHandlers(h)...)
}

func (e *GinEngine) Handle(method, p string, h ...zen.Handler) {
	e.engine.Handle(method, p, wrapHandlers(h)...)
}

func (e *GinEngine) StaticFile(relativePath, filepath string) {
	e.engine.StaticFile(relativePath, filepath)
}

func (e *GinEngine) Static(relativePath, root string) {
	e.engine.Static(relativePath, root)
}

func (e *GinEngine) StaticFS(relativePath string, fs http.FileSystem) {
	e.engine.StaticFS(relativePath, fs)
}

func (e *GinEngine) Group(prefix string, mw ...zen.Handler) zen.RouterGroup {
	g := e.engine.Group(prefix, wrapHandlers(mw)...)
	return &ginRouterGroup{group: g}
}

// --- Engine lifecycle ---

func (e *GinEngine) Start(addr string) error {
	e.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	e.server = &http.Server{
		Addr:              addr,
		Handler:           e.engine,
		ReadTimeout:       e.serverOptions.ReadTimeout,
		ReadHeaderTimeout: e.serverOptions.ReadHeaderTimeout,
		WriteTimeout:      e.serverOptions.WriteTimeout,
		IdleTimeout:       e.serverOptions.IdleTimeout,
		MaxHeaderBytes:    e.serverOptions.MaxHeaderBytes,
	}
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

func (g *ginRouterGroup) GET(p string, h ...zen.Handler)     { g.group.GET(p, wrapHandlers(h)...) }
func (g *ginRouterGroup) POST(p string, h ...zen.Handler)    { g.group.POST(p, wrapHandlers(h)...) }
func (g *ginRouterGroup) PUT(p string, h ...zen.Handler)     { g.group.PUT(p, wrapHandlers(h)...) }
func (g *ginRouterGroup) DELETE(p string, h ...zen.Handler)  { g.group.DELETE(p, wrapHandlers(h)...) }
func (g *ginRouterGroup) PATCH(p string, h ...zen.Handler)   { g.group.PATCH(p, wrapHandlers(h)...) }
func (g *ginRouterGroup) HEAD(p string, h ...zen.Handler)    { g.group.HEAD(p, wrapHandlers(h)...) }
func (g *ginRouterGroup) OPTIONS(p string, h ...zen.Handler) { g.group.OPTIONS(p, wrapHandlers(h)...) }
func (g *ginRouterGroup) Use(mw ...zen.Handler)              { g.group.Use(wrapHandlers(mw)...) }

func (g *ginRouterGroup) Any(p string, h ...zen.Handler) {
	g.group.Any(p, wrapHandlers(h)...)
}

func (g *ginRouterGroup) Handle(method, p string, h ...zen.Handler) {
	g.group.Handle(method, p, wrapHandlers(h)...)
}

func (g *ginRouterGroup) StaticFile(relativePath, filepath string) {
	g.group.StaticFile(relativePath, filepath)
}

func (g *ginRouterGroup) Static(relativePath, root string) {
	g.group.Static(relativePath, root)
}

func (g *ginRouterGroup) StaticFS(relativePath string, fs http.FileSystem) {
	g.group.StaticFS(relativePath, fs)
}

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

func requestIDFromContext(c *gin.Context) string {
	if value, exists := c.Get(requestIDKey); exists {
		if requestID, ok := value.(string); ok && requestID != "" {
			return requestID
		}
	}
	return ""
}

func nextRequestID() string {
	return uuid.New().String()
}

func normalizeServerOptions(options ServerOptions) ServerOptions {
	defaults := DefaultServerOptions()

	if options.ReadTimeout <= 0 {
		options.ReadTimeout = defaults.ReadTimeout
	}
	if options.ReadHeaderTimeout <= 0 {
		options.ReadHeaderTimeout = defaults.ReadHeaderTimeout
	}
	if options.WriteTimeout <= 0 {
		options.WriteTimeout = defaults.WriteTimeout
	}
	if options.IdleTimeout <= 0 {
		options.IdleTimeout = defaults.IdleTimeout
	}
	if options.MaxHeaderBytes <= 0 {
		options.MaxHeaderBytes = defaults.MaxHeaderBytes
	}

	return options
}

// NewEngineFromConfig creates a GinEngine whose access log middleware is driven by
// the [http_log] config section instead of the built-in static logger.
// When http_log.enabled is false (or the section is absent) the engine falls back to
// the default accessLogMiddleware so that requests are still always logged.
func NewEngineFromConfig(logger zen.Logger, cfg zen.Config) *GinEngine {
	return NewEngineFromConfigWithOptions(logger, cfg, DefaultServerOptions())
}

// NewEngineFromConfigWithOptions is like NewEngineFromConfig but accepts explicit server timeouts.
func NewEngineFromConfigWithOptions(logger zen.Logger, cfg zen.Config, options ServerOptions) *GinEngine {
	if cfg != nil && cfg.GetBool("system.debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.New()
	if logger == nil {
		logger = discardLogger{}
	}
	_ = g.SetTrustedProxies(nil)
	engine := &GinEngine{
		engine:        g,
		logger:        logger,
		serverOptions: normalizeServerOptions(options),
	}
	accessLog := NewLoggerFromConfig(cfg, logger)
	if accessLog == nil {
		// http_log disabled — use lightweight built-in logger so nothing is lost
		accessLog = engine.accessLogMiddleware()
	}
	g.Use(engine.requestContextMiddleware(), accessLog, NewRecovery(engine.logger))
	return engine
}

// Factory returns a zen.Option that lazily builds a GinEngine using the app logger.
// Use this with zen.New() instead of calling NewEngine manually:
//
//	app := zen.New(
//	    zen.WithConfig(cfg),
//	    zlog.Factory(),
//	    zgin.Factory(),
//	)
func Factory() zen.Option {
	return zen.WithEngineFactory(func(logger zen.Logger) zen.Engine {
		return NewEngine(logger)
	})
}

// FactoryFromConfig returns a zen.Option that lazily builds a GinEngine with access
// logging configured from the [http_log] config section.
// Use this instead of Factory() when you want config-driven request logging:
//
//	cfg, _ := loadConfig(path)
//	app := zen.New(
//	    zen.WithConfig(cfg),
//	    zlog.Factory(),
//	    zgin.FactoryFromConfig(cfg),
//	)
func FactoryFromConfig(cfg zen.Config) zen.Option {
	return zen.WithEngineFactory(func(logger zen.Logger) zen.Engine {
		return NewEngineFromConfig(logger, cfg)
	})
}
