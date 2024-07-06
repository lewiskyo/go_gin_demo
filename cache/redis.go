package cache

import (
	"sync"

	"github.com/go-redis/redis/v8"
)

var redisOnce sync.Once       // 保证只初始化1次
var redisClient *redis.Client // 单节点Redis

func Redis() *redis.Client {
	if redisClient == nil {
		redisOnce.Do(func() {
			rdb := redis.NewClient(&redis.Options{
				Addr:     "0.0.0.0:6379", // Redis地址
				Password: "",             // Redis密码，如果没有则为空字符串
				DB:       0,              // 使用默认DB
			})
			redisClient = rdb
		})
	}

	return redisClient
}
