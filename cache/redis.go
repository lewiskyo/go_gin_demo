package cache

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

// 定义全局变量和互斥锁
var (
	rOnce sync.Once     // 保证只初始化1次
	conn  *redis.Client // 单节点Redis
)

func GetRedis() *redis.Client {
	return conn
}

func InitRedis() {
	rOnce.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "0.0.0.0:6379", // Redis地址
			Password: "",             // Redis密码，如果没有则为空字符串
			DB:       0,              // 使用默认DB
		})
		conn = rdb
	})

	fmt.Println("redis connection pool initialized")
}
