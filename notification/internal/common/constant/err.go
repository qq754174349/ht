package constant

import "github.com/qq754174349/ht-frame/common/error"

var (
	NoLog     = error.Template{Code: 600, Msg: "未登录"}
	NoReg     = error.Template{Code: 610, Msg: "未注册"}
	RepeatReg = error.Template{Code: 611, Msg: "重复注册"}
)
