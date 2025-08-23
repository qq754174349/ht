package web

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qq754174349/ht/ht-frame/autoconfigure"
	"github.com/qq754174349/ht/ht-frame/common/utils"
	baseConfig "github.com/qq754174349/ht/ht-frame/config"
	"github.com/qq754174349/ht/ht-frame/consul"
	"github.com/qq754174349/ht/ht-frame/logger"
	"github.com/qq754174349/ht/ht-frame/web/config"
	"github.com/qq754174349/ht/ht-frame/web/middlewares"
)

type Engine struct {
	gin.Engine
	timeout time.Duration
}

type RouteGroup struct {
	gin.RouterGroup
}

type AutoConfig struct{}

func init() {
	autoconfigure.Register(&AutoConfig{})
	appCig := baseConfig.GetAppCfg()
	gin.DefaultWriter = logger.Writer()
	gin.DefaultErrorWriter = logger.Writer()
	if appCig.Active == "pro" {
		gin.SetMode(gin.ReleaseMode)
	} else if appCig.Active == "test" {
		gin.SetMode(gin.TestMode)
	}
}

func (AutoConfig) Close() error {
	return nil
}

func Default(opts ...gin.OptionFunc) *gin.Engine {
	engine := gin.New()
	engine.Use(
		middlewares.TraceIDMiddleware(),
		middlewares.HeaderPropagationMiddleware("X-New-Access-Token", "X-Trace-Id"),
		middlewares.RequestInfoLogger([]string{"/health"}),
		middlewares.Prometheus(),
		middlewares.ErrorHandler(),
		middlewares.Cors(),
	)
	engine.With(opts...)
	return engine
}

func DefaultRun(opts ...gin.OptionFunc) error {
	engine := Default(opts...)
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	srv := &http.Server{
		Addr:    ":" + config.Get().Web.Port,
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
	consulRegister()
	// 处理优雅关闭
	GracefulShutdown(srv)

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("捕获到异常: %v", r)
		}
	}()
	return nil
}

func Run(opts ...gin.OptionFunc) error {
	engine := gin.New()
	engine.With(opts...)

	srv := &http.Server{
		Addr:    ":" + config.Get().Web.Port,
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
	autoconfigure.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("web服务器关闭错误: %v", err)
	}
	logger.Info("web服务已完全关闭")
}

func consulRegister() {
	port, _ := strconv.Atoi(config.Get().Web.Port)
	err := consul.Register(&api.AgentServiceRegistration{
		ID:   baseConfig.GetAppCfg().AppName,
		Name: baseConfig.GetAppCfg().AppName,
		Port: port,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%s/health", utils.GetOutboundIP(), config.Get().Web.Port),
			Interval: "20s",
			Timeout:  "3s",
		},
		Tags: []string{
			"traefik.enable=true",
			fmt.Sprintf("traefik.http.routers.%s.middlewares=strip-prefix@file,auth@file", baseConfig.GetAppCfg().AppName),
		},
	})
	if err != nil {
		panic(err)
	}
}
