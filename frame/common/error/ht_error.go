package error

import (
	"encoding/json"
)

type HtError struct {
	Code int
	Msg  string
}

func (e *HtError) Error() string {
	res, _ := json.Marshal(e)
	return string(res)
}

func NewHtError(code int, msg string) *HtError {
	return &HtError{Code: code, Msg: msg}
}

func NewHtErrorFromMsg(msg string) *HtError {
	return &HtError{Code: FAILURE.Code, Msg: msg}
}

func NewBaseHtError() *HtError {
	return &HtError{Code: FAILURE.Code, Msg: FAILURE.Msg}
}

func NewHtErrorFromTemplate(template Template) *HtError {
	return &HtError{Code: template.Code, Msg: template.Msg}
}
