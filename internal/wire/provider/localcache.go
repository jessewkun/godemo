package provider

import "github.com/jessewkun/gocommon/db/localcache"

// ProvideLocalCacheManager 创建一个本地缓存管理器
func ProvideLocalCacheManager() *localcache.Manager {
	return localcache.NewManager()
}
