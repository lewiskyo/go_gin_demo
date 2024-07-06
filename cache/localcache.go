package cache

import (
	"sync"
	"time"

	gcache "github.com/patrickmn/go-cache"
)

type Cache struct {
	CM *gcache.Cache
}

var cache Cache
var once sync.Once

func LocalCache() (*Cache, error) {
	once.Do(func() {
		gcache := gcache.New(5*time.Minute, 10*time.Minute)

		cache.CM = gcache
	})
	return &cache, nil
}
