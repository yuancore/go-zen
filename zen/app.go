package zen

import (
	"context"
	"net/http"
)

// Context abstracts the HTTP request/response context.
type Context interface {
	Param(key string) string
	Query(key string) string
	Header(key string) string
	SetHeader(key, value string)
	JSON(code int, v any)
	String(code int, format string, values ...any)
	BindJSON(v any) error
	Status(code int)
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Set(key string, value any)
	Get(key string) (any, bool)
	Next()
	Abort()
}

// Handler is the function signature for HTTP handlers.
type Handler func(Context)

// Middleware is a semantic alias for Handler.
type Middleware = Handler

// RouterGroup is a group of routes sharing a prefix and middleware.
type RouterGroup interface {
	GET(path string, h ...Handler)
	POST(path string, h ...Handler)
	PUT(path string, h ...Handler)
	DELETE(path string, h ...Handler)
	PATCH(path string, h ...Handler)
	Group(prefix string, mw ...Handler) RouterGroup
	Use(mw ...Handler)
}

// Engine abstracts an HTTP server.
type Engine interface {
	RouterGroup
	Start(addr string) error
	Stop(ctx context.Context) error
	AddHealthCheck(name string, check func() error)
}

// Config abstracts application configuration.
type Config interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetStringSlice(key string) []string
	Sub(key string) Config
	Unmarshal(key string, v any) error
}

// Logger abstracts structured logging.
type Logger interface {
	Debug(msg string, kv ...any)
	Info(msg string, kv ...any)
	Warn(msg string, kv ...any)
	Error(msg string, kv ...any)
	Fatal(msg string, kv ...any)
	With(kv ...any) Logger
}

// Module is a pluggable unit of functionality.
type Module interface {
	Name() string
	Init(app *App) error
	Start() error
	Stop(ctx context.Context) error
	Depends() []string
}
