package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"go_gin_demo/bo"
	"go_gin_demo/cache"
	"go_gin_demo/model.go"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// QueryReq 查询用户最近1次行为请求结构体
type QueryReq struct {
	Uid int64 `form:"uid" binding:"required"`
}

// QueryRsp 查询用户最近1次行为日志响应体
type QueryRsp struct {
	Code    uint64   `json:"code"`
	Data    []bo.Msg `json:"data"`
	Message string   `json:"message"`
}

func (ctl *BaseController) Query(ctx *gin.Context) {
	var req QueryReq
	if err := ctx.BindQuery(&req); err != nil {
		log.Printf("bind query err: %s", err.Error())
		ctx.JSON(http.StatusOK, QueryRsp{Code: 1, Message: "uri param error"})
		return
	}

	redisClient := cache.GetRedis()
	if redisClient == nil {
		log.Printf("get redis instance fail")
		ctx.JSON(http.StatusOK, QueryRsp{Code: 1, Message: "get redis instance fail"})
		return
	}

	uid := req.Uid
	key := fmt.Sprintf("%s:%d", bo.MSG_KEY_PREFIX, uid)
	flag := false
	var from string
	rsp := QueryRsp{Code: 0, Message: "ok"}

	localCache := cache.GetLocalCache()
	localCacheVal, exist := localCache.CM.Get(key)
	if exist && localCacheVal != nil {
		msg, _ := localCacheVal.(bo.Msg)
		rsp.Data = append(rsp.Data, msg)
		flag = true
		from = "local cache"
	}

	if !flag {
		val, err := redisClient.Get(context.TODO(), key).Result()
		if err == nil {
			var msg bo.Msg
			err = json.Unmarshal([]byte(val), &msg)
			if err == nil {
				rsp.Data = append(rsp.Data, msg)
				flag = true
				from = "redis"
			}
		}
	}

	if !flag {
		msg, err := model.GetLastMsg(uid)
		if err == nil {
			rsp.Data = append(rsp.Data, msg)
			flag = true
			from = "db"

			// 回写到local cache 和 Redis
			localCache.CM.Add(key, msg, 60*time.Second)
			val, err := json.Marshal(msg)
			if err == nil {
				redisClient.SetEX(context.TODO(), key, val, time.Second*60)
			}
		}
	}

	if !flag {
		ctx.JSON(http.StatusOK, QueryRsp{Code: 1, Message: "user last action msg not exist"})
		return
	}

	fmt.Printf("get user last action success, uid: %d, from: %s", uid, from)

	ctx.JSON(http.StatusOK, rsp)
}
