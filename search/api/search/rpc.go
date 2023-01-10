package search

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"jam3.com/common/discovery"
	"jam3.com/common/logs"
	"jam3.com/search/config"
	searchServiceV1 "jam3.com/search/pgk/service/search.service.v1"
	"log"
)

var SearchServiceClient searchServiceV1.SearchServiceClient

func InitRpcSearchClient() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.Dial("etcd:///search", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect:%v", err)
	}
	SearchServiceClient = searchServiceV1.NewSearchServiceClient(conn)
}
