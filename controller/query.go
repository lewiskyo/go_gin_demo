package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"go_gin_demo/bo"
	"go_gin_demo/cache"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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

	redisClient := cache.Redis()
	if redisClient == nil {
		log.Printf("get redis instance fail")
		ctx.JSON(http.StatusOK, QueryRsp{Code: 1, Message: "get redis instance fail"})
		return
	}

	uid := req.Uid
	key := fmt.Sprintf("uid:%d", uid)

	// 先从localcache获取
	localCache, _ := cache.LocalCache()
	localCacheVal, exist := localCache.CM.Get(key)
	if exist && localCacheVal != nil {
		log.Println("get from local cache")
		rsp := QueryRsp{Code: 0, Message: "ok"}
		msg, _ := localCacheVal.(bo.Msg)
		rsp.Data = append(rsp.Data, msg)
		ctx.JSON(http.StatusOK, rsp)
		return
	}

	val, err := redisClient.Get(context.TODO(), key).Result()
	if err == redis.Nil {
		log.Println("redis val not exist")
		ctx.JSON(http.StatusOK, QueryRsp{Code: 1, Message: "last access info not exist"})
		return
	}

	if err != nil {
		log.Printf("get redis val fail, err: %s", err.Error())
		ctx.JSON(http.StatusOK, QueryRsp{Code: 1, Message: "get redis val fail"})
		return
	}

	var msg bo.Msg
	err = json.Unmarshal([]byte(val), &msg)
	if err != nil {
		log.Printf("user last access info err")
		ctx.JSON(http.StatusOK, QueryRsp{Code: 1, Message: "user last access info err"})
		return
	}

	rsp := QueryRsp{Code: 0, Message: "ok"}
	rsp.Data = append(rsp.Data, msg)

	ctx.JSON(http.StatusOK, rsp)
}
