package main

import (
	"flag"
	"github.com/qq754174349/ht-frame/autoconfigure"
	"github.com/qq754174349/ht-frame/grpc/service"
	_ "github.com/qq754174349/ht-frame/mysql"
	_ "github.com/qq754174349/ht-frame/redis"
	"github.com/qq754174349/ht-frame/web"
	"ht-notification/internal/routes"
	"ht-notification/internal/service/mail"
)

func main() {
	active := flag.String("active", "", "指定配置环境，例如 dev、prod 等")
	flag.Parse()

	service.RegisterRegistrant(&mail.Register{})

	autoconfigure.Bootstrap(*active)
	err := web.Run(routes.RegisterRoutes)

	if err != nil {
		return
	}
}
