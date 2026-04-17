package zen

import (
	"context"
	"time"
)

// DefaultCacheServiceName is the container key used for the default zen.Cache.
const DefaultCacheServiceName = "cache"

// ZMember represents a sorted-set member with its associated score.
type ZMember struct {
	Score  float64
	Member string
}

// StringOps covers basic key/value cache operations (Redis strings).
type StringOps interface {
	// Get returns the string value stored at key.
	// Returns ("", ErrCacheNil) when the key does not exist.
	Get(ctx context.Context, key string) (string, error)
	// MGet returns the values of all given keys in order.
	// Missing keys are returned as empty strings.
	MGet(ctx context.Context, keys ...string) ([]string, error)
	// Set stores value at key with the given expiration (0 = no expiry).
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	// SetNX sets value only when the key does not already exist.
	// Returns true if the value was set.
	SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error)
	// Del removes one or more keys and returns the number of keys deleted.
	Del(ctx context.Context, keys ...string) (int64, error)
	// Exists reports how many of the given keys exist.
	Exists(ctx context.Context, keys ...string) (int64, error)
	// Expire sets a timeout on key.
	Expire(ctx context.Context, key string, expiration time.Duration) error
	// TTL returns the remaining time-to-live of a key.
	TTL(ctx context.Context, key string) (time.Duration, error)
	// Incr increments the integer stored at key by one.
	Incr(ctx context.Context, key string) (int64, error)
	// IncrBy increments the integer stored at key by delta.
	IncrBy(ctx context.Context, key string, delta int64) (int64, error)
}

// ListOps covers Redis list operations.
type ListOps interface {
	// LPush inserts values at the head of a list.
	LPush(ctx context.Context, key string, values ...any) error
	// RPush inserts values at the tail of a list.
	RPush(ctx context.Context, key string, values ...any) error
	// LPop removes and returns the first element of a list.
	// Returns ("", ErrCacheNil) when the list is empty.
	LPop(ctx context.Context, key string) (string, error)
	// RPop removes and returns the last element of a list.
	// Returns ("", ErrCacheNil) when the list is empty.
	RPop(ctx context.Context, key string) (string, error)
	// LRange returns a slice of elements between start and stop (0-indexed, inclusive).
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	// LLen returns the number of elements in a list.
	LLen(ctx context.Context, key string) (int64, error)
	// LRem removes the first count occurrences of element from a list.
	LRem(ctx context.Context, key string, count int64, element any) (int64, error)
}

// SetOps covers Redis set operations.
type SetOps interface {
	// SAdd adds members to a set. Returns the number of new members.
	SAdd(ctx context.Context, key string, members ...any) (int64, error)
	// SRem removes members from a set. Returns the number of members removed.
	SRem(ctx context.Context, key string, members ...any) (int64, error)
	// SMembers returns all members of a set.
	SMembers(ctx context.Context, key string) ([]string, error)
	// SIsMember reports whether member is in the set.
	SIsMember(ctx context.Context, key string, member any) (bool, error)
	// SCard returns the number of members in a set.
	SCard(ctx context.Context, key string) (int64, error)
	// SInter returns members that exist in every given set.
	SInter(ctx context.Context, keys ...string) ([]string, error)
	// SUnion returns members that exist in at least one of the given sets.
	SUnion(ctx context.Context, keys ...string) ([]string, error)
}

// ZSetOps covers Redis sorted set operations.
type ZSetOps interface {
	// ZAdd adds members with scores to a sorted set.
	ZAdd(ctx context.Context, key string, members ...ZMember) error
	// ZRem removes members from a sorted set. Returns the number removed.
	ZRem(ctx context.Context, key string, members ...any) (int64, error)
	// ZRange returns elements by index range in ascending order.
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	// ZRangeWithScores returns elements with scores by index range in ascending order.
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]ZMember, error)
	// ZRangeByScore returns members with scores between min and max.
	// Use "-inf" and "+inf" for unbounded ranges.
	ZRangeByScore(ctx context.Context, key, min, max string) ([]string, error)
	// ZScore returns the score of a member. Returns (0, ErrCacheNil) if not found.
	ZScore(ctx context.Context, key string, member string) (float64, error)
	// ZCard returns the number of elements in a sorted set.
	ZCard(ctx context.Context, key string) (int64, error)
	// ZRank returns the 0-based ascending rank of a member. Returns (0, ErrCacheNil) if absent.
	ZRank(ctx context.Context, key string, member string) (int64, error)
	// ZRevRank returns the 0-based descending rank of a member.
	ZRevRank(ctx context.Context, key string, member string) (int64, error)
	// ZIncrBy increments the score of a member and returns the new score.
	ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error)
}

// HashOps covers Redis hash operations.
type HashOps interface {
	// HSet sets one or more field-value pairs on a hash.
	// values must alternate: field, value, field, value, ...
	HSet(ctx context.Context, key string, values ...any) error
	// HGet returns the value of a single hash field.
	// Returns ("", ErrCacheNil) when the field does not exist.
	HGet(ctx context.Context, key string, field string) (string, error)
	// HMGet returns the values of the specified fields in order.
	HMGet(ctx context.Context, key string, fields ...string) ([]any, error)
	// HGetAll returns all field-value pairs of a hash.
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	// HDel removes one or more fields from a hash.
	HDel(ctx context.Context, key string, fields ...string) error
	// HExists reports whether a hash field exists.
	HExists(ctx context.Context, key string, field string) (bool, error)
	// HLen returns the number of fields in a hash.
	HLen(ctx context.Context, key string) (int64, error)
	// HIncrBy increments the integer value of a hash field by delta.
	HIncrBy(ctx context.Context, key string, field string, delta int64) (int64, error)
}

// Cache combines all cache operation interfaces.
//
// The default implementation is provided by adapter/cache/zredis.
// To replace it with a custom backend (e.g. Memcached, in-memory), implement
// this interface and call zen.App.RegisterCache.
type Cache interface {
	StringOps
	ListOps
	SetOps
	ZSetOps
	HashOps

	// Ping checks the cache connection health.
	Ping(ctx context.Context) error
	// Close releases all resources held by this cache instance.
	Close() error
}
