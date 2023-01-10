package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"jam3.com/common"
	"jam3.com/user/config"
	"jam3.com/user/pgk/dao"
	"jam3.com/user/pgk/model"
	"net/http"
	"strings"
	"time"
)

type MyClaims struct {
	Uid int64 `json:"uid"`
	jwt.StandardClaims
}

// 定义过期时间
const TokenExpireDuration = time.Hour * 2

//定义secret
var MySecret = []byte("token-secret")

//生成jwt
func GenToken(uid int64) (string, error) {
	c := MyClaims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    config.C.SC.Name,
		},
	}
	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	//使用指定的secret签名并获得完成的编码后的字符串token
	return token.SignedString(MySecret)
}

//解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

//基于JWT认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中的auth为空",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)

		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中的auth格式错误",
			})
			//阻止调用后续的函数
			c.Abort()
			return
		}
		mc, err := ParseToken(parts[1])
		fmt.Println(mc)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的token",
			})
			c.Abort()
			return
		}
		//将当前请求的username信息保存到请求的上下文c上
		c.Set("uid", mc.Uid)
		//后续的处理函数可以通过c.Get("username")来获取请求的用户信息
		c.Next()
	}

}

func AuthHandler(c *gin.Context) {
	resp := &common.Result{}
	username, _ := c.GetQuery("username")
	password, _ := c.GetQuery("password")
	if username == "" || password == "" {
		c.JSON(http.StatusOK, resp.Fail(common.BusinessCode(model.UsernameOrPwd.Code), model.UsernameOrPwd.Msg))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var user *dao.User
	user = dao.Login(ctx, username, password)
	if user.Token != "" {
		//生成token
		tokenString, _ := GenToken(user.Uid)
		newUser := dao.User{
			Token: tokenString,
		}
		dao.Update(ctx, user.Uid, newUser)
		data := struct {
			Token string `json:"token"`
		}{Token: tokenString}
		c.JSON(
			http.StatusOK,
			resp.Success(data),
		)

		return
	}

	c.JSON(http.StatusOK, resp.Fail(common.BusinessCode(model.JwtAuthFail.Code), model.JwtAuthFail.Msg))
	return
}
