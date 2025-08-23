package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	_ "github.com/qq754174349/ht/ht-frame/grpc/service"
	log "github.com/qq754174349/ht/ht-frame/logger"
	"github.com/qq754174349/ht/ht-frame/web"
	"github.com/qq754174349/ht/ht-frame/web/router"
	"github.com/qq754174349/ht/user/internal/interface/http"
)

func main() {
	flag.Parse()

	err := web.DefaultRun(func(c *gin.Engine) {
		http.RegisterRoutes(router.New(c))
	})

	if err != nil {
		log.Fatal(err)
		return
	}
}
