package result

import (
	"context"
	"encoding/json"
)

type Result[T any] struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    T      `json:"data"`
	TraceId string `json:"traceId"`
}

func NewResult[T any](ctx context.Context, code int, msg string, data T) *Result[T] {
	traceId := (ctx).Value("traceID")
	if traceId == nil {
		traceId = ""
	}
	return &Result[T]{Code: code, Msg: msg, TraceId: traceId.(string), Data: data}
}

func NewBaseSuccessResult(ctx context.Context) *Result[struct{}] {
	traceId := (ctx).Value("traceID")
	if traceId == nil {
		traceId = ""
	}
	return NewTemplateResult(ctx, SUCCESS)
}

func NewSuccessResult[T any](ctx context.Context, data T) *Result[T] {
	return NewResult[T](ctx, SUCCESS.Code, SUCCESS.Msg, data)
}

func NewFailResult(ctx context.Context, msg string) *Result[struct{}] {
	return NewResult[struct{}](ctx, FAILURE.Code, msg, struct{}{})
}

func NewTemplateResult(ctx context.Context, template Template) *Result[struct{}] {
	return NewResult[struct{}](ctx, template.Code, template.Msg, struct{}{})
}

func (Result *Result[T]) ToString() string {
	marshal, _ := json.Marshal(Result)
	return string(marshal)
}
