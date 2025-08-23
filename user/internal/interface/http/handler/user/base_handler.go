package user

import (
	"strconv"

	"github.com/qq754174349/ht/user/internal/interface/dto/req"
	"github.com/qq754174349/ht/user/internal/usecase/user"

	"github.com/gin-gonic/gin"
	"github.com/qq754174349/ht/ht-frame/common/result"
)

var userUseCase = user.NewUserUseCase()

// GetCurrentUser 获取用户信息
func GetCurrentUser(ctx *gin.Context) {
	userIdStr := ctx.GetHeader("X-User-Id")
	if userIdStr == "" {
		result.FailDefault(ctx)
		return
	}
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		result.FailDefault(ctx)
		return
	}

	resp, err := userUseCase.GetUserById(ctx, userId)
	if err != nil {
		result.FailByErr(ctx, err)
	}
	result.Success(ctx, resp)
}

// WechatUserLogin 微信小程序登录
func WechatUserLogin(ctx *gin.Context) {
	code := ctx.Query("code")
	jwt, err := userUseCase.WechatLogin(ctx, code)
	if err != nil {
		result.FailByErr(ctx, err)
	} else {
		result.Success(ctx, jwt)
	}
}

// WechatUserReg 微信小程序注册
func WechatUserReg(ctx *gin.Context) {
	body := req.WechatRegReq{}
	ctx.BindJSON(&body)
	userUseCase.WechatReg(ctx, body)
}

// MailReg 邮箱注册
func MailReg(ctx *gin.Context) {
	body := req.EmailRegReq{}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		result.FailByErr(ctx, err)
		return
	}
	err = userUseCase.EMailReg(ctx, body)
	if err != nil {
		ctx.Writer.WriteString(err.Error())
		return
	}

	result.SuccessEmpty(ctx)
}

// SessionCreate 创建用户会话
func SessionCreate(ctx *gin.Context) {
	body := req.SessionCreateReq{}
	ctx.ShouldBindJSON(&body)

	token, err := userUseCase.SessionCreate(ctx, body)
	if err != nil {
		result.FailByErr(ctx, err)
		return
	}

	ctx.Header("X-New-Access-Token", token)
	result.SuccessEmpty(ctx)
}

// Activate 用户激活
func Activate(ctx *gin.Context) {
	token := ctx.Query("token")
	err := userUseCase.UserActivate(ctx, token)
	if err != nil {
		result.FailByErr(ctx, err)
		return
	}

	result.SuccessEmpty(ctx)
}
