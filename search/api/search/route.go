package search

import (
	"github.com/gin-gonic/gin"
	"jam3.com/search/router"
	"log"
)

func init() {
	log.Println("init search router")
	router.Register(&RouterSearch{})
}

type RouterSearch struct {
}

func (t *RouterSearch) Route(r *gin.Engine) {
	star := NewStar()
	news := NewNews()
	InitRpcSearchClient()
	r.GET("/star", star.Star)
	r.GET("/news/tips", news.tips)
	r.GET("/news/search", news.newsSearch)
}
