package zredis

import (
	"fmt"
	"strings"
	"time"
)

const (
	// DefaultConfigKey 是读取 Redis 配置的 key，对应 TOML 中的 [[redis]] 数组。
	// DefaultConfigKey is the config section key used to load Redis settings.
	// It maps to the [[redis]] array in TOML config files.
	DefaultConfigKey = "redis"

	// DefaultServiceName 是默认 zen.Cache 在容器中的注册键。
	// DefaultServiceName is the container key used for the default zen.Cache.
	DefaultServiceName = "cache"

	// DefaultManagerServiceName 是 *Manager 在容器中的注册键。
	// DefaultManagerServiceName is the container key used for the *Manager.
	DefaultManagerServiceName = "redis_manager"

	defaultPoolSize     = 10
	defaultMinIdleConns = 0 // go-redis 默认不保留空闲连接 / go-redis default: no min idle connections
	defaultMaxRetries   = 3
	defaultDialTimeout  = 5 // seconds
	defaultReadTimeout  = 3 // seconds
	defaultWriteTimeout = 3 // seconds
	defaultPingTimeout  = 5 // seconds
)

// configProvider 是该组件所需的最小配置接口。
// configProvider is the minimal config surface required by this component.
type configProvider interface {
	IsSet(key string) bool
	Unmarshal(key string, v any) error
}

// InstanceConfig 描述单个 Redis 连接配置，与 TOML [[redis]] 数组格式兼容：
//
//	[[redis]]
//	name     = "default"
//	address  = "127.0.0.1:6379"
//	password = ""
//	db       = 0
//
// 高并发场景建议设置 pool_size 和 min_idle_conns：
//
//	pool_size     = 100   # 连接池最大连接数
//	min_idle_conns = 20   # 保持最小空闲连接，降低冷连接建立延迟
//
// InstanceConfig describes a single Redis connection.
// It is compatible with the [[redis]] TOML array format.
// For high-concurrency workloads, tune pool_size and min_idle_conns:
//
//	pool_size      = 100  # max connections in the pool
//	min_idle_conns = 20   # keep warm connections to reduce cold-start latency
type InstanceConfig struct {
	Name         string `mapstructure:"name"`
	Address      string `mapstructure:"address"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"` // 最小空闲连接数 / minimum idle connections kept warm
	MaxRetries   int    `mapstructure:"max_retries"`
	DialTimeout  int    `mapstructure:"dial_timeout"`  // seconds
	ReadTimeout  int    `mapstructure:"read_timeout"`  // seconds
	WriteTimeout int    `mapstructure:"write_timeout"` // seconds
	// PingTimeout 是启动连通性检查的超时时间（秒），设为 -1 禁用。
	// PingTimeout is the deadline for the startup connectivity check.
	// Set to -1 to disable the ping. Defaults to 5 seconds.
	PingTimeout int `mapstructure:"ping_timeout"` // seconds
}

// nestedSettings 用于嵌套格式（redis.instances）：
//
//	[redis]
//	default = "main"
//	[[redis.instances]]
//	name = "main"
//
// nestedSettings is used when configs are grouped under a parent key.
type nestedSettings struct {
	Default   string           `mapstructure:"default"`
	Instances []InstanceConfig `mapstructure:"instances"`
}

func (c InstanceConfig) effectiveAddress() string {
	return strings.TrimSpace(c.Address)
}

// validate 校验实例配置的必填字段。
// 配置文件字段名写错（如把 address 写成 host）时，mapstructure 会静默赋零值，
// 此处显式校验可保证启动时 fail-fast，而非连接到意料之外的 Redis 实例。
//
// validate checks required fields for an instance config.
// When a config key name is misspelled (e.g. "host" instead of "address"),
// mapstructure silently assigns zero values; this validation ensures fail-fast
// startup instead of silently connecting to an unintended Redis instance.
func (c InstanceConfig) validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("redis: instance name is required — set [[redis]] name = \"...\" in config")
	}
	if strings.TrimSpace(c.Address) == "" {
		return fmt.Errorf("redis: address is required for instance %q — check [[redis]] address = \"host:port\" in config", c.Name)
	}
	return nil
}

func (c InstanceConfig) effectivePoolSize() int {
	if c.PoolSize <= 0 {
		return defaultPoolSize
	}
	return c.PoolSize
}

// effectiveMinIdleConns 返回最小空闲连接数。
// 高并发场景建议设置为 pool_size 的 20%~50%，保持连接预热。
// effectiveMinIdleConns returns the minimum number of idle connections to keep.
// For high concurrency, set this to 20-50% of pool_size to reduce cold-start latency.
func (c InstanceConfig) effectiveMinIdleConns() int {
	if c.MinIdleConns < 0 {
		return 0
	}
	return c.MinIdleConns
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

// effectivePingTimeout 返回启动健康检查的超时时间，返回 0 表示禁用。
// effectivePingTimeout returns the ping deadline; 0 means ping is disabled.
func (c InstanceConfig) effectivePingTimeout() time.Duration {
	if c.PingTimeout < 0 {
		return 0 // disabled / 已禁用
	}
	if c.PingTimeout == 0 {
		return time.Duration(defaultPingTimeout) * time.Second
	}
	return time.Duration(c.PingTimeout) * time.Second
}

// loadSettings 从配置中读取 Redis 实例列表，支持两种格式：
//
// 格式一 — 扁平数组（与 quickstart 配置兼容）：
//
//	[[redis]]
//	name = "default"
//	address = "127.0.0.1:6379"
//
// 格式二 — 嵌套结构（含显式 default 字段）：
//
//	[redis]
//	default = "main"
//	[[redis.instances]]
//	name = "main"
//	address = "127.0.0.1:6379"
//
// 返回 (instances, defaultName, error)。
//
// loadSettings reads Redis instance configs from the given config provider.
// It supports two formats: flat [[redis]] array or nested [redis] with instances.
// Returns (instances, defaultName, error).
func loadSettings(cfg configProvider, key string) ([]InstanceConfig, string, error) {
	if cfg == nil || !cfg.IsSet(key) {
		return nil, "", nil
	}

	// 优先尝试嵌套格式 (redis.instances)。
	// Try nested format first (redis.instances).
	var nested nestedSettings
	if err := cfg.Unmarshal(key, &nested); err == nil && len(nested.Instances) > 0 {
		defaultName := nested.Default
		if defaultName == "" && len(nested.Instances) > 0 {
			defaultName = nested.Instances[0].Name
		}
		return nested.Instances, defaultName, nil
	}

	// 回退到扁平数组格式 [[redis]]。
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
