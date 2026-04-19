package zredis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yuancore/go-zen/zen"
)

// RedisCache 包装 *redis.Client 并实现 zen.Cache 接口。
// RedisCache wraps *redis.Client and implements zen.Cache.
type RedisCache struct {
	rdb *redis.Client
}

// 编译期接口检查。
// compile-time interface check.
var _ zen.Cache = (*RedisCache)(nil)

// newClient 根据 InstanceConfig 创建 Redis 客户端。
// MinIdleConns 在高并发场景下可保持连接预热，降低冷连接建立延迟。
// newClient creates a Redis client from InstanceConfig.
// MinIdleConns keeps connections warm to reduce cold-start latency under high concurrency.
func newClient(cfg InstanceConfig) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.effectiveAddress(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.effectivePoolSize(),
		MinIdleConns: cfg.effectiveMinIdleConns(),
		MaxRetries:   cfg.effectiveMaxRetries(),
		DialTimeout:  cfg.effectiveDialTimeout(),
		ReadTimeout:  cfg.effectiveReadTimeout(),
		WriteTimeout: cfg.effectiveWriteTimeout(),
	})
	return &RedisCache{rdb: rdb}
}

// mapErr 将 redis.Nil 转换为 ErrNil，便于调用方统一使用 errors.Is。
// mapErr translates redis.Nil to ErrNil so callers can use errors.Is uniformly.
func mapErr(err error) error {
	if errors.Is(err, redis.Nil) {
		return ErrNil
	}
	return err
}

// ---------- Lifecycle ----------

// Ping 检测 Redis 连接是否正常。
// Ping checks the Redis connection.
func (c *RedisCache) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Close 关闭底层 Redis 客户端连接。
// Close closes the underlying Redis client.
func (c *RedisCache) Close() error {
	return c.rdb.Close()
}

// ---------- StringOps ----------

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	v, err := c.rdb.Get(ctx, key).Result()
	return v, mapErr(err)
}

func (c *RedisCache) MGet(ctx context.Context, keys ...string) ([]string, error) {
	vals, err := c.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	result := make([]string, len(vals))
	for i, v := range vals {
		if v != nil {
			result[i], _ = v.(string)
		}
	}
	return result, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

func (c *RedisCache) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	return c.rdb.SetNX(ctx, key, value, expiration).Result()
}

func (c *RedisCache) Del(ctx context.Context, keys ...string) (int64, error) {
	return c.rdb.Del(ctx, keys...).Result()
}

func (c *RedisCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.rdb.Exists(ctx, keys...).Result()
}

func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.rdb.Expire(ctx, key, expiration).Err()
}

func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.rdb.TTL(ctx, key).Result()
}

func (c *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return c.rdb.Incr(ctx, key).Result()
}

func (c *RedisCache) IncrBy(ctx context.Context, key string, delta int64) (int64, error) {
	return c.rdb.IncrBy(ctx, key, delta).Result()
}

// ---------- ListOps ----------

func (c *RedisCache) LPush(ctx context.Context, key string, values ...any) error {
	return c.rdb.LPush(ctx, key, values...).Err()
}

func (c *RedisCache) RPush(ctx context.Context, key string, values ...any) error {
	return c.rdb.RPush(ctx, key, values...).Err()
}

func (c *RedisCache) LPop(ctx context.Context, key string) (string, error) {
	v, err := c.rdb.LPop(ctx, key).Result()
	return v, mapErr(err)
}

func (c *RedisCache) RPop(ctx context.Context, key string) (string, error) {
	v, err := c.rdb.RPop(ctx, key).Result()
	return v, mapErr(err)
}

func (c *RedisCache) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.rdb.LRange(ctx, key, start, stop).Result()
}

func (c *RedisCache) LLen(ctx context.Context, key string) (int64, error) {
	return c.rdb.LLen(ctx, key).Result()
}

func (c *RedisCache) LRem(ctx context.Context, key string, count int64, element any) (int64, error) {
	return c.rdb.LRem(ctx, key, count, element).Result()
}

// ---------- SetOps ----------

