package result

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	error2 "github.com/qq754174349/ht/ht-frame/common/error"
)

type Result struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"traceId"`
}

func New(ctx context.Context, code int, msg string, data interface{}) *Result {
	traceId := (ctx).Value("traceID")
	if traceId == nil {
		traceId = ""
	}
	return &Result{Code: code, Msg: msg, TraceId: traceId.(string), Data: data}
}

func Success(ctx *gin.Context, data interface{}) {
	traceId := (ctx).Value("traceID")
	if traceId == nil {
		traceId = ""
	}
	ctx.JSON(http.StatusOK, New(ctx, error2.SUCCESS.Code, error2.SUCCESS.Msg, data))
}

func SuccessEmpty(ctx *gin.Context) {
	traceId := (ctx).Value("traceID")
	if traceId == nil {
		traceId = ""
	}
	ctx.JSON(http.StatusOK, New(ctx, error2.SUCCESS.Code, error2.SUCCESS.Msg, nil))
}

func FailDefault(ctx *gin.Context) {
	Write(ctx, http.StatusOK, error2.FAILURE.Code, error2.FAILURE.Msg, nil)
}

func Fail(ctx *gin.Context, code int, msg string) {
	Write(ctx, http.StatusOK, code, msg, nil)
}

func FailByTemplate(ctx *gin.Context, template error2.Template) {
	Write(ctx, http.StatusOK, template.Code, template.Msg, nil)
}

func FailByErr(ctx *gin.Context, err error) {
	var htError *error2.HtError
	if errors.As(err, &htError) {
		Fail(ctx, htError.Code, htError.Msg)
		return
	}

	FailDefault(ctx)
}

func FailWithHttpCode(ctx *gin.Context, httpCode int, code int, msg string) {
	Write(ctx, httpCode, code, msg, nil)
}

func Write(ctx *gin.Context, httpCode int, code int, msg string, data interface{}) {
	ctx.JSON(httpCode, New(ctx, code, msg, data))
}
