package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-blog/models"
	"gin-blog/pkg/gredis"
	"gin-blog/pkg/logging"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"gin-blog/routers"
)

func init() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()
	util.Setup()
}

func main() {
	/*endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeOut
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeOut
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HTTPPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}*/

	router := routers.InitRouter()
	// 通过setting.RunMode设置运行模式
	gin.SetMode(setting.ServerSetting.RunMode)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		// 请求的读取操作在超时前的最大持续时间
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		// 回复的写入操作在超时前的最大持续时间
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		// 请求的头域最大长度，如为0则用DefaultMaxHeaderBytes 512M
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}