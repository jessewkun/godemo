package provider

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jessewkun/gocommon/cache"
)

type MainCache struct{ *redis.Client }

type MainCacheName string

var MainCacheNameValue MainCacheName = "main"

func ProvideMainCache(name MainCacheName) MainCache {
	conn, err := cache.GetConn(string(name))
	if err != nil {
		panic(fmt.Errorf("get cache conn error: %s", err))
	}
	return MainCache{conn}
}
