package user

import (
	"github.com/gin-gonic/gin"
	"github.com/qq754174349/ht-frame/common/result"
	req2 "ht-user/internal/dto/req"
	"ht-user/internal/service/user"
	"time"
)

// WechatUserLogin 微信小程序登录
func WechatUserLogin(ctx *gin.Context) {
	code := ctx.Query("code")

	jwt, err := user.WechatLogin(ctx, code)
	if err != nil {
		ctx.Writer.WriteString(err.Error())
	} else {
		ctx.Writer.WriteString(result.NewSuccessResult(ctx, jwt).ToString())
	}
}

// WechatUserReg 微信小程序注册
func WechatUserReg(ctx *gin.Context) {
	req := req2.WechatRegReq{}
	ctx.BindJSON(&req)
	user.WechatReg(ctx, req)
}

// MailReg 邮箱注册
func MailReg(ctx *gin.Context) {
	req := req2.EmailRegReq{}
	ctx.ShouldBindJSON(&req)

	time.Sleep(10 * time.Second)
	user.EMailReg(ctx, req)

	ctx.Writer.WriteString(result.NewBaseSuccessResult(ctx).ToString())
}
