package zgin

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuancore/go-zen/zen"
)

// LoggerConfig controls the behavior of the Logger middleware.
type LoggerConfig struct {
	// SkipPaths are URL paths to exclude from logging.
	// Supports exact match (e.g. "/health") and wildcard prefix (e.g. "/static/*").
	SkipPaths []string
	// SkipMethods are HTTP methods to exclude from logging (case-insensitive).
	SkipMethods []string
	// HeaderWhitelist is the list of request headers to record.
	// ["*"] logs all headers; empty list disables header logging.
	HeaderWhitelist []string
	// EnableRequestBody enables request body capture. Disabled by default.
	// Avoid enabling for endpoints that accept large payloads or file uploads.
	EnableRequestBody bool
	// EnableResponseBody enables response body capture. Disabled by default.
	EnableResponseBody bool
	// MaxBodyLogSize is the maximum bytes included in a logged body field.
	// Default: 8 KB. 0 means unlimited.
	MaxBodyLogSize int64
}

// LoggerOption is a functional option for LoggerConfig.
type LoggerOption func(*LoggerConfig)

// WithSkipPaths adds URL paths to exclude from access logging.
func WithSkipPaths(paths ...string) LoggerOption {
	return func(c *LoggerConfig) { c.SkipPaths = append(c.SkipPaths, paths...) }
}

// WithSkipMethods adds HTTP methods to exclude from access logging.
func WithSkipMethods(methods ...string) LoggerOption {
	return func(c *LoggerConfig) { c.SkipMethods = append(c.SkipMethods, methods...) }
}

// WithRequestBody enables request body capture in log fields.
func WithRequestBody() LoggerOption {
	return func(c *LoggerConfig) { c.EnableRequestBody = true }
}

// WithResponseBody enables response body capture in log fields.
func WithResponseBody() LoggerOption {
	return func(c *LoggerConfig) { c.EnableResponseBody = true }
}

// WithHeaderWhitelist sets the request headers to include in the log entry.
// Pass "*" as the sole element to log all headers.
func WithHeaderWhitelist(headers ...string) LoggerOption {
	return func(c *LoggerConfig) { c.HeaderWhitelist = headers }
}

// WithMaxBodyLogSize overrides the maximum body bytes written to a log entry.
func WithMaxBodyLogSize(n int64) LoggerOption {
	return func(c *LoggerConfig) { c.MaxBodyLogSize = n }
}

// ---------- response body capture ----------

type responseCapture struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseCapture) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

const maxPooledCap = 64 << 10 // discard buffers that grew beyond 64 KB

var respBufPool = sync.Pool{
	New: func() any { return bytes.NewBuffer(make([]byte, 0, 4<<10)) },
}

func acquireRespBuf() *bytes.Buffer { return respBufPool.Get().(*bytes.Buffer) }

func releaseRespBuf(b *bytes.Buffer) {
	if b == nil || b.Cap() > maxPooledCap {
		return
	}
	b.Reset()
	respBufPool.Put(b)
}

// ---------- config-driven constructor ----------

// HttpLogConfig maps to the [http_log] section in the application config file.
//
//	[http_log]
//	enabled          = true
//	request_body     = true
//	response_body    = false
//	skip_paths       = ["/health", "/static/*"]
//	header_whitelist = ["*"]
type HttpLogConfig struct {
	Enabled         bool     `mapstructure:"enabled"`
	RequestBody     bool     `mapstructure:"request_body"`
	ResponseBody    bool     `mapstructure:"response_body"`
	SkipPaths       []string `mapstructure:"skip_paths"`
	HeaderWhitelist []string `mapstructure:"header_whitelist"`
}

// NewLoggerFromConfig builds a NewLogger middleware from the [http_log] config section.
// Returns nil when http_log.enabled is false or the section is absent.
//
//	mw := zgin.NewLoggerFromConfig(cfg, logger)
//	if mw != nil {
//	    engine.Raw().Use(mw)
//	}
func NewLoggerFromConfig(cfg zen.Config, logger zen.Logger) gin.HandlerFunc {
	if cfg == nil || !cfg.GetBool("http_log.enabled") {
		return nil
	}
	var hcfg HttpLogConfig
	if cfg.IsSet("http_log") {
		_ = cfg.Unmarshal("http_log", &hcfg)
	}

	opts := []LoggerOption{WithSkipPaths(hcfg.SkipPaths...)}
	if len(hcfg.HeaderWhitelist) > 0 {
		opts = append(opts, WithHeaderWhitelist(hcfg.HeaderWhitelist...))
	}
	if hcfg.RequestBody {
		opts = append(opts, WithRequestBody())
	}
	if hcfg.ResponseBody {
		opts = append(opts, WithResponseBody())
	}
	return NewLogger(logger, opts...)
}

// LoggerHandler wraps NewLogger as a zen.Handler for explicit middleware registration.
//
//	app.Middleware(zgin.LoggerHandler(logger,
//	    zgin.WithSkipPaths("/health"),
//	    zgin.WithHeaderWhitelist("*"),
//	))
func LoggerHandler(logger zen.Logger, opts ...LoggerOption) zen.Handler {
	h := NewLogger(logger, opts...)
	return func(c zen.Context) {
		if gc, ok := c.(*ginContext); ok {
			h(gc.c)
		} else {
			c.Next()
		}
	}
}

// NewLoggerHandlerFromConfig wraps NewLoggerFromConfig as a zen.Handler.
// Returns nil when http_log.enabled is false.
//
//	if mw := zgin.NewLoggerHandlerFromConfig(cfg, logger); mw != nil {
//	    app.Middleware(mw)
//	}
func NewLoggerHandlerFromConfig(cfg zen.Config, logger zen.Logger) zen.Handler {
	h := NewLoggerFromConfig(cfg, logger)
	if h == nil {
		return nil
	}
	return func(c zen.Context) {
		if gc, ok := c.(*ginContext); ok {
			h(gc.c)
		} else {
			c.Next()
		}
	}
}

