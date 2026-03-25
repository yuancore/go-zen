package main

import (
	"log"
	"path/filepath"

	"github.com/yuancore/go-zen/examples/quickstart/boot/serve"
	"github.com/yuancore/go-zen/zen"
)

var defaultConfigPath = filepath.Join("config", "config.toml")

func main() {
	if err := serve.RunCLI(defaultConfigPath); err != nil {
		log.Fatalf("quickstart run failed: %v", err)
	}
}

func resolveConfigPath(configPath string) (string, error) {
	return serve.ResolveConfigPath(configPath)
}

func configBaseDir(configPath string) string {
	return serve.ConfigBaseDir(configPath)
}

func normalizeLoggerPaths(cfg zen.Config, configDir string) {
	serve.NormalizeLoggerPaths(cfg, configDir)
}
