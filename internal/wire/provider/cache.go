package provider

import (
	"github.com/go-redis/redis/v8"
	"github.com/jessewkun/gocommon/cache"
)

type MainCache struct{ *redis.Client }

type MainCacheName string

var MainCacheNameValue MainCacheName = "main"

func ProvideMainCache(name MainCacheName) MainCache {
	return MainCache{cache.GetConn(string(name))}
}
