package zen

import (
	"net/http"
)

// Context abstracts the HTTP request/response context.
// It wraps the underlying engine's context but exposes a framework-neutral API.
type Context interface {
	// Request params
	Param(key string) string
	Query(key string) string
	DefaultQuery(key, defaultValue string) string
	Header(key string) string
	SetHeader(key, value string)

	// Request body
	BindJSON(v any) error
	BindQuery(v any) error
	ShouldBind(v any) error

	// Response
	JSON(code int, v any)
	String(code int, format string, values ...any)
	Data(code int, contentType string, data []byte)
	Status(code int)
	Redirect(code int, location string)

	// Underlying
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	ClientIP() string
	FullPath() string

	// Store
	Set(key string, value any)
	Get(key string) (any, bool)
	MustGet(key string) any

	// Flow control
	Next()
	Abort()
	AbortWithStatusJSON(code int, v any)
	IsAborted() bool
}

// Handler is the function signature for HTTP handlers.
type Handler func(Context)

// Middleware is a semantic alias for Handler.
type Middleware = Handler
