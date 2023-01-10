package main

import (
	"github.com/gin-gonic/gin"
	srv "jam3.com/common"
	_ "jam3.com/search/api"
	"jam3.com/search/config"
	_ "jam3.com/search/docs"
	"jam3.com/search/router"
)

func main() {
	r := gin.Default()

	router.InitRouter(r)

	grpc := router.RegisterGrpc()
	router.RegisterEtcdServer()
	stop := func() {
		grpc.Stop()
	}
	srv.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
}
