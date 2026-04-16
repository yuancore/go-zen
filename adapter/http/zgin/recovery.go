package zgin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuancore/go-zen/zen"
)

const recoveryBodyMaxSize = 4 << 20 // 4 MB

// NewRecovery returns a production-grade panic recovery gin middleware.
//
// On panic it logs: request_id, method, path, client_ip, panic value,
// first user-code location (panic_at), trimmed stack frames, and request body.
// The body is cached before dispatch so it is available even after the handler consumes it.
//
// Usage:
//
//	engine.Raw().Use(zgin.NewRecovery(logger))
func NewRecovery(logger zen.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		body := cacheBodyForRecovery(c)
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			stack := debug.Stack()
			logger.Error("http panic recovered",
				"request_id", requestIDFromContext(c),
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"client_ip", c.ClientIP(),
				"panic", fmt.Sprintf("%v", rec),
				"panic_at", panicLocation(stack),
				"stack", panicStack(stack, 20),
				"request_body", body,
			)
			if !c.IsAborted() {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// RecoveryHandler wraps NewRecovery as a zen.Handler for use with app.Middleware().
//
//	app.Middleware(zgin.RecoveryHandler(logger))
func RecoveryHandler(logger zen.Logger) zen.Handler {
	h := NewRecovery(logger)
	return func(c zen.Context) {
		if gc, ok := c.(*ginContext); ok {
			h(gc.c)
		} else {
			c.Next()
		}
	}
}

// cacheBodyForRecovery reads the request body before the handler consumes it,
// restores it, and returns a log-friendly value (parsed JSON or plain string).
func cacheBodyForRecovery(c *gin.Context) any {
	if c.Request == nil || c.Request.Body == nil {
		return nil
	}
	data, err := io.ReadAll(io.LimitReader(c.Request.Body, recoveryBodyMaxSize))
	c.Request.Body = io.NopCloser(bytes.NewReader(data)) // always restore
	if err != nil || len(data) == 0 {
		return nil
	}
	return parseBodyForLog(data, recoveryBodyMaxSize)
}

// panicLocation returns the first user-code file:line from a stack trace.
// It skips runtime, vendor, and module cache frames to point directly at
// the line that caused the panic.
func panicLocation(stack []byte) string {
	lines := bytes.Split(stack, []byte("\n"))
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(string(lines[i]))
		if line == "" || !strings.Contains(line, ".go:") {
			continue
		}
		if isSystemFrame(line) {
			continue
		}
		fn := ""
		if i > 0 {
			fn = strings.TrimSpace(string(lines[i-1]))
		}
		if fn != "" {
			return fn + " @ " + line
		}
		return line
	}
	return "unknown"
}

// panicStack returns up to max non-empty lines from the stack trace.
func panicStack(stack []byte, max int) []string {
	lines := bytes.Split(stack, []byte("\n"))
	out := make([]string, 0, max)
	for _, l := range lines {
		if s := strings.TrimSpace(string(l)); s != "" {
			out = append(out, s)
			if len(out) >= max {
				break
			}
		}
	}
	return out
}

// isSystemFrame reports whether a stack file line belongs to runtime or dependencies.
func isSystemFrame(path string) bool {
	return strings.Contains(path, "runtime/") ||
		strings.Contains(path, "/vendor/") ||
		strings.Contains(path, "/pkg/mod/")
}
