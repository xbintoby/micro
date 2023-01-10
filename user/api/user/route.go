package user

import (
	"github.com/gin-gonic/gin"
	"jam3.com/user/pgk/middleware"
	"jam3.com/user/router"
	"log"
)

func init() {
	log.Println("init user router")
	router.Register(&RouterUser{})
}

type RouterUser struct {
}

func (t *RouterUser) Route(r *gin.Engine) {
	h := New()
	InitRpcUserClient()

	// curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJoYW8iLCJleHAiOjE2NzI3NDg1NzAsImlzcyI6InByb2plY3QtdXNlciJ9.As7FTbDCieFaNDAkolo8tYCanqytUIDyVMStAOORfsQ" http://localhost/user/info?uid=3
	r.GET("/user/info", middleware.JWTAuthMiddleware(), h.UserInfo)
	r.GET("/user/login", middleware.QPS, h.Login)
	r.GET("/user/auth", middleware.QPS, middleware.AuthHandler)
}
