package zlog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yuancore/go-zen/zen"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	RequestIDKey = "request_id"
	TraceIDKey   = "trace_id"
	SpanIDKey    = "span_id"
	UserIDKey    = "user_id"
)

var (
	defaultContextKeys = []string{RequestIDKey, TraceIDKey, SpanIDKey, UserIDKey}
	globalLogger       atomic.Pointer[ZapLogger]
	fallbackOnce       sync.Once
	fallbackLogger     *ZapLogger
)

// ZapLogger wraps *zap.Logger to implement zen.Logger with low overhead.
type ZapLogger struct {
	l           *zap.Logger
	contextKeys []string
}

var _ zen.Logger = (*ZapLogger)(nil)

// Sampling controls log sampling for high-throughput production workloads.
type Sampling struct {
	Initial    int `mapstructure:"initial"`
	Thereafter int `mapstructure:"thereafter"`
}

// Options controls how the zap logger is built.
type Options struct {
	Switch            bool           `mapstructure:"switch"`
	Level             string         `mapstructure:"level"`
	Encoding          string         `mapstructure:"encoding"`
	Format            string         `mapstructure:"format"`
	Development       bool           `mapstructure:"development"`
	DisableCaller     bool           `mapstructure:"disable_caller"`
	DisableStacktrace bool           `mapstructure:"disable_stacktrace"`
	DisableSampling   bool           `mapstructure:"disable_sampling"`
	OutputPaths       []string       `mapstructure:"output_paths"`
	ErrorOutputPaths  []string       `mapstructure:"error_output_paths"`
	FilePath          string         `mapstructure:"file_path"`
	Path              string         `mapstructure:"path"`
	Console           bool           `mapstructure:"console"`
	MaxSize           int            `mapstructure:"max_size"`
	MaxBackups        int            `mapstructure:"max_backups"`
	MaxAge            int            `mapstructure:"max_age"`
	Compress          bool           `mapstructure:"compress"`
	Name              string         `mapstructure:"name"`
	ServiceName       string         `mapstructure:"service_name"`
	InitialFields     map[string]any `mapstructure:"initial_fields"`
	ContextKeys       []string       `mapstructure:"context_keys"`
	Sampling          Sampling       `mapstructure:"sampling"`
}

// DefaultOptions returns production-safe logger defaults.
func DefaultOptions() Options {
	return Options{
		Switch:           true,
		Level:            "info",
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		MaxSize:          100,
		MaxBackups:       10,
		MaxAge:           30,
		ContextKeys:      append([]string(nil), defaultContextKeys...),
		Sampling: Sampling{
			Initial:    100,
			Thereafter: 100,
		},
	}
}

// NewLogger creates a production-ready ZapLogger with default options.
func NewLogger() *ZapLogger {
	logger, err := NewLoggerWithOptions(DefaultOptions())
	if err != nil {
		panic(err)
	}
	return logger
}

// MustNewLoggerWithOptions builds a logger from explicit options and panics on error.
func MustNewLoggerWithOptions(opts Options) *ZapLogger {
	logger, err := NewLoggerWithOptions(opts)
	if err != nil {
		panic(err)
	}
	return logger
}

// New builds a logger from the "logger" config section.
func New(cfg zen.Config) (*ZapLogger, error) {
	if cfg != nil {
		switch {
		case cfg.IsSet("logger"):
			return NewLoggerFromConfigKey(cfg, "logger")
		case cfg.IsSet("log"):
			return NewLoggerFromConfigKey(cfg, "log")
		}
	}
	return NewLoggerWithOptions(DefaultOptions())
}

// New builds a logger from config and panics on error.
func MustNewLoggerFromConfig(cfg zen.Config) *ZapLogger {
	logger, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return logger
}

