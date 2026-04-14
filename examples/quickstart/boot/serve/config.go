package serve

import (
	"fmt"

	"github.com/yuancore/go-zen/adapter/config"
	"github.com/yuancore/go-zen/zen"
)

// loadConfig reads the config file at the given path.
func loadConfig(path string) (zen.Config, error) {
	cfg := config.New()
	if err := cfg.Load(path); err != nil {
		return nil, fmt.Errorf("load config %s: %w", path, err)
	}
	return cfg, nil
}
