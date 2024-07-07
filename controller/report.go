package controller

import (
	"fmt"
	"go_gin_demo/bo"
	"go_gin_demo/cache"
	"go_gin_demo/model.go"
	"go_gin_demo/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// ReportReq 行为日志上报请求体结构
type ReportReq struct {
	Timestamp int64  `json:"timestamp"`
	Uid       int64  `json:"uid"`
	Type      uint8  `json:"type"`
	Region    string `json:"region"`
	Device    string `json:"device"`
	Ip        string `json:"ip"`
	Network   string `json:"network"`
	Version   uint64 `json:"version"`
}

// ReportReq 行为日志上报响应体结构
type ReportRsp struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
}

func (ctl *BaseController) Report(ctx *gin.Context) {
	var req ReportReq
	if err := ctx.BindQuery(&req); err != nil {
		log.Printf("bind query err: %s", err.Error())
		ctx.JSON(http.StatusOK, ReportRsp{Code: 1, Message: "uri param error"})
		return
	}

	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		log.Printf("bind body err: %s", err.Error())
		ctx.JSON(http.StatusOK, ReportRsp{Code: 1, Message: "form param error"})
		return
	}

	msg := bo.Msg{
		Timestamp: req.Timestamp,
		Uid:       req.Uid,
		Type:      req.Type,
		Region:    req.Region,
		Device:    req.Device,
		Ip:        req.Ip,
		Network:   req.Network,
		Version:   req.Version,
	}

	utils.MsgQueue.Enqueue(msg)

	localCache := cache.GetLocalCache()
	key := fmt.Sprintf("%s:%d", bo.MSG_KEY_PREFIX, req.Uid)
	localCache.CM.Add(key, msg, 60*time.Second)
	model.InsertMsg(msg)

	ctx.JSON(http.StatusOK, ReportRsp{Code: 0, Message: "ok"})
}
