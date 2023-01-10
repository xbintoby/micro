package middleware

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	// 配置一条限流规则
	_, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:        "/user/auth",
			Threshold:       1,
			ControlBehavior: flow.Reject,
		},
		{
			Resource:        "/user/login",
			Threshold:       1,
			ControlBehavior: flow.Reject,
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func QPS(context *gin.Context) {
	zap.L().Debug(context.Request.URL.Path)
	zap.L().Debug(context.Request.URL.Query().Get("username"))
	e, err := sentinel.Entry(context.Request.URL.Path)
	if err != nil {
		errStr := context.Request.URL.Path + " QPS limited "
		context.AbortWithStatusJSON(400, gin.H{"error": errStr})
	} else {
		e.Exit()
		context.Next()
	}
	return
}
