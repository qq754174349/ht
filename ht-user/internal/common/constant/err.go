package constant

import "github.com/qq754174349/ht/ht-frame/common/result"

var (
	NoLog          = result.Template{Code: 600, Msg: "未登录"}
	NoReg          = result.Template{Code: 610, Msg: "未注册"}
	RepeatReg      = result.Template{Code: 611, Msg: "重复注册"}
	NoActivate     = result.Template{Code: 620, Msg: "未激活"}
	ActivateExpire = result.Template{Code: 621, Msg: "激活链接失效"}
	RepeatActivate = result.Template{Code: 622, Msg: "重复激活"}
	NoUser         = result.Template{Code: 630, Msg: "用户不存在"}
)