// ---------- middleware ----------

// NewLogger returns a production-ready HTTP access log gin middleware.
//
// Logged fields: request_id, method, path, route, status, latency_ms,
// client_ip, resp_bytes. Request/response bodies are opt-in.
//
// Code-driven example:
//
//	engine.Raw().Use(zgin.NewLogger(logger,
//	    zgin.WithSkipPaths("/health", "/metrics"),
//	    zgin.WithSkipMethods("OPTIONS"),
//	))
//
// Config-driven example (reads [http_log] section):
//
//	engine.Raw().Use(zgin.NewLoggerFromConfig(cfg, logger))
func NewLogger(logger zen.Logger, opts ...LoggerOption) gin.HandlerFunc {
	cfg := LoggerConfig{MaxBodyLogSize: 8 << 10}
	for _, opt := range opts {
		opt(&cfg)
	}

	// Pre-compute skip sets for fast per-request lookup.
	exactSkip := make(map[string]struct{}, len(cfg.SkipPaths))
	var prefixSkip []string
	for _, p := range cfg.SkipPaths {
		if strings.HasSuffix(p, "/*") {
			if base := strings.TrimSuffix(p, "/*"); base != "" {
				prefixSkip = append(prefixSkip, base)
			}
		} else {
			exactSkip[p] = struct{}{}
		}
	}

	skipMethods := make(map[string]struct{}, len(cfg.SkipMethods))
	for _, m := range cfg.SkipMethods {
		skipMethods[strings.ToUpper(m)] = struct{}{}
	}

	// Pre-compute header whitelist once — zero alloc per request.
	logAllHeaders := len(cfg.HeaderWhitelist) == 1 && cfg.HeaderWhitelist[0] == "*"
	headerWL := make(map[string]struct{}, len(cfg.HeaderWhitelist))
	for _, h := range cfg.HeaderWhitelist {
		headerWL[http.CanonicalHeaderKey(h)] = struct{}{}
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Fast skip checks — no allocations on skipped requests.
		if _, ok := skipMethods[c.Request.Method]; ok {
			c.Next()
			return
		}
		if _, ok := exactSkip[path]; ok {
			c.Next()
			return
		}
		for _, pfx := range prefixSkip {
			if path == pfx || strings.HasPrefix(path, pfx+"/") {
				c.Next()
				return
			}
		}

		start := time.Now()

		// Optionally snapshot the request body before it is consumed downstream.
		reqBody := snapshotRequestBody(c, cfg.EnableRequestBody, cfg.MaxBodyLogSize)

		// Optionally intercept response writes.
		var capture *bytes.Buffer
		if cfg.EnableResponseBody {
			capture = acquireRespBuf()
			c.Writer = &responseCapture{ResponseWriter: c.Writer, body: capture}
		}

		c.Next()

		status := c.Writer.Status()
		fields := make([]any, 0, 18)
		fields = append(fields,
			"request_id", requestIDFromContext(c),
			"method", c.Request.Method,
			"path", path,
			"route", resolveRoute(c),
			"status", status,
			"latency", time.Since(start).String(),
			"client_ip", c.ClientIP(),
			"resp_bytes", c.Writer.Size(),
		)
		if len(c.Errors) > 0 {
			fields = append(fields, "errors", c.Errors.String())
		}
		if logAllHeaders {
			fields = append(fields, "headers", c.Request.Header)
		} else if len(headerWL) > 0 {
			fields = append(fields, "headers", filteredHeaders(c.Request.Header, headerWL))
		}
		if reqBody != nil {
			fields = append(fields, "request_body", reqBody)
		}
		if capture != nil {
			fields = append(fields, "response_body", parseBodyForLog(capture.Bytes(), cfg.MaxBodyLogSize))
			releaseRespBuf(capture)
		}

		switch {
		case status >= http.StatusInternalServerError:
			logger.Error("http request", fields...)
		case status >= http.StatusBadRequest:
			logger.Warn("http request", fields...)
		default:
			logger.Info("http request", fields...)
		}
	}
}

// snapshotRequestBody reads the full request body, restores it on c.Request.Body
// for downstream handlers, and returns a parsed value ready for structured logging.
func snapshotRequestBody(c *gin.Context, enabled bool, maxSize int64) any {
	if !enabled || c.Request.Body == nil {
		return nil
	}
	data, err := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewReader(data)) // restore unconditionally
	if err != nil || len(data) == 0 {
		return nil
	}
	return parseBodyForLog(data, maxSize)
}

// parseBodyForLog tries to unmarshal data as JSON so zap writes it as a nested
// object instead of an escaped string. Falls back to a plain (truncated) string.
func parseBodyForLog(data []byte, maxSize int64) any {
	if len(data) == 0 {
		return nil
	}
	if maxSize > 0 && int64(len(data)) > maxSize {
		return string(data[:maxSize]) + "...(truncated)"
	}
	var v any
	if err := json.Unmarshal(data, &v); err == nil {
		return v
	}
	return string(data)
}

// filteredHeaders returns a shallow copy of h containing only keys in wl.
// Iterates over the whitelist (typically small) rather than all headers.
func filteredHeaders(h http.Header, wl map[string]struct{}) http.Header {
	out := make(http.Header, len(wl))
	for k := range wl {
		if vals := h[k]; len(vals) > 0 {
			out[k] = vals
		}
	}
	return out
}

// resolveRoute returns the matched route template, falling back to the request path.
func resolveRoute(c *gin.Context) string {
	if route := c.FullPath(); route != "" {
		return route
	}
	return c.Request.URL.Path
}
