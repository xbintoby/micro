package main

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/gin-gonic/gin"
	srv "jam3.com/common"
	_ "jam3.com/user/api"
	"jam3.com/user/config"
	_ "jam3.com/user/docs"
	"jam3.com/user/router"
	"log"
)

func main() {
	err := sentinel.InitWithConfigFile("config/sentinel.yaml")
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()

	router.InitRouter(r)

	grpc := router.RegisterGrpc()
	router.RegisterEtcdServer()
	stop := func() {
		grpc.Stop()
	}
	srv.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
}