func (c *RedisCache) SAdd(ctx context.Context, key string, members ...any) (int64, error) {
	return c.rdb.SAdd(ctx, key, members...).Result()
}

func (c *RedisCache) SRem(ctx context.Context, key string, members ...any) (int64, error) {
	return c.rdb.SRem(ctx, key, members...).Result()
}

func (c *RedisCache) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.rdb.SMembers(ctx, key).Result()
}

func (c *RedisCache) SIsMember(ctx context.Context, key string, member any) (bool, error) {
	return c.rdb.SIsMember(ctx, key, member).Result()
}

func (c *RedisCache) SCard(ctx context.Context, key string) (int64, error) {
	return c.rdb.SCard(ctx, key).Result()
}

func (c *RedisCache) SInter(ctx context.Context, keys ...string) ([]string, error) {
	return c.rdb.SInter(ctx, keys...).Result()
}

func (c *RedisCache) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	return c.rdb.SUnion(ctx, keys...).Result()
}

// ---------- ZSetOps ----------

func (c *RedisCache) ZAdd(ctx context.Context, key string, members ...zen.ZMember) error {
	zm := make([]redis.Z, len(members))
	for i, m := range members {
		zm[i] = redis.Z{Score: m.Score, Member: m.Member}
	}
	return c.rdb.ZAdd(ctx, key, zm...).Err()
}

func (c *RedisCache) ZRem(ctx context.Context, key string, members ...any) (int64, error) {
	return c.rdb.ZRem(ctx, key, members...).Result()
}

func (c *RedisCache) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.rdb.ZRange(ctx, key, start, stop).Result()
}

func (c *RedisCache) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]zen.ZMember, error) {
	vals, err := c.rdb.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	result := make([]zen.ZMember, len(vals))
	for i, v := range vals {
		result[i] = zen.ZMember{Score: v.Score, Member: fmt.Sprint(v.Member)}
	}
	return result, nil
}

func (c *RedisCache) ZRangeByScore(ctx context.Context, key, min, max string) ([]string, error) {
	return c.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{Min: min, Max: max}).Result()
}

func (c *RedisCache) ZScore(ctx context.Context, key string, member string) (float64, error) {
	v, err := c.rdb.ZScore(ctx, key, member).Result()
	return v, mapErr(err)
}

func (c *RedisCache) ZCard(ctx context.Context, key string) (int64, error) {
	return c.rdb.ZCard(ctx, key).Result()
}

func (c *RedisCache) ZRank(ctx context.Context, key string, member string) (int64, error) {
	v, err := c.rdb.ZRank(ctx, key, member).Result()
	return v, mapErr(err)
}

func (c *RedisCache) ZRevRank(ctx context.Context, key string, member string) (int64, error) {
	v, err := c.rdb.ZRevRank(ctx, key, member).Result()
	return v, mapErr(err)
}

func (c *RedisCache) ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	return c.rdb.ZIncrBy(ctx, key, increment, member).Result()
}

// ---------- HashOps ----------

func (c *RedisCache) HSet(ctx context.Context, key string, values ...any) error {
	return c.rdb.HSet(ctx, key, values...).Err()
}

func (c *RedisCache) HGet(ctx context.Context, key string, field string) (string, error) {
	v, err := c.rdb.HGet(ctx, key, field).Result()
	return v, mapErr(err)
}

func (c *RedisCache) HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	return c.rdb.HMGet(ctx, key, fields...).Result()
}

func (c *RedisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.rdb.HGetAll(ctx, key).Result()
}

func (c *RedisCache) HDel(ctx context.Context, key string, fields ...string) error {
	return c.rdb.HDel(ctx, key, fields...).Err()
}

func (c *RedisCache) HExists(ctx context.Context, key string, field string) (bool, error) {
	return c.rdb.HExists(ctx, key, field).Result()
}

func (c *RedisCache) HLen(ctx context.Context, key string) (int64, error) {
	return c.rdb.HLen(ctx, key).Result()
}

func (c *RedisCache) HIncrBy(ctx context.Context, key string, field string, delta int64) (int64, error) {
	return c.rdb.HIncrBy(ctx, key, field, delta).Result()
}
