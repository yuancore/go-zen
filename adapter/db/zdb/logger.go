package zdb

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/yuancore/go-zen/zen"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var tablePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bfrom\s+["'\[\]` + "`" + `]*([a-zA-Z0-9_.]+)["'\]\[` + "`" + `]*`),
	regexp.MustCompile(`(?i)\bupdate\s+["'\[\]` + "`" + `]*([a-zA-Z0-9_.]+)["'\]\[` + "`" + `]*`),
	regexp.MustCompile(`(?i)\binsert\s+into\s+["'\[\]` + "`" + `]*([a-zA-Z0-9_.]+)["'\]\[` + "`" + `]*`),
	regexp.MustCompile(`(?i)\bdelete\s+from\s+["'\[\]` + "`" + `]*([a-zA-Z0-9_.]+)["'\]\[` + "`" + `]*`),
	regexp.MustCompile(`(?i)\breplace\s+into\s+["'\[\]` + "`" + `]*([a-zA-Z0-9_.]+)["'\]\[` + "`" + `]*`),
}

type sqlLogger struct {
	logger                    zen.Logger
	logLevel                  gormlogger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

func newSQLLogger(logger zen.Logger, cfg resolvedLogConfig) gormlogger.Interface {
	if logger == nil {
		logger = noopLogger{}
	}
	return &sqlLogger{
		logger:                    logger,
		logLevel:                  cfg.level,
		slowThreshold:             cfg.slowThreshold,
		ignoreRecordNotFoundError: cfg.ignoreRecordNotFoundError,
	}
}

func (l *sqlLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	cloned := *l
	cloned.logLevel = level
	return &cloned
}

func (l *sqlLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.logLevel < gormlogger.Info {
		return
	}
	l.logger.Info(fmt.Sprintf(msg, args...), contextFields(ctx)...)
}

func (l *sqlLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.logLevel < gormlogger.Warn {
		return
	}
	l.logger.Warn(fmt.Sprintf(msg, args...), contextFields(ctx)...)
}

func (l *sqlLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.logLevel < gormlogger.Error {
		return
	}
	l.logger.Error(fmt.Sprintf(msg, args...), contextFields(ctx)...)
}

func (l *sqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel == gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := queryFields(ctx, sql, rows, elapsed)

	switch {
	case err != nil && l.logLevel >= gormlogger.Error && (!l.ignoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		fields = append(fields, "err", err)
		l.logger.Error("db.error", append([]any{"event", "db.error"}, fields...)...)
	case l.slowThreshold > 0 && elapsed > l.slowThreshold && l.logLevel >= gormlogger.Warn:
		fields = append(fields, "slow_threshold_ms", durationMS(l.slowThreshold))
		l.logger.Warn("db.slow", append([]any{"event", "db.slow"}, fields...)...)
	case l.logLevel >= gormlogger.Info:
		l.logger.Info("db.query", append([]any{"event", "db.query"}, fields...)...)
	}
}

type noopLogger struct{}

func (noopLogger) Debug(string, ...any) {}
func (noopLogger) Info(string, ...any)  {}
func (noopLogger) Warn(string, ...any)  {}
func (noopLogger) Error(string, ...any) {}
func (noopLogger) Fatal(string, ...any) {}
func (noopLogger) With(...any) zen.Logger {
	return noopLogger{}
}

func contextFields(ctx context.Context) []any {
	if ctx == nil {
		return nil
	}

	keys := [...]string{"request_id", "trace_id", "span_id"}
	fields := make([]any, 0, len(keys)*2)
	for _, key := range keys {
		if value, ok := lookupContextString(ctx, key); ok {
			fields = append(fields, key, value)
		}
	}
	return fields
}

func lookupContextString(ctx context.Context, key string) (string, bool) {
	if ctx == nil || key == "" {
		return "", false
	}

	if value, ok := ctx.Value(key).(string); ok && value != "" {
		return value, true
	}

	return "", false
}

func queryFields(ctx context.Context, sql string, rows int64, elapsed time.Duration) []any {
	sql = compactSQL(sql)
	fields := make([]any, 0, 12)
	fields = append(fields, "sql", sql, "rows", rows, "duration_ms", durationMS(elapsed))
	if table := tableName(sql); table != "" {
		fields = append(fields, "table", table)
	}
	return append(fields, contextFields(ctx)...)
}

func compactSQL(sql string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(sql)), " ")
}

func tableName(sql string) string {
	for _, pattern := range tablePatterns {
		matches := pattern.FindStringSubmatch(sql)
		if len(matches) == 2 {
			return strings.Trim(matches[1], "\"'[]`")
		}
	}
	return ""
}

func durationMS(d time.Duration) float64 {
	return float64(d.Microseconds()) / 1000
}
