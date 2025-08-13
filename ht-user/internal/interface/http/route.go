package http

import (
	"github.com/qq754174349/ht/ht-frame/web/router"
	"github.com/qq754174349/ht/ht-user/internal/interface/http/handler/user"
)

func RegisterRoutes(router *router.Registrar) {
	apiGroup := router.Group("/api")
	{
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/wechat/login", user.WechatUserLogin).
				POST("/wechat/reg", user.WechatUserReg).
				POST("/mail/reg", user.MailReg).
				GET("/activate", user.Activate)
		}
	}
	authGroup := router.Group("api/auth")
	authGroup.POST("/wechat/reg", user.MailReg)
}
