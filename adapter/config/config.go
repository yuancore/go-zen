package config

import (
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