// NewLoggerFromConfigKey builds a logger from the provided config section.
func NewLoggerFromConfigKey(cfg zen.Config, key string) (*ZapLogger, error) {
	opts := DefaultOptions()
	if cfg != nil && key != "" && cfg.IsSet(key) {
		if err := cfg.Unmarshal(key, &opts); err != nil {
			return nil, fmt.Errorf("unmarshal logger config %q: %w", key, err)
		}
	}
	if cfg != nil && key == "" {
		if err := cfg.Unmarshal("", &opts); err != nil {
			return nil, fmt.Errorf("unmarshal logger config: %w", err)
		}
	}
	return NewLoggerWithOptions(opts)
}

// NewLoggerWithOptions builds a logger from explicit options and installs it as the default logger.
func NewLoggerWithOptions(opts Options) (*ZapLogger, error) {
	logger, err := buildLogger(opts)
	if err != nil {
		return nil, err
	}
	SetDefault(logger)
	return logger, nil
}

// SetDefault replaces the process-wide default logger used by package helpers.
func SetDefault(logger *ZapLogger) {
	if logger == nil || logger.l == nil {
		return
	}
	globalLogger.Store(logger)
}

// Default returns the process-wide logger, falling back to a safe production logger.
func Default() *ZapLogger {
	if logger := globalLogger.Load(); logger != nil {
		return logger
	}
	fallbackOnce.Do(func() {
		logger, err := buildLogger(DefaultOptions())
		if err != nil {
			logger = &ZapLogger{
				l:           zap.NewNop(),
				contextKeys: append([]string(nil), defaultContextKeys...),
			}
		}
		fallbackLogger = logger
	})
	return fallbackLogger
}

// L is a short alias for Default.
func L() *ZapLogger { return Default() }

// ContextWithValue stores a string field in context so WithContext can pick it up.
func ContextWithValue(ctx context.Context, key, value string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if strings.TrimSpace(key) == "" || value == "" {
		return ctx
	}
	ctx = context.WithValue(ctx, contextKey(key), value)
	return context.WithValue(ctx, key, value)
}

// ContextWithRequestID stores request_id in context.
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return ContextWithValue(ctx, RequestIDKey, requestID)
}

// ContextWithTraceID stores trace_id in context.
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return ContextWithValue(ctx, TraceIDKey, traceID)
}

// ContextWithSpanID stores span_id in context.
func ContextWithSpanID(ctx context.Context, spanID string) context.Context {
	return ContextWithValue(ctx, SpanIDKey, spanID)
}

// WithContext enriches the logger with fields extracted from context.
func WithContext(ctx context.Context) *ZapLogger { return Default().WithContext(ctx) }

// Ctx is a short alias for WithContext.
func Ctx(ctx context.Context) *ZapLogger { return WithContext(ctx) }

// Debug logs with the default logger.
func Debug(msg string, kv ...any) { Default().Debug(msg, kv...) }

// Info logs with the default logger.
func Info(msg string, kv ...any) { Default().Info(msg, kv...) }

// Warn logs with the default logger.
func Warn(msg string, kv ...any) { Default().Warn(msg, kv...) }

// Error logs with the default logger.
func Error(msg string, kv ...any) { Default().Error(msg, kv...) }

// Fatal logs with the default logger and exits the process.
func Fatal(msg string, kv ...any) { Default().Fatal(msg, kv...) }

// Sync flushes the default logger.
func Sync() error { return Default().Sync() }

// Raw exposes the underlying *zap.Logger for typed-field logging paths.
func (z *ZapLogger) Raw() *zap.Logger { return z.base().l }

// Named returns a named child logger.
func (z *ZapLogger) Named(name string) *ZapLogger {
	base := z.base()
	name = strings.TrimSpace(name)
	if name == "" {
		return base
	}
	return &ZapLogger{
		l:           base.l.Named(name),
		contextKeys: append([]string(nil), base.contextKeys...),
	}
}

