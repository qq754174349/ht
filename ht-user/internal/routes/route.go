package routes

import (
	"github.com/qq754174349/ht-frame/web/router"
	"ht-user/internal/controller/user"
)

func RegisterRoutes(router *router.Registrar) {
	apiGroup := router.Group("/api")
	{
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/wechat/login", user.WechatUserLogin).
				POST("/wechat/reg", user.WechatUserReg).
				GET("/mail/reg", user.MailReg)
		}
	}
	authGroup := router.Group("api/auth")
	authGroup.POST("/wechat/reg", user.MailReg)
}
