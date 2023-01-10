package login_service_v1

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"jam3.com/common/errs"
	"jam3.com/user/pgk/dao"
	"jam3.com/user/pgk/model"
	"jam3.com/user/pgk/repo"
	"log"
	"strconv"
	"time"
)

type LoginService struct {
	UnimplementedUserServiceServer
	cache repo.Cache
	db    *gorm.DB
}

func New() *LoginService {
	return &LoginService{
		cache: dao.Rc,
	}
}
func (ls *LoginService) Login(ctx context.Context, msg *LoginMessage) (*LoginResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	username := msg.Username
	password := msg.Password
	newUser := dao.Login(c, username, password)
	if newUser.Token == "" {
		return &LoginResponse{Token: newUser.Token}, errs.GrpcError(model.TokenIsNull)
	}
	return &LoginResponse{Token: newUser.Token}, nil
}
func (ls *LoginService) GetUserinfo(ctx context.Context, msg *UserMessage) (*UserResponse, error) {
	zap.L().Info("Get param")
	//1.get param
	uid := msg.Uid
	userid, _ := strconv.Atoi(uid)
	zap.L().Debug("Get param" + uid)
	//2.check param
	if userid < 0 {
		zap.L().Error("uid is unsiged int")
		return nil, errs.GrpcError(model.NoLegalUid)
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	key := "user" + uid
	cVal, err := ls.cache.Get(ctx, key)
	if err != nil {
		log.Printf("get cache %s err: %s \n", key, err)
	}

	getFromDb := func() string {
		var user *dao.User
		user = dao.GetInfo(ctx, int64(userid))

		err := ls.cache.Put(c, key, user.Username, 15*time.Minute)
		if err != nil {
			log.Printf("put cache %s err: %s \n", key, err)
		}
		return user.Username
	}
	val := ""
	if cVal == "" {
		val = getFromDb()
	} else {
		val = cVal
	}

	return &UserResponse{Username: val}, nil
}
