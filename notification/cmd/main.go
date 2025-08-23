package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/qq754174349/ht/ht-frame/logger"
	"github.com/qq754174349/ht/ht-frame/web/router"
	"github.com/qq754174349/ht/notification/internal/interface/grpc"
	"github.com/qq754174349/ht/notification/internal/interface/http"

	"github.com/qq754174349/ht/ht-frame/web"
)

func main() {
	flag.Parse()

	// 注册 grpc
	grpc.Register()

	err := web.DefaultRun(func(c *gin.Engine) {
		http.RegisterRoutes(router.New(c))
	})

	if err != nil {
		logger.Fatal(err)
		return
	}
}
