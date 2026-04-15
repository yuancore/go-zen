package main

import (
	"path/filepath"
	"strings"

	viperadapter "github.com/yuancore/go-zen/adapter/config"
)

// configBaseDir returns the project root directory inferred from the config
// file path.  Assumes the config file lives inside a directory named "config",
// so the project root is the parent of that directory.
//
// Example:
//
//	configBaseDir("/app/config/config.toml") // → "/app"
func configBaseDir(configPath string) string {
	dir := filepath.Dir(configPath)
	if filepath.Base(dir) == "config" {
		return filepath.Dir(dir)
	}
	return dir
}

// normalizeLoggerPaths converts relative file paths inside the [log] config
// section to absolute paths anchored at configDir.  Non-file paths such as
// "stdout" or "stderr" are left unchanged.
func normalizeLoggerPaths(cfg *viperadapter.ViperConfig, configDir string) {
	resolvePath := func(p string) string {
		if filepath.IsAbs(p) || p == "stdout" || p == "stderr" || strings.HasPrefix(p, "tcp://") {
			return p
		}
		return filepath.Join(configDir, filepath.FromSlash(p))
	}

	if path := cfg.GetString("log.path"); path != "" {
		cfg.Set("log.path", resolvePath(path))
	}

	if paths := cfg.GetStringSlice("log.output_paths"); len(paths) > 0 {
		for i, p := range paths {
			paths[i] = resolvePath(p)
		}
		cfg.Set("log.output_paths", paths)
	}

	if paths := cfg.GetStringSlice("log.error_output_paths"); len(paths) > 0 {
		for i, p := range paths {
			paths[i] = resolvePath(p)
		}
		cfg.Set("log.error_output_paths", paths)
	}
}
