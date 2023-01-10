package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"jam3.com/common"
	"jam3.com/common/errs"
	"jam3.com/user/pgk/dao"
	"jam3.com/user/pgk/model"
	"jam3.com/user/pgk/repo"
	loginServiceV1 "jam3.com/user/pgk/service/login.service.v1"
	"net/http"
	"strconv"
	"time"
)

// gin-swagger middleware
// swagger embed files

type HandlerUser struct {
	cache repo.Cache
}

func New() *HandlerUser {
	return &HandlerUser{
		cache: dao.Rc,
	}
}

// @用户信息
// @Description getuserinfo
// @host localhost
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /user/info [get]
func (h *HandlerUser) UserInfo(ctx *gin.Context) {

	resp := &common.Result{}
	zap.L().Info("Get param")
	//1.get param
	uid, _ := ctx.GetQuery("uid")
	userid, _ := strconv.Atoi(uid)
	zap.L().Debug("Get param uid : " + uid)
	//2.check param
	if userid < 0 {
		zap.L().Error("uid is unsiged int")
		ctx.JSON(http.StatusOK, resp.Fail(common.BusinessCode(model.NoLegal.Code), "uid is unsiged int"))
		return
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	userResponse, err := LoginServiceClient.GetUserinfo(c, &loginServiceV1.UserMessage{Uid: uid})
	if err != nil {
		//fromError, _ := status.FromError(err)
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, resp.Fail(code, msg))
	}
	ctx.JSON(http.StatusOK, resp.Success(userResponse.Username))
}

// @登录
// @Description login
// @host localhost
// @Param username path string true "用户名"
// @Param password path string true "密码"
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /user/login [get]
func (h *HandlerUser) Login(ctx *gin.Context) {
	resp := &common.Result{}
	username, _ := ctx.GetQuery("username")
	password, _ := ctx.GetQuery("password")
	if username == "" || password == "" {
		zap.L().Error("username or password is null")
		ctx.JSON(http.StatusOK, resp.Fail(common.BusinessCode(model.UsernameOrPwd.Code), model.UsernameOrPwd.Msg))
		return
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	loginResponse, err := LoginServiceClient.Login(c,
		&loginServiceV1.LoginMessage{Username: username, Password: password})
	if err != nil {
		//fromError, _ := status.FromError(err)
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, resp.Fail(code, msg))
		return
	}
	if loginResponse.Token == "" {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, resp.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, resp.Success(loginResponse.Token))
}
