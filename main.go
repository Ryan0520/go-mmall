package main

import (
	log "github.com/sirupsen/logrus"

	"fmt"
	"github.com/Ryan0520/go-mmall/models"
	"github.com/Ryan0520/go-mmall/pkg/gredis"
	"github.com/Ryan0520/go-mmall/pkg/setting"
	"github.com/Ryan0520/go-mmall/routers"
	"github.com/fvbock/endless"
	"syscall"
)

func main() {
	models.Setup()
	gredis.SetUp()

	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
