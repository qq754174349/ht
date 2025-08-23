package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/qq754174349/ht/gateway/internal/infrastructure/http/middleware"
	"github.com/qq754174349/ht/gateway/internal/interface/http"
	"github.com/qq754174349/ht/ht-frame/logger"
	"github.com/qq754174349/ht/ht-frame/web"
	"github.com/qq754174349/ht/ht-frame/web/router"
)

func main() {
	flag.Parse()
	err := web.DefaultRun(func(c *gin.Engine) {
		c.Use(middleware.PreflightMiddleware())
		http.RegisterRoutes(router.New(c))
	})

	if err != nil {
		logger.Fatal(err)
		return
	}
}
