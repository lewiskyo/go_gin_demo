package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go_gin_demo/bo"
	"go_gin_demo/cache"
	"go_gin_demo/utils"
	"log"
	"time"
)

func InitMsgConsumer() {
	threads := 5

	for i := 0; i < threads; i++ {
		go func() {
			for {
				rawMsg, ok := utils.MsgQueue.Dequeue()
				if !ok {
					// 无消息消费, 休眠10毫秒
					time.Sleep(10 * time.Millisecond)
					continue
				}

				msg, ok := rawMsg.(bo.Msg)
				if !ok {
					continue
				}

				redisClient := cache.Redis()
				if redisClient == nil {
					log.Printf("get redis instance fail")
					continue
				}

				uid := msg.Uid
				key := fmt.Sprintf("uid:%d", uid)
				val, err := json.Marshal(msg)
				if err != nil {
					log.Printf("msg json encode fail")
					continue
				}
				redisClient.SetEX(context.TODO(), key, val, time.Second*60)

				log.Printf("consume msg: %+v", msg)
			}
		}()
	}
}
