package routes

import (
	"github.com/qq754174349/ht-frame/web/router"
	"ht-notification/internal/controller/mail"
)

func RegisterRoutes(router *router.Registrar) {
	apiGroup := router.Group("/api")
	{
		mailGroup := apiGroup.Group("/mail")
		{
			mailGroup.GET("/send", mail.Send)
		}
	}
}
