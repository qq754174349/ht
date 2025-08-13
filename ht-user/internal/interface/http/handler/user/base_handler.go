package user

import (
	req3 "github.com/qq754174349/ht/ht-user/internal/interface/dto/req"
	"github.com/qq754174349/ht/ht-user/internal/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/qq754174349/ht/ht-frame/common/result"
)

var userUseCase = user.NewUserUseCase()

// WechatUserLogin 微信小程序登录
func WechatUserLogin(ctx *gin.Context) {
	code := ctx.Query("code")
	jwt, err := userUseCase.WechatLogin(ctx, code)
	if err != nil {
		ctx.Writer.WriteString(err.Error())
	} else {
		ctx.Writer.WriteString(result.NewSuccessResult(ctx, jwt).ToString())
	}
}

// WechatUserReg 微信小程序注册
func WechatUserReg(ctx *gin.Context) {
	req := req3.WechatRegReq{}
	ctx.BindJSON(&req)
	userUseCase.WechatReg(ctx, req)
}

// MailReg 邮箱注册
func MailReg(ctx *gin.Context) {
	req := req3.EmailRegReq{}
	ctx.ShouldBindJSON(&req)

	err := userUseCase.EMailReg(ctx, req)
	if err != nil {
		ctx.Writer.WriteString(err.Error())
		return
	}

	ctx.Writer.WriteString(result.NewBaseSuccessResult(ctx).ToString())
}

// Activate 用户激活
func Activate(ctx *gin.Context) {
	token := ctx.Query("token")
	err := userUseCase.UserActivate(ctx, token)
	if err != nil {
		ctx.Writer.WriteString(err.Error())
	}

	ctx.Writer.WriteString(result.NewBaseSuccessResult(ctx).ToString())
}
