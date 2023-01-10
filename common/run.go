package common

import (
	"context"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(r *gin.Engine, srvName string, addr string, stop func()) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if srvName == "project-search" {
		r.Static("/static", "./static/")
		r.LoadHTMLGlob("templates/*")
		r.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "index",
			})
		})
		r.GET("/news", func(c *gin.Context) {
			c.HTML(http.StatusOK, "news.html", gin.H{
				"title": "news",
			})
		})
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	// graceful restart
	go func() {
		log.Printf("%s running in  %s \n", srvName, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting Down project %s ...", srvName)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if stop != nil {
		stop()
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("%s Shutdown , cause by :", err, srvName)

	}
	select {
	case <-ctx.Done():
		log.Println("wait timeout ...")
	}
	log.Printf("%s stop success ...\n", srvName)
}
