package viperadapter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfigMergesFilesInOrder(t *testing.T) {
	dir := t.TempDir()
	basePath := filepath.Join(dir, "base.toml")
	overridePath := filepath.Join(dir, "override.toml")

	if err := os.WriteFile(basePath, []byte(`
[app]
name = "go-zen"

[logger]
level = "info"
output_paths = ["stdout"]
`), 0o600); err != nil {
		t.Fatalf("write base config: %v", err)
	}

	if err := os.WriteFile(overridePath, []byte(`
[logger]
level = "debug"

[db]
host = "127.0.0.1"
`), 0o600); err != nil {
		t.Fatalf("write override config: %v", err)
	}

	cfg := NewConfig(basePath, overridePath)

	if got := cfg.GetString("app.name"); got != "go-zen" {
		t.Fatalf("expected app.name=go-zen, got %q", got)
	}
	if got := cfg.GetString("logger.level"); got != "debug" {
		t.Fatalf("expected logger.level=debug, got %q", got)
	}
	if got := cfg.GetString("db.host"); got != "127.0.0.1" {
		t.Fatalf("expected db.host=127.0.0.1, got %q", got)
	}

	cfg.Set("logger.level", "warn")
	if got := cfg.GetString("logger.level"); got != "warn" {
		t.Fatalf("expected runtime override logger.level=warn, got %q", got)
	}
}

func TestNewConfigSupportsEnvOverride(t *testing.T) {
	t.Setenv("LOGGER_LEVEL", "error")

	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(configPath, []byte(`
[logger]
level = "info"
`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg := NewConfig(configPath)
	if got := cfg.GetString("logger.level"); got != "error" {
		t.Fatalf("expected env override logger.level=error, got %q", got)
	}
}
