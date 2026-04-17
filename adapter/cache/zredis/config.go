package zredis

import (
	"fmt"
	"strings"
	"time"
)

const (
	// DefaultConfigKey is the config section key used to load Redis settings.
	// It maps to the [[redis]] array in TOML config files.
	DefaultConfigKey = "redis"

	// DefaultServiceName is the container key used for the default zen.Cache.
	DefaultServiceName = "cache"

	// DefaultManagerServiceName is the container key used for the *Manager.
	DefaultManagerServiceName = "redis_manager"

	defaultPoolSize     = 10
	defaultMaxRetries   = 3
	defaultDialTimeout  = 5 // seconds
	defaultReadTimeout  = 3 // seconds
	defaultWriteTimeout = 3 // seconds
	defaultPingTimeout  = 5 // seconds
)

// configProvider is the minimal config surface required by this component.
type configProvider interface {
	IsSet(key string) bool
	Unmarshal(key string, v any) error
}

// InstanceConfig describes a single Redis connection.
// It is compatible with the existing [[redis]] TOML array format:
//
//	[[redis]]
//	name     = "default"
//	address  = "127.0.0.1:6379"
//	password = ""
//	db       = 0
type InstanceConfig struct {
	Name         string `mapstructure:"name"`
	Address      string `mapstructure:"address"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MaxRetries   int    `mapstructure:"max_retries"`
	DialTimeout  int    `mapstructure:"dial_timeout"`  // seconds
	ReadTimeout  int    `mapstructure:"read_timeout"`  // seconds
	WriteTimeout int    `mapstructure:"write_timeout"` // seconds
	// PingTimeout is the deadline for the startup connectivity check.
	// Set to -1 to disable the ping. Defaults to 5 seconds.
	PingTimeout int `mapstructure:"ping_timeout"` // seconds
}

// nestedSettings is used when configs are grouped under a parent key:
//
//	[redis]
//	default = "main"
//	[[redis.instances]]
//	name = "main"
type nestedSettings struct {
	Default   string           `mapstructure:"default"`
	Instances []InstanceConfig `mapstructure:"instances"`
}

func (c InstanceConfig) effectiveAddress() string {
	if strings.TrimSpace(c.Address) == "" {
		return "127.0.0.1:6379"
	}
	return c.Address
}

func (c InstanceConfig) effectivePoolSize() int {
	if c.PoolSize <= 0 {
		return defaultPoolSize
	}
	return c.PoolSize
}

func (c InstanceConfig) effectiveMaxRetries() int {
	if c.MaxRetries < 0 {
		return 0
	}
	if c.MaxRetries == 0 {
		return defaultMaxRetries
	}
	return c.MaxRetries
}

func (c InstanceConfig) effectiveDialTimeout() time.Duration {
	if c.DialTimeout <= 0 {
		return time.Duration(defaultDialTimeout) * time.Second
	}
	return time.Duration(c.DialTimeout) * time.Second
}

func (c InstanceConfig) effectiveReadTimeout() time.Duration {
	if c.ReadTimeout <= 0 {
		return time.Duration(defaultReadTimeout) * time.Second
	}
	return time.Duration(c.ReadTimeout) * time.Second
}

func (c InstanceConfig) effectiveWriteTimeout() time.Duration {
	if c.WriteTimeout <= 0 {
		return time.Duration(defaultWriteTimeout) * time.Second
	}
	return time.Duration(c.WriteTimeout) * time.Second
}

// effectivePingTimeout returns the ping deadline for startup health-check.
// Returns 0 when ping is explicitly disabled (PingTimeout == -1).
func (c InstanceConfig) effectivePingTimeout() time.Duration {
	if c.PingTimeout < 0 {
		return 0 // disabled
	}
	if c.PingTimeout == 0 {
		return time.Duration(defaultPingTimeout) * time.Second
	}
	return time.Duration(c.PingTimeout) * time.Second
}

// loadSettings reads Redis instance configs from the given config provider.
// It supports two formats:
//
// Format 1 – flat array (matches the existing quickstart config):
//
//	[[redis]]
//	name = "default"
//	address = "127.0.0.1:6379"
//
// Format 2 – nested with explicit default:
//
//	[redis]
//	default = "main"
//	[[redis.instances]]
//	name = "main"
//	address = "127.0.0.1:6379"
//
// Returns (instances, defaultName, error).
func loadSettings(cfg configProvider, key string) ([]InstanceConfig, string, error) {
	if cfg == nil || !cfg.IsSet(key) {
		return nil, "", nil
	}

	// Try nested format first (redis.instances).
	var nested nestedSettings
	if err := cfg.Unmarshal(key, &nested); err == nil && len(nested.Instances) > 0 {
		defaultName := nested.Default
		if defaultName == "" && len(nested.Instances) > 0 {
			defaultName = nested.Instances[0].Name
		}
		return nested.Instances, defaultName, nil
	}

	// Fall back to flat array format: [[redis]] entries.
	var instances []InstanceConfig
	if err := cfg.Unmarshal(key, &instances); err != nil {
		return nil, "", fmt.Errorf("unmarshal redis config %q: %w", key, err)
	}

	defaultName := ""
	if len(instances) > 0 {
		defaultName = instances[0].Name
	}
	return instances, defaultName, nil
}
