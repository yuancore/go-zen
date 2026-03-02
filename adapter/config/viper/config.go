package viperadapter

import (
	"github.com/spf13/viper"
	"github.com/yuancore/go-zen/zen"
)

// ViperConfig wraps spf13/viper to implement zen.Config.
type ViperConfig struct {
	v *viper.Viper
}

var _ zen.Config = (*ViperConfig)(nil)

// NewConfig reads a YAML/TOML/JSON config file at the given path.
func NewConfig(path string) *ViperConfig {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv()
	_ = v.ReadInConfig() // best-effort; missing file uses defaults/env
	return &ViperConfig{v: v}
}

func (c *ViperConfig) GetString(key string) string        { return c.v.GetString(key) }
func (c *ViperConfig) GetInt(key string) int              { return c.v.GetInt(key) }
func (c *ViperConfig) GetBool(key string) bool            { return c.v.GetBool(key) }
func (c *ViperConfig) GetStringSlice(key string) []string { return c.v.GetStringSlice(key) }

func (c *ViperConfig) Sub(key string) zen.Config {
	sub := c.v.Sub(key)
	if sub == nil {
		return &ViperConfig{v: viper.New()}
	}
	return &ViperConfig{v: sub}
}

func (c *ViperConfig) Unmarshal(key string, v any) error {
	if key == "" {
		return c.v.Unmarshal(v)
	}
	sub := c.v.Sub(key)
	if sub == nil {
		return nil
	}
	return sub.Unmarshal(v)
}
