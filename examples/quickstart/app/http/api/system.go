package api

import (
	"runtime"

	"github.com/yuancore/go-zen/zen"
)

// SystemController handles infrastructure-level endpoints.
type SystemController struct{}

// NewSystemController returns a new SystemController.
func NewSystemController() *SystemController {
	return &SystemController{}
}

// Ping returns a brief liveness check response.
func (ctrl *SystemController) Ping(c zen.Context) {
	c.JSON(200, map[string]any{
		"status":     "ok",
		"go":         runtime.Version(),
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"goroutines": runtime.NumGoroutine(),
	})
}
