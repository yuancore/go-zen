package ant

import (
	"github.com/yuancore/go-zen/db/aredis"
	"github.com/yuancore/go-zen/os/config"
)

// initRedis
func initRedis() {
	redis := config.GetMaps("redis")

	if len(redis) > 0 {
		aredis.New(redis)
	}
}

// Redis Select a different redis
func Redis(name ...string) *aredis.ClientRedis {
	key := ""

	if len(name) > 0 {
		key = name[0]
	}

	val, ok := aredis.Client[key]
	if ok {
		return val
	} else {
		key = config.GetString("redis.0.name")
	}

	return aredis.Client[key]
}
