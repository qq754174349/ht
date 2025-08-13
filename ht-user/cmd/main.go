package main

import (
	"flag"

	"github.com/qq754174349/ht/ht-user/internal/interface/http"

	"github.com/qq754174349/ht/ht-frame/autoconfigure"
	"github.com/qq754174349/ht/ht-frame/web"
)

func main() {
	active := flag.String("active", "", "指定配置环境，例如 dev、prod 等")
	flag.Parse()

	autoconfigure.Bootstrap(*active)

	err := web.Run(http.RegisterRoutes)

	if err != nil {
		return
	}
}
