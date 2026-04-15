package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/yuancore/go-zen/zen"
)

// ViperConfig wraps spf13/viper to implement zen.Config.
type ViperConfig struct {
	mu        sync.RWMutex
	v         *viper.Viper
	files     []*viper.Viper
	overrides map[string]any
}

var _ zen.Config = (*ViperConfig)(nil)

// New reads and merges one or more local config files.
// Missing files are ignored so defaults and environment variables still work.
func New(paths ...string) *ViperConfig {
	c := &ViperConfig{
		v:         newBaseViper(),
		overrides: make(map[string]any),
	}
	_ = c.Load(paths...) // best-effort; missing file uses defaults/env
	return c
}

func newBaseViper() *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return v
}

// Load reads and merges one or more local config files in order.
func (c *ViperConfig) Load(paths ...string) error {
	for _, path := range paths {
		if err := c.AddConfigFile(path); err != nil {
			return err
		}
	}
	return nil
}

// AddConfigFile loads a local file, merges it into the config, and watches it for changes.
func (c *ViperConfig) AddConfigFile(path string) error {
	fileViper := newBaseViper()
	fileViper.SetConfigFile(path)
	if err := fileViper.ReadInConfig(); err != nil {
		return err
	}

	c.mu.Lock()
	c.files = append(c.files, fileViper)
	err := c.rebuildLocked()
	c.mu.Unlock()
	if err != nil {
		return err
	}

	fileViper.OnConfigChange(func(fsnotify.Event) {
		if err := fileViper.ReadInConfig(); err != nil {
			return
		}

		c.mu.Lock()
		_ = c.rebuildLocked()
		c.mu.Unlock()
	})
	fileViper.WatchConfig()
	return nil
}

func (c *ViperConfig) rebuildLocked() error {
	merged := newBaseViper()
	for _, fileViper := range c.files {
		if err := merged.MergeConfigMap(fileViper.AllSettings()); err != nil {
			return err
		}
	}
	for key, value := range c.overrides {
		merged.Set(key, value)
	}
	c.v = merged
	return nil
}

func (c *ViperConfig) GetString(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v.GetString(key)
}

func (c *ViperConfig) GetInt(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v.GetInt(key)
}

func (c *ViperConfig) GetBool(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v.GetBool(key)
}

func (c *ViperConfig) GetFloat64(key string) float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v.GetFloat64(key)
}

func (c *ViperConfig) GetStringSlice(key string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v.GetStringSlice(key)
}

func (c *ViperConfig) GetStringMap(key string) map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v.GetStringMap(key)
}

func (c *ViperConfig) IsSet(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v.IsSet(key)
}

func (c *ViperConfig) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.overrides[key] = value
	c.v.Set(key, value)
}

func (c *ViperConfig) Sub(key string) zen.Config {
	c.mu.RLock()
	defer c.mu.RUnlock()

	sub := c.v.Sub(key)
	if sub == nil {
		return &ViperConfig{
			v:         newBaseViper(),
			overrides: make(map[string]any),
		}
	}
	return &ViperConfig{
		v:         sub,
		overrides: make(map[string]any),
	}
}

func (c *ViperConfig) Unmarshal(key string, v any) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if key == "" {
		return c.v.Unmarshal(v)
	}
	sub := c.v.Sub(key)
	if sub == nil {
		return nil
	}
	return sub.Unmarshal(v)
}

// Raw returns the underlying *viper.Viper for advanced use.
func (c *ViperConfig) Raw() *viper.Viper { return c.v }

// ---------- Environment-based loading ----------

// DefaultEnvVar is the environment variable used to select the active environment.
const DefaultEnvVar = "APP_ENV"

// NewEnv loads config using environment-based file merging.
//
// It reads the APP_ENV environment variable (default "dev") and loads:
//
//	{dir}/config.toml          ← base (shared settings)
//	{dir}/config_{env}.toml    ← env overlay (overrides base)
//
// Example:
//
//	cfg := config.NewEnv("./config")
//	// APP_ENV=prod  →  loads config.toml + config_prod.toml
//	// APP_ENV=      →  loads config.toml + config_dev.toml  (default)
func NewEnv(dir string) *ViperConfig {
	return NewEnvVar(dir, DefaultEnvVar, "dev")
}

// NewEnvVar is like NewEnv but uses a custom env variable name and default environment name.
//
//	cfg := config.NewEnvVar("./config", "GO_ENV", "dev")
func NewEnvVar(dir, envVar, defaultEnv string) *ViperConfig {
	env := strings.TrimSpace(os.Getenv(envVar))
	if env == "" {
		env = defaultEnv
	}

	base := filepath.Join(dir, "config.toml")
	overlay := filepath.Join(dir, fmt.Sprintf("config_%s.toml", env))

	c := &ViperConfig{
		v:         newBaseViper(),
		overrides: make(map[string]any),
	}
	// base config is required; overlay is optional (ignored if missing)
	_ = c.Load(base)
	if _, err := os.Stat(overlay); err == nil {
		_ = c.Load(overlay)
	}
	return c
}

// Factory returns a zen.Option that builds a ViperConfig from explicit file paths.
//
//	app := zen.New(config.Factory("./config/config.toml"), ...)
func Factory(paths ...string) zen.Option {
	return zen.WithConfig(New(paths...))
}

// EnvFactory returns a zen.Option that builds a ViperConfig via environment-based loading.
// The active environment is read from APP_ENV (default "dev").
//
//	app := zen.New(config.EnvFactory("./config"), zlog.Factory(), zgin.Factory())
func EnvFactory(dir string) zen.Option {
	return zen.WithConfig(NewEnv(dir))
}
