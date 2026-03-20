package zapadapter

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	viperadapter "github.com/yuancore/go-zen/adapter/config/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestWithContextAddsStructuredFields(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	logger := &ZapLogger{
		l:           zap.New(core),
		contextKeys: []string{RequestIDKey, TraceIDKey},
	}

	ctx := ContextWithRequestID(context.Background(), "req-1")
	ctx = ContextWithTraceID(ctx, "trace-1")

	logger.WithContext(ctx).Info("create order", "order_id", 42, zap.String("tenant", "acme"))

	entries := recorded.AllUntimed()
	if len(entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(entries))
	}

	fields := entries[0].ContextMap()
	if got := fields[RequestIDKey]; got != "req-1" {
		t.Fatalf("expected request_id=req-1, got %#v", got)
	}
	if got := fields[TraceIDKey]; got != "trace-1" {
		t.Fatalf("expected trace_id=trace-1, got %#v", got)
	}
	switch got := fields["order_id"].(type) {
	case int:
		if got != 42 {
			t.Fatalf("expected order_id=42, got %#v", got)
		}
	case int64:
		if got != 42 {
			t.Fatalf("expected order_id=42, got %#v", got)
		}
	default:
		t.Fatalf("expected numeric order_id, got %#v", fields["order_id"])
	}
	if got := fields["tenant"]; got != "acme" {
		t.Fatalf("expected tenant=acme, got %#v", got)
	}
}

func TestKVToFieldsSupportsZapFieldsAndOddValues(t *testing.T) {
	fields := kvToFields([]any{
		zap.String("tenant", "acme"),
		"request_id", "req-1",
		"lonely",
	})

	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
	if fields[0].Key != "tenant" {
		t.Fatalf("expected first field key tenant, got %q", fields[0].Key)
	}
	if fields[1].Key != "request_id" {
		t.Fatalf("expected second field key request_id, got %q", fields[1].Key)
	}
	if fields[2].Key != "field_2" {
		t.Fatalf("expected odd field key field_2, got %q", fields[2].Key)
	}
}

func TestNewLoggerFromConfigSupportsLegacyLogSection(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")

	if err := os.WriteFile(configPath, []byte(`
[log]
format = "json"
level = "all"
service_name = "quickstart"
output_paths = ["stdout"]
`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg := viperadapter.NewConfig(configPath)
	logger, err := NewLoggerFromConfig(cfg)
	if err != nil {
		t.Fatalf("new logger from legacy config: %v", err)
	}

	if !logger.Raw().Core().Enabled(zap.DebugLevel) {
		t.Fatalf("expected level=all to enable debug logging")
	}
}

func TestNewLoggerFromConfigHonorsSwitch(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")

	if err := os.WriteFile(configPath, []byte(`
[log]
switch = false
level = "debug"
output_paths = ["stdout"]
`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg := viperadapter.NewConfig(configPath)
	logger, err := NewLoggerFromConfig(cfg)
	if err != nil {
		t.Fatalf("new logger with switch=false: %v", err)
	}

	if logger.Raw().Core().Enabled(zap.ErrorLevel) {
		t.Fatalf("expected switch=false to disable logger output")
	}
}
