package serve

import (
	"fmt"
	"path/filepath"

	"github.com/yuancore/go-zen/adapter/config"
	"github.com/yuancore/go-zen/zen"
)

// loadConfig reads a single config file at the given path.
// Deprecated: prefer loadEnvConfig for environment-based loading.
func loadConfig(path string) (zen.Config, error) {
	cfg := config.New()
	if err := cfg.Load(path); err != nil {
		return nil, fmt.Errorf("load config %s: %w", path, err)
	}
	return cfg, nil
}

// loadEnvConfig loads config using environment-based merging from a directory.
// It reads APP_ENV (default "dev") and merges:
//
//	{dir}/config.toml           ← shared base
//	{dir}/config_{APP_ENV}.toml ← env overlay
func loadEnvConfig(dir string) (zen.Config, error) {
	basePath := filepath.Join(dir, "config.toml")
	cfg := config.NewEnv(dir)
	if !cfg.IsSet("system") {
		// try the plain config path as a safety fallback
		_ = cfg.Load(basePath)
	}
	return cfg, nil
}
