package zapadapter

import (
	"os"

	"github.com/yuancore/go-zen/zen"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger wraps *zap.SugaredLogger to implement zen.Logger.
type ZapLogger struct {
	s *zap.SugaredLogger
}

var _ zen.Logger = (*ZapLogger)(nil)

// NewLogger creates a production-ready ZapLogger.
func NewLogger() *ZapLogger {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zap.InfoLevel,
	)
	l := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	return &ZapLogger{s: l.Sugar()}
}

func (z *ZapLogger) Debug(msg string, kv ...any) { z.s.Debugw(msg, kv...) }
func (z *ZapLogger) Info(msg string, kv ...any)  { z.s.Infow(msg, kv...) }
func (z *ZapLogger) Warn(msg string, kv ...any)  { z.s.Warnw(msg, kv...) }
func (z *ZapLogger) Error(msg string, kv ...any) { z.s.Errorw(msg, kv...) }
func (z *ZapLogger) Fatal(msg string, kv ...any) { z.s.Fatalw(msg, kv...) }

func (z *ZapLogger) With(kv ...any) zen.Logger {
	return &ZapLogger{s: z.s.With(kv...)}
}

// Sync flushes any buffered log entries.
func (z *ZapLogger) Sync() error { return z.s.Sync() }