// WithContext enriches the logger with fields extracted from context.
func (z *ZapLogger) WithContext(ctx context.Context) *ZapLogger {
	base := z.base()
	if ctx == nil {
		return base
	}
	fields := contextFields(ctx, base.contextKeys)
	if len(fields) == 0 {
		return base
	}
	return &ZapLogger{
		l:           base.l.With(fields...),
		contextKeys: append([]string(nil), base.contextKeys...),
	}
}

// Ctx is a short alias for WithContext.
func (z *ZapLogger) Ctx(ctx context.Context) *ZapLogger { return z.WithContext(ctx) }

func (z *ZapLogger) Debug(msg string, kv ...any) { z.log(zap.DebugLevel, msg, kv...) }
func (z *ZapLogger) Info(msg string, kv ...any)  { z.log(zap.InfoLevel, msg, kv...) }
func (z *ZapLogger) Warn(msg string, kv ...any)  { z.log(zap.WarnLevel, msg, kv...) }
func (z *ZapLogger) Error(msg string, kv ...any) { z.log(zap.ErrorLevel, msg, kv...) }
func (z *ZapLogger) Fatal(msg string, kv ...any) { z.log(zap.FatalLevel, msg, kv...) }

func (z *ZapLogger) With(kv ...any) zen.Logger {
	base := z.base()
	return &ZapLogger{
		l:           base.l.With(kvToFields(kv)...),
		contextKeys: append([]string(nil), base.contextKeys...),
	}
}

// Sync flushes any buffered log entries.
func (z *ZapLogger) Sync() error {
	err := z.base().l.Sync()
	if err == nil {
		return nil
	}
	text := strings.ToLower(err.Error())
	if strings.Contains(text, "invalid argument") ||
		strings.Contains(text, "inappropriate ioctl") ||
		strings.Contains(text, "bad file descriptor") {
		return nil
	}
	return err
}

type contextKey string

func buildLogger(opts Options) (*ZapLogger, error) {
	opts = normalizeOptions(opts)
	if !opts.Switch {
		return &ZapLogger{
			l:           zap.NewNop(),
			contextKeys: append([]string(nil), opts.ContextKeys...),
		}, nil
	}

	level, err := parseLevel(opts.Level)
	if err != nil {
		return nil, err
	}

	output, err := buildWriteSyncer(resolveOutputPaths(opts), opts)
	if err != nil {
		return nil, err
	}

	errorOutput, err := buildWriteSyncer(normalizePaths(opts.ErrorOutputPaths, "stderr"), opts)
	if err != nil {
		return nil, err
	}

	encoder, err := buildEncoder(opts.Encoding)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(encoder, output, zap.NewAtomicLevelAt(level))
	if shouldSample(opts) {
		core = zapcore.NewSamplerWithOptions(core, time.Second, opts.Sampling.Initial, opts.Sampling.Thereafter)
	}

	zapOptions := []zap.Option{zap.ErrorOutput(errorOutput)}
	if opts.Development {
		zapOptions = append(zapOptions, zap.Development())
	}
	if !opts.DisableCaller {
		zapOptions = append(zapOptions, zap.AddCaller())
	}
	if !opts.DisableStacktrace {
		stackLevel := zapcore.ErrorLevel
		if opts.Development {
			stackLevel = zapcore.WarnLevel
		}
		zapOptions = append(zapOptions, zap.AddStacktrace(stackLevel))
	}

	if fields := initialFields(opts); len(fields) > 0 {
		zapOptions = append(zapOptions, zap.Fields(fields...))
	}

	logger := zap.New(core, zapOptions...)
	if name := strings.TrimSpace(opts.Name); name != "" {
		logger = logger.Named(name)
	}

	return &ZapLogger{
		l:           logger,
		contextKeys: append([]string(nil), opts.ContextKeys...),
	}, nil
}

