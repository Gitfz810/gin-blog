package main

import (
	"fmt"
	"log"
	"syscall"

	"github.com/fvbock/endless"

	"gin-blog/pkg/setting"
	"gin-blog/routers"
)

func main() {
	endless.DefaultReadTimeOut = setting.ReadTimeOut
	endless.DefaultWriteTimeOut = setting.WriteTimeOut
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.HTTPPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}


	/*router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		// 请求的读取操作在超时前的最大持续时间
		ReadTimeout:    setting.ReadTimeOut,
		// 回复的写入操作在超时前的最大持续时间
		WriteTimeout:   setting.WriteTimeOut,
		// 请求的头域最大长度，如为0则用DefaultMaxHeaderBytes 512M
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()*/
}