package http

import (
	"github.com/qq754174349/ht/ht-frame/web/router"
	"github.com/qq754174349/ht/notification/internal/interface/http/handler"
)

func RegisterRoutes(router *router.Registrar) {
	apiGroup := router.Group("/api")
	{
		mailGroup := apiGroup.Group("/mail")
		{
			mailGroup.GET("/SendTextMail", handler.SendTextMail)
		}
	}
}