func normalizeOptions(opts Options) Options {
	defaults := DefaultOptions()

	if strings.TrimSpace(opts.Level) == "" {
		opts.Level = defaults.Level
	}
	if strings.TrimSpace(opts.Encoding) == "" && strings.TrimSpace(opts.Format) != "" {
		opts.Encoding = opts.Format
	}
	if strings.TrimSpace(opts.Encoding) == "" {
		opts.Encoding = defaults.Encoding
	}
	if opts.MaxSize <= 0 {
		opts.MaxSize = defaults.MaxSize
	}
	if opts.MaxBackups <= 0 {
		opts.MaxBackups = defaults.MaxBackups
	}
	if opts.MaxAge <= 0 {
		opts.MaxAge = defaults.MaxAge
	}
	if len(opts.ContextKeys) == 0 {
		opts.ContextKeys = append([]string(nil), defaults.ContextKeys...)
	}
	if opts.Sampling.Initial <= 0 {
		opts.Sampling.Initial = defaults.Sampling.Initial
	}
	if opts.Sampling.Thereafter <= 0 {
		opts.Sampling.Thereafter = defaults.Sampling.Thereafter
	}
	if len(trimPaths(opts.OutputPaths)) == 0 && strings.TrimSpace(opts.FilePath) == "" && strings.TrimSpace(opts.Path) == "" {
		opts.OutputPaths = append([]string(nil), defaults.OutputPaths...)
	}
	if len(trimPaths(opts.ErrorOutputPaths)) == 0 {
		opts.ErrorOutputPaths = append([]string(nil), defaults.ErrorOutputPaths...)
	}

	return opts
}

func shouldSample(opts Options) bool {
	return !opts.Development &&
		!opts.DisableSampling &&
		opts.Sampling.Initial > 0 &&
		opts.Sampling.Thereafter > 0
}

func parseLevel(raw string) (zapcore.Level, error) {
	levelText := strings.ToLower(strings.TrimSpace(raw))
	if levelText == "" {
		levelText = DefaultOptions().Level
	}
	if levelText == "all" {
		levelText = "debug"
	}
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(levelText)); err != nil {
		return zapcore.InfoLevel, fmt.Errorf("parse logger level %q: %w", raw, err)
	}
	return level, nil
}

func buildEncoder(raw string) (zapcore.Encoder, error) {
	encoding := strings.ToLower(strings.TrimSpace(raw))
	switch encoding {
	case "", "json":
		return zapcore.NewJSONEncoder(newEncoderConfig(false)), nil
	case "console":
		return zapcore.NewConsoleEncoder(newEncoderConfig(true)), nil
	default:
		return nil, fmt.Errorf("unsupported logger encoding %q", raw)
	}
}

func newEncoderConfig(console bool) zapcore.EncoderConfig {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	if console {
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
		return cfg
	}
	cfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	return cfg
}

func resolveOutputPaths(opts Options) []string {
	paths := trimPaths(opts.OutputPaths)
	filePath := strings.TrimSpace(firstNonEmpty(opts.FilePath, opts.Path))
	if filePath != "" && !containsPath(paths, filePath) {
		paths = append(paths, filePath)
	}
	if opts.Console && !containsStdIO(paths) {
		paths = append(paths, "stdout")
	}
	if len(paths) == 0 {
		return []string{"stdout"}
	}
	return paths
}

func normalizePaths(paths []string, fallback string) []string {
	paths = trimPaths(paths)
	if len(paths) == 0 {
		return []string{fallback}
	}
	return paths
}

