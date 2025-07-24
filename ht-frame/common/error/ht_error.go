package error

import (
	"context"
	"encoding/json"
	result2 "github.com/qq754174349/ht-frame/common/result"
)

type HtError struct {
	Code    int
	Msg     string
	context context.Context
}

func (e *HtError) Error() string {
	res, _ := json.Marshal(result2.NewResult(e.context, e.Code, e.Msg, nil))
	return string(res)
}

func NewHtError(ctx context.Context, code int, msg string) *HtError {
	return &HtError{Code: code, Msg: msg, context: ctx}
}

func NewHtErrorFromMsg(ctx context.Context, msg string) *HtError {
	return &HtError{Code: result2.FAILURE.Code, Msg: msg, context: ctx}
}

func NewHtErrorFromTemplate(ctx context.Context, template result2.Template) *HtError {
	return &HtError{Code: template.Code, Msg: template.Msg, context: ctx}
}
