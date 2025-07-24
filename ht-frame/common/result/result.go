package result

import (
	"context"
	"encoding/json"
)

type Result struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"traceId"`
}

func NewResult(ctx context.Context, code int, msg string, data interface{}) *Result {
	traceId := (ctx).Value("traceID")
	if traceId == nil {
		traceId = ""
	}
	return &Result{Code: code, Msg: msg, TraceId: traceId.(string), Data: data}
}

func NewBaseSuccessResult(ctx context.Context) *Result {
	traceId := (ctx).Value("traceID")
	if traceId == nil {
		traceId = ""
	}
	return NewTemplateResult(ctx, SUCCESS)
}

func NewSuccessResult(ctx context.Context, data interface{}) *Result {
	return NewResult(ctx, SUCCESS.Code, SUCCESS.Msg, data)
}

func NewFailResult(ctx context.Context, msg string) *Result {
	return NewResult(ctx, FAILURE.Code, msg, nil)
}

func NewTemplateResult(ctx context.Context, template Template) *Result {
	return NewResult(ctx, template.Code, template.Msg, nil)
}

func (Result *Result) ToString() string {
	marshal, _ := json.Marshal(Result)
	return string(marshal)
}
