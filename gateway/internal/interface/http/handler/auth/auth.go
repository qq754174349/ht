package auth

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qq754174349/ht/gateway/internal/interface/dto/req"
	"github.com/qq754174349/ht/gateway/internal/usecase/auth"
	"github.com/qq754174349/ht/ht-frame/common/result"
)

var authUseCase = auth.NewUseCase()

func Auth(ctx *gin.Context) {
	accessToken := ctx.GetHeader("X-Access-Token")
	path := ctx.GetHeader("X-Forwarded-Uri")
	resp, err := authUseCase.Auth(ctx, &req.AuthReq{AccessToken: accessToken, Path: path})
	if err != nil {
		result.FailWithHttpCode(ctx, 401, 401, "token invalid or expired")
		return
	}
	if resp == nil {
		ctx.Status(200)
		return
	}

	if resp.UserId != 0 {
		ctx.Header("X-User-Id", strconv.FormatInt(resp.UserId, 10))
	}
	if resp.NewAccessToken != "" {
		ctx.Header("X-New-Access-Token", resp.NewAccessToken)
	}
}
