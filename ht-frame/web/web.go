package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qq754174349/ht-frame/autoconfigure"
	"github.com/qq754174349/ht-frame/common/utils"
	baseConfig "github.com/qq754174349/ht-frame/config"
	"github.com/qq754174349/ht-frame/consul"
	_ "github.com/qq754174349/ht-frame/grpc/service"
	"github.com/qq754174349/ht-frame/logger"
	"github.com/qq754174349/ht-frame/web/middlewares"
	"github.com/qq754174349/ht-frame/web/router"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Engine struct {
	gin.Engine
	timeout time.Duration
}

type RouteGroup struct {
	gin.RouterGroup
}

type AutoConfig struct{}

var config *Web

type Web struct {
	Web Config
}

type Config struct {
	Port    string
	Timeout time.Duration
}

func init() {
	err := autoconfigure.Register(AutoConfig{})
	if err != nil {
		logger.Fatal("web 自动配置注册失败")
	}
}

func (AutoConfig) Init() error {
	config = &Web{}
	autoconfigure.ConfigRead(config)
	appCig := baseConfig.GetAppCfg()
	gin.DefaultWriter = logger.Writer()
	gin.DefaultErrorWriter = logger.Writer()
	if appCig.Active == "pro" {
		gin.SetMode(gin.ReleaseMode)
	} else if appCig.Active == "test" {
		gin.SetMode(gin.TestMode)
	}
	return nil
}

func Default(opts ...gin.OptionFunc) *gin.Engine {
	engine := gin.New()
	engine.Use(
		middlewares.GenerateTraceID(),
		middlewares.RequestInfoLogger([]string{"/health"}),
		middlewares.Prometheus(),
		middlewares.ErrorHandler(),
		gin.Recovery(),
	)
	engine.With(opts...)
	return engine
}

func Run(regRoutes func(engine *router.Registrar), opts ...gin.OptionFunc) error {
	engine := Default(opts...)
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	regRoutes(router.New(engine, config.Web.Timeout))
	srv := &http.Server{
		Addr:    ":" + config.Web.Port,
		Handler: engine,
	}

	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return fmt.Errorf("监听端口失败: %w", err)
	}
	addr := listener.Addr().(*net.TCPAddr)
	logger.Infof("service start success, port: %d", addr.Port)
	go func() {
		if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("web服务器运行错误: %v", err)
		}
	}()

	// 注册consul
	//consulRegister()
	// 处理优雅关闭
	GracefulShutdown(srv)
	return nil
}

func GracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("web接收到关闭信号，开始优雅关闭...")
	consul.Deregister(baseConfig.GetAppCfg().AppName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("web服务器关闭错误: %v", err)
	}
	logger.Info("web服务已完全关闭")
}

func consulRegister() {
	port, _ := strconv.Atoi(config.Web.Port)
	err := consul.Register(&api.AgentServiceRegistration{
		ID:   baseConfig.GetAppCfg().AppName,
		Name: baseConfig.GetAppCfg().AppName,
		Port: port,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%s/health", utils.GetOutboundIP(), config.Web.Port),
			Interval: "20s",
			Timeout:  "3s",
		},
	})
	if err != nil {
		panic(err)
	}
}
