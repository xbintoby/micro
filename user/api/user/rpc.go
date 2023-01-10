package user

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"jam3.com/common/discovery"
	"jam3.com/common/logs"
	"jam3.com/user/config"
	loginServiceV1 "jam3.com/user/pgk/service/login.service.v1"
	"log"
)

var LoginServiceClient loginServiceV1.UserServiceClient

func InitRpcUserClient() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.Dial("etcd:///user", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect:%v", err)
	}
	LoginServiceClient = loginServiceV1.NewUserServiceClient(conn)
}
