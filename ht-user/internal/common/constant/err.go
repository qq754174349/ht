package constant

import "github.com/qq754174349/ht/ht-frame/common/result"

var (
	NoLog     = result.Template{Code: 600, Msg: "未登录"}
	NoReg     = result.Template{Code: 610, Msg: "未注册"}
	RepeatReg = result.Template{Code: 611, Msg: "重复注册"}
)
