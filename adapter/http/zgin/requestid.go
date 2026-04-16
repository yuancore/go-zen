package zgin

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuancore/go-zen/zen"
)

// Generator is a function that produces a unique request ID string.
type Generator func() string

// RequestIDConfig holds options for the RequestID middleware.
type RequestIDConfig struct {
	// HeaderKey is the HTTP header name used to carry the request ID.
	// Default: "X-Request-ID"
	HeaderKey string
	// Generator produces a new ID when the incoming request has none.
	// Default: nextRequestID() (nanosecond timestamp + atomic sequence, no external dep)
	Generator Generator
	// Handler is an optional callback invoked after the ID is resolved.
	// Use it to store the ID in a custom location (e.g. logger context).
	Handler func(c *gin.Context, id string)
}

// RequestIDOption is a functional option for RequestIDConfig.
type RequestIDOption func(*RequestIDConfig)

// WithIDGenerator overrides the default ID generator.
//
//	import "github.com/google/uuid"
//	zgin.WithIDGenerator(uuid.NewString)
func WithIDGenerator(g Generator) RequestIDOption {
	return func(c *RequestIDConfig) { c.Generator = g }
}

// WithIDHeader overrides the header key (default "X-Request-ID").
func WithIDHeader(key string) RequestIDOption {
	return func(c *RequestIDConfig) { c.HeaderKey = key }
}

// WithIDHandler registers a callback invoked after ID resolution.
func WithIDHandler(h func(c *gin.Context, id string)) RequestIDOption {
	return func(c *RequestIDConfig) { c.Handler = h }
}

// NewRequestID returns a gin middleware that ensures every request carries a
// unique identifier. It reads the ID from the incoming header (if present),
// falls back to Generator, writes it back to both the request and response
// headers, and stores it in the gin/go context under the "request_id" key.
//
// This middleware is already embedded in the engine via requestContextMiddleware.
// Use NewRequestID only when you need a custom generator or a handler callback.
//
//	engine.Raw().Use(zgin.NewRequestID(
//	    zgin.WithIDGenerator(uuid.NewString),
//	    zgin.WithIDHeader("X-Request-ID"),
//	))
func NewRequestID(opts ...RequestIDOption) gin.HandlerFunc {
	cfg := RequestIDConfig{
		HeaderKey: requestIDHeader, // "X-Request-ID" from engine.go constants
		Generator: nextRequestID,
	}
	for _, o := range opts {
		o(&cfg)
	}
	headerKey := cfg.HeaderKey

	return func(c *gin.Context) {
		id := strings.TrimSpace(c.GetHeader(headerKey))
		if id == "" {
			id = cfg.Generator()
			c.Request.Header.Set(headerKey, id)
		}
		if cfg.Handler != nil {
			cfg.Handler(c, id)
		}
		// Propagate to gin context, Go context, and response header.
		c.Set(requestIDKey, id)
		ctx := context.WithValue(c.Request.Context(), requestIDKey, id)
		c.Request = c.Request.WithContext(ctx)
		c.Header(headerKey, id)
		c.Next()
	}
}

// RequestIDHandler wraps NewRequestID as a zen.Handler for app.Middleware().
//
//	app.Middleware(zgin.RequestIDHandler(zgin.WithIDGenerator(uuid.NewString)))
func RequestIDHandler(opts ...RequestIDOption) zen.Handler {
	h := NewRequestID(opts...)
	return func(c zen.Context) {
		if gc, ok := c.(*ginContext); ok {
			h(gc.c)
		} else {
			c.Next()
		}
	}
}

// GetRequestID returns the request ID from a zen.Context.
//
//	id := zgin.GetRequestID(c)
func GetRequestID(c zen.Context) string {
	v, ok := c.Get(requestIDKey)
	if !ok {
		return ""
	}
	id, _ := v.(string)
	return id
}
