package wechat

import (
	"github.com/qq754174349/ht/ht-frame/orm/mysql"
)

type UserWechatInfo struct {
	mysql.Model
	// 用户id
	UserId int64
	// 微信昵称
	NickName string `gorm:"type:varchar(50)"`
	// 微信头像
	AvatarUrl string `gorm:"type:varchar(255)"`
	// 微信openId
	OpenId *string `gorm:"type:varchar(32)"`
}
