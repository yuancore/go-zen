package serve

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	zconfig "github.com/yuancore/go-zen/adapter/config/viper"
	"github.com/yuancore/go-zen/zen"
)

var defaultConfigCandidates = []string{
	filepath.Join("config", "config.toml"),
	filepath.Join("examples", "quickstart", "config", "config.toml"),
}

func loadConfig(configPath string) (zen.Config, string, error) {
	resolvedConfigPath, err := ResolveConfigPath(configPath)
	if err != nil {
		return nil, "", err
	}

	baseDir := ConfigBaseDir(resolvedConfigPath)
	config := zconfig.NewConfig()
	if err := config.Load(resolvedConfigPath); err != nil {
		return nil, "", fmt.Errorf("load config: %w", err)
	}

	NormalizeLoggerPaths(config, baseDir)
	return config, baseDir, nil
}

func ResolveConfigPath(configPath string) (string, error) {
	for _, candidate := range configCandidates(configPath) {
		absPath, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		if stat, err := os.Stat(absPath); err == nil && !stat.IsDir() {
			return absPath, nil
		}
	}
	return "", fmt.Errorf("config file not found")
}

func configCandidates(configPath string) []string {
	candidates := make([]string, 0, len(defaultConfigCandidates)+1)
	appendCandidate := func(path string) {
		path = strings.TrimSpace(path)
		if path != "" {
			candidates = append(candidates, path)
		}
	}

	appendCandidate(configPath)
	for _, candidate := range defaultConfigCandidates {
		appendCandidate(candidate)
	}
	return candidates
}

func ConfigBaseDir(configPath string) string {
	if strings.TrimSpace(configPath) == "" {
		return "."
	}

	dir := filepath.Dir(configPath)
	if strings.EqualFold(filepath.Base(dir), "config") {
		return filepath.Dir(dir)
	}
	return dir
}

func NormalizeLoggerPaths(cfg zen.Config, baseDir string) {
	normalizeLoggerSectionPaths(cfg, "log", baseDir)
	normalizeLoggerSectionPaths(cfg, "logger", baseDir)
}

func normalizeLoggerSectionPaths(cfg zen.Config, section, baseDir string) {
	if cfg == nil || !cfg.IsSet(section) {
		return
	}

	normalizePath := func(key string) {
		fullKey := section + "." + key
		value := strings.TrimSpace(cfg.GetString(fullKey))
		if value != "" {
			cfg.Set(fullKey, resolveLocalPath(baseDir, value))
		}
	}

	normalizePaths := func(key string) {
		fullKey := section + "." + key
		values := cfg.GetStringSlice(fullKey)
		if len(values) == 0 {
			return
		}

		normalized := make([]string, 0, len(values))
		for _, value := range values {
			normalized = append(normalized, resolveLocalPath(baseDir, value))
		}
		cfg.Set(fullKey, normalized)
	}

	normalizePath("path")
	normalizePath("file_path")
	normalizePaths("output_paths")
	normalizePaths("error_output_paths")
}

func resolveLocalPath(baseDir, raw string) string {
	raw = strings.TrimSpace(raw)
	switch {
	case raw == "":
		return raw
	case isStandardOutput(raw):
		return raw
	case filepath.IsAbs(raw):
		return filepath.Clean(raw)
	case strings.Contains(raw, "://"):
		return raw
	default:
		return filepath.Clean(filepath.Join(baseDir, raw))
	}
}

func isStandardOutput(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "stdout", "stderr":
		return true
	default:
		return false
	}
}

func appName(cfg zen.Config) string {
	if cfg == nil {
		return "go-zen-quickstart"
	}
	if name := strings.TrimSpace(cfg.GetString("system.app_name")); name != "" {
		return name
	}
	return "go-zen-quickstart"
}

func serverAddress(cfg zen.Config) string {
	address := "8080"
	if cfg != nil {
		if configured := strings.TrimSpace(cfg.GetString("system.address")); configured != "" {
			address = configured
		}
	}
	if strings.HasPrefix(address, ":") {
		return address
	}
	return ":" + address
}
