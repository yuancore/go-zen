package serve

import (
	"fmt"

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
//	{dir}/config_{APP_ENV}.toml ← env overlay (optional)
//
// Returns an error if the base config.toml cannot be read, so misconfigured
// paths fail loudly at startup instead of running with empty configuration.
func loadEnvConfig(dir string) (zen.Config, error) {
	cfg := config.NewEnv(dir)
	if !cfg.IsSet("system") {
		return nil, fmt.Errorf(
			"load config: no [system] section found after loading %s/config.toml — "+
				"check that the directory exists and contains a valid config.toml", dir)
	}
	return cfg, nil
}