func trimPaths(paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, path := range paths {
		trimmed := strings.TrimSpace(path)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}

func containsStdIO(paths []string) bool {
	for _, path := range paths {
		switch strings.ToLower(strings.TrimSpace(path)) {
		case "stdout", "stderr":
			return true
		}
	}
	return false
}

func containsPath(paths []string, target string) bool {
	for _, path := range paths {
		if strings.EqualFold(strings.TrimSpace(path), target) {
			return true
		}
	}
	return false
}

func buildWriteSyncer(paths []string, opts Options) (zapcore.WriteSyncer, error) {
	writers := make([]zapcore.WriteSyncer, 0, len(paths))
	for _, path := range paths {
		trimmed := strings.TrimSpace(path)
		switch strings.ToLower(trimmed) {
		case "":
			continue
		case "stdout":
			writers = append(writers, zapcore.AddSync(os.Stdout))
		case "stderr":
			writers = append(writers, zapcore.AddSync(os.Stderr))
		default:
			if err := ensureLogDir(trimmed); err != nil {
				return nil, err
			}
			writers = append(writers, zapcore.AddSync(&lumberjack.Logger{
				Filename:   trimmed,
				MaxSize:    opts.MaxSize,
				MaxBackups: opts.MaxBackups,
				MaxAge:     opts.MaxAge,
				Compress:   opts.Compress,
				LocalTime:  true,
			}))
		}
	}
	if len(writers) == 0 {
		return zapcore.AddSync(os.Stdout), nil
	}
	return zapcore.NewMultiWriteSyncer(writers...), nil
}

func ensureLogDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create log dir %q: %w", dir, err)
	}
	return nil
}

func initialFields(opts Options) []zap.Field {
	fields := make([]zap.Field, 0, len(opts.InitialFields)+1)

	if service := strings.TrimSpace(opts.ServiceName); service != "" {
		if _, exists := opts.InitialFields["service"]; !exists {
			fields = append(fields, zap.String("service", service))
		}
	}

	keys := make([]string, 0, len(opts.InitialFields))
	for key := range opts.InitialFields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fields = append(fields, zap.Any(key, opts.InitialFields[key]))
	}

	return fields
}

func (z *ZapLogger) base() *ZapLogger {
	if z != nil && z.l != nil {
		return z
	}
	return Default()
}

func (z *ZapLogger) log(level zapcore.Level, msg string, kv ...any) {
	base := z.base()
	if ce := base.l.Check(level, msg); ce != nil {
		ce.Write(kvToFields(kv)...)
	}
}

func kvToFields(kv []any) []zap.Field {
	if len(kv) == 0 {
		return nil
	}

	fields := make([]zap.Field, 0, len(kv))
	for i := 0; i < len(kv); {
		switch value := kv[i].(type) {
		case zap.Field:
			fields = append(fields, value)
			i++
		case []zap.Field:
			fields = append(fields, value...)
			i++
		default:
			if i+1 >= len(kv) {
				fields = append(fields, zap.Any("field_"+strconv.Itoa(len(fields)), kv[i]))
				i++
				continue
			}
			fields = append(fields, zap.Any(fieldKey(kv[i], len(fields)), kv[i+1]))
			i += 2
		}
	}
	return fields
}

func fieldKey(key any, index int) string {
	switch value := key.(type) {
	case string:
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	case fmt.Stringer:
		if text := strings.TrimSpace(value.String()); text != "" {
			return text
		}
	case error:
		if text := strings.TrimSpace(value.Error()); text != "" {
			return text
		}
	}
	text := strings.TrimSpace(fmt.Sprint(key))
	if text == "" || text == "<nil>" {
		return "field_" + strconv.Itoa(index)
	}
	return text
}

func contextFields(ctx context.Context, keys []string) []zap.Field {
	fields := make([]zap.Field, 0, len(keys))
	seen := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		if value, ok := lookupContextString(ctx, key); ok {
			fields = append(fields, zap.String(key, value))
		}
	}
	return fields
}

func lookupContextString(ctx context.Context, key string) (string, bool) {
	if ctx == nil || key == "" {
		return "", false
	}
	if value, ok := asString(ctx.Value(contextKey(key))); ok {
		return value, true
	}
	if value, ok := asString(ctx.Value(key)); ok {
		return value, true
	}
	return "", false
}

func asString(value any) (string, bool) {
	switch v := value.(type) {
	case string:
		if v != "" {
			return v, true
		}
	case fmt.Stringer:
		if text := v.String(); text != "" {
			return text, true
		}
	}
	return "", false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
