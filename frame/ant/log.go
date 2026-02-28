package ant

import (
	"github.com/yuancore/go-zen/os/alog"
	"go.uber.org/zap"
)

// Log Get Log content
func Log() *zap.Logger {
	return alog.Write
}
