package http

import (
	"github.com/qq754174349/ht/gateway/internal/interface/http/handler/auth"
	"github.com/qq754174349/ht/ht-frame/web/router"
)

func RegisterRoutes(router *router.Registrar) {
	router.GET("/auth", auth.Auth)

}
