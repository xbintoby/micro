package search

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

type HandlerStar struct {
	cache repo.Cache
}

func NewStar() *HandlerStar {
	return &HandlerStar{
		cache: dao.Rc,
	}
}

// @人名
// @Description login
// @host localhost
// @Param term path string true "用户名"
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /star [get]
func (h *HandlerStar) Star(c *gin.Context) {
	resp := &common.Result{}
	zap.L().Info("Get param")
	term := c.DefaultQuery("term", "")
	if term == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stars, err := SearchServiceClient.StarQuery(ctx, &seachServiceV1.StarMessage{Term: term})
	fmt.Println(stars)
	if err != nil {
		//fromError, _ := status.FromError(err)
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, resp.Fail(code, msg))
		return
	}

	c.JSON(http.StatusOK, stars.Stars)
}
