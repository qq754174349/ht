package client

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/qq754174349/ht-frame/consul"
	log "github.com/qq754174349/ht-frame/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetConn(serviceName string) *grpc.ClientConn {
	conn, err := grpc.NewClient(fmt.Sprintf("consul://%s/%s", consul.GetConfig().Consul.Addr, serviceName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
	)
	if err != nil {
		log.Errorf("获取grpc服务连接失败, %s", err)
		return nil
	}

	return conn
}
