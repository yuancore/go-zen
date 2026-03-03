package zen

import "context"

// Config abstracts application configuration.
// The default implementation uses spf13/viper.
type Config interface {
	// Basic getters
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]any

	// Sub returns a sub-tree of the configuration.
	Sub(key string) Config

	// Unmarshal decodes a config section into a struct.
	// Pass empty key to unmarshal the entire config.
	Unmarshal(key string, v any) error

	// IsSet checks if a key exists.
	IsSet(key string) bool

	// Set sets a config value at runtime.
	Set(key string, value any)
}

// Logger abstracts structured logging.
// The default implementation uses go.uber.org/zap.
type Logger interface {
	Debug(msg string, kv ...any)
	Info(msg string, kv ...any)
	Warn(msg string, kv ...any)
	Error(msg string, kv ...any)
	Fatal(msg string, kv ...any)
	With(kv ...any) Logger
}

// Engine abstracts the HTTP server engine (routing + lifecycle).
// The default implementation uses gin-gonic/gin via the adapter package.
type Engine interface {
	// Routing
	GET(path string, h ...Handler)
	POST(path string, h ...Handler)
	PUT(path string, h ...Handler)
	DELETE(path string, h ...Handler)
	PATCH(path string, h ...Handler)
	Use(mw ...Handler)
	Group(prefix string, mw ...Handler) RouterGroup

	// Lifecycle
	Start(addr string) error
	Stop(ctx context.Context) error
}
