package constant

import (
	error2 "github.com/qq754174349/ht/ht-frame/common/error"
)

var (
	NoLog          = error2.Template{Code: 600, Msg: "未登录"}
	NoReg          = error2.Template{Code: 610, Msg: "未注册"}
	RepeatReg      = error2.Template{Code: 611, Msg: "重复注册"}
	NoActivate     = error2.Template{Code: 620, Msg: "用户未激活"}
	ActivateExpire = error2.Template{Code: 621, Msg: "激活链接失效"}
	RepeatActivate = error2.Template{Code: 622, Msg: "重复激活"}
	NoUser         = error2.Template{Code: 630, Msg: "用户不存在"}
	UserNamePwdErr = error2.Template{Code: 631, Msg: "用户名或密码错误"}
)
