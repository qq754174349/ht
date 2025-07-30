package service

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/qq754174349/ht/ht-frame/autoconfigure"
	"github.com/qq754174349/ht/ht-frame/common/utils"
	baseConfig "github.com/qq754174349/ht/ht-frame/config"
	"github.com/qq754174349/ht/ht-frame/consul"
	grpcCfg "github.com/qq754174349/ht/ht-frame/grpc/config"
	"github.com/qq754174349/ht/ht-frame/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	grpcServer  *grpc.Server
	registrants []Registrant
)

type Registrant interface {
	Register(server *grpc.Server)
}

func RegisterRegistrant(r Registrant) {
	registrants = append(registrants, r)
}

func ApplyRegistrants() {
	for _, r := range registrants {
		r.Register(grpcServer)
	}
}

type AutoConfig struct{}

type AutoReg struct {
}

func init() {
	err := autoconfigure.Register(AutoConfig{})
	if err != nil {
		logger.Fatal("grpc 自动配置注册失败")
	}
}

func (AutoConfig) Init() error {
	appCig := baseConfig.GetAppCfg()
	config := grpcCfg.Get()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Grpc.Port))
	if err != nil {
		log.Fatalf("端口监听失败: %v", err)
	}

	grpcServer = grpc.NewServer()

	// 健康检查
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	ApplyRegistrants()

	// 启动服务
	go func() {
		log.Printf("gRPC服务启动，监听端口: %d", config.Grpc.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 确保服务启动后注册
	time.Sleep(100 * time.Millisecond)
	serviceID := fmt.Sprintf("%s-grpc-%s-%d", appCig.AppName, utils.GetOutboundIP(), config.Grpc.Port)
	if err := consul.Register(&api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    appCig.AppName + "-grpc",
		Address: utils.GetOutboundIP(),
		Port:    config.Grpc.Port,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", utils.GetOutboundIP(), config.Grpc.Port),
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30m",
		},
		Tags: []string{"grpc", "v1"},
	}); err != nil {
		log.Fatalf("Consul注册失败: %v", err)
	}

	// 优雅停止
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("grpc接收到终止信号，开始优雅停止...")

		// 1. 标记不健康
		healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

		// 2. 取消注册(带重试)
		consul.Deregister(serviceID)

		// 3. 停止服务
		done := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(done)
		}()

		// 超时控制
		select {
		case <-done:
			log.Println("grpc服务已优雅停止")
		case <-time.After(10 * time.Second):
			log.Println("grpc优雅停止超时，强制终止")
			grpcServer.Stop()
		}
	}()

	return nil
}

func (AutoConfig) Close() error {
	return nil
}

func GetServer() *grpc.Server {
	return grpcServer
}
