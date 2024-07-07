package cache

import (
	"sync"
	"time"

	gcache "github.com/patrickmn/go-cache"
)

type Cache struct {
	CM *gcache.Cache
}

// 定义全局变量和互斥锁
var (
	once  sync.Once
	cache *Cache
)

func InitLocalCache() {
	once.Do(func() {
		c := gcache.New(5*time.Minute, 10*time.Minute)
		cache = new(Cache)
		cache.CM = c
	})
}

func GetLocalCache() *Cache {
	return cache
}
