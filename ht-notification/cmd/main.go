package main

import (
	"flag"

	"github.com/qq754174349/ht/ht-notification/internal/interface/grpc"
	"github.com/qq754174349/ht/ht-notification/internal/interface/http"

	"github.com/qq754174349/ht/ht-frame/autoconfigure"
	_ "github.com/qq754174349/ht/ht-frame/orm/mysql"
	_ "github.com/qq754174349/ht/ht-frame/redis"
	"github.com/qq754174349/ht/ht-frame/web"
)

func main() {
	active := flag.String("active", "", "指定配置环境，例如 dev、prod 等")
	flag.Parse()

	// 注册 grpc
	grpc.Register()

	autoconfigure.Bootstrap(*active)
	err := web.Run(http.RegisterRoutes)

	if err != nil {
		return
	}
}
