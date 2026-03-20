package main

import (
	"path/filepath"
	"testing"

	viperadapter "github.com/yuancore/go-zen/adapter/config/viper"
)

func TestNormalizeLoggerPathsResolvesRelativeToConfigDir(t *testing.T) {
	configDir := filepath.Join(t.TempDir(), "config")
	cfg := viperadapter.NewConfig()
	cfg.Set("log", map[string]any{
		"path":               "./log/app.log",
		"output_paths":       []string{"stdout", "./log/access.log"},
		"error_output_paths": []string{"stderr", "./log/error.log"},
	})

	normalizeLoggerPaths(cfg, configDir)

	if got, want := cfg.GetString("log.path"), filepath.Join(configDir, "log", "app.log"); got != want {
		t.Fatalf("expected log.path=%q, got %q", want, got)
	}

	paths := cfg.GetStringSlice("log.output_paths")
	if len(paths) != 2 {
		t.Fatalf("expected 2 output paths, got %d", len(paths))
	}
	if paths[0] != "stdout" {
		t.Fatalf("expected stdout to remain unchanged, got %q", paths[0])
	}
	if want := filepath.Join(configDir, "log", "access.log"); paths[1] != want {
		t.Fatalf("expected resolved output path %q, got %q", want, paths[1])
	}

	errorPaths := cfg.GetStringSlice("log.error_output_paths")
	if len(errorPaths) != 2 {
		t.Fatalf("expected 2 error output paths, got %d", len(errorPaths))
	}
	if errorPaths[0] != "stderr" {
		t.Fatalf("expected stderr to remain unchanged, got %q", errorPaths[0])
	}
	if want := filepath.Join(configDir, "log", "error.log"); errorPaths[1] != want {
		t.Fatalf("expected resolved error path %q, got %q", want, errorPaths[1])
	}
}

func TestConfigBaseDirUsesParentOfConfigFolder(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config", "config.toml")

	if got := configBaseDir(configPath); got != root {
		t.Fatalf("expected config base dir %q, got %q", root, got)
	}
}
