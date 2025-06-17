package provider

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	gocommonredis "github.com/jessewkun/gocommon/db/redis"
)

type MainCache struct{ *redis.Client }

type MainCacheName string

var MainCacheNameValue MainCacheName = "main"

func ProvideMainCache(name MainCacheName) MainCache {
	conn, err := gocommonredis.GetConn(string(name))
	if err != nil {
		panic(fmt.Errorf("get cache conn error: %s", err))
	}
	return MainCache{conn}
}
