package search

import (
	"context"
	"github.com/gin-gonic/gin"
	"jam3.com/common"
	"jam3.com/common/errs"
	"jam3.com/search/pgk/dao"
	"jam3.com/search/pgk/repo"
	seachServiceV1 "jam3.com/search/pgk/service/search.service.v1"
	"net/http"
	"time"
)

// gin-swagger middleware
// swagger embed files

type HandlerNews struct {
	cache repo.Cache
}

func NewNews() *HandlerNews {
	return &HandlerNews{
		cache: dao.Rc,
	}
}

// @新闻标题
// @Description login
// @host localhost
// @Param term path string true "标题"
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /tips [get]
func (h *HandlerNews) tips(c *gin.Context) {
	resp := &common.Result{}

	term := c.DefaultQuery("term", "")
	if term == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	articles, err := SearchServiceClient.NewsQuery(ctx, &seachServiceV1.NewsMessage{Term: term})

	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, resp.Fail(code, msg))
		return
	}

	c.JSON(http.StatusOK, articles.Title)
}

// @相关新闻
// @Description login
// @host localhost
// @Param term path string true "标题"
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /tips [get]
func (h *HandlerNews) newsSearch(c *gin.Context) {
	resp := &common.Result{}

	text := c.DefaultQuery("text", "")
	if text == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	articles, err := SearchServiceClient.SearchNews(ctx, &seachServiceV1.ArticleMessage{Text: text})

	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, resp.Fail(code, msg))
		return
	}

	c.JSON(http.StatusOK, articles.Arts)
}
