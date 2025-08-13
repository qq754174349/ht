package base

import (
	"time"

	"github.com/qq754174349/ht/ht-frame/orm/mysql"
)

type UserBaseInfo struct {
	mysql.Model
	//用户名
	Username *string `gorm:"type:varchar(16)"`
	// 邮箱
	Email *string `gorm:"type:varchar(255)"`
	// 国家码
	CountryCode *string `gorm:"type:varchar(10)"`
	// 手机号
	Phone *string `gorm:"type:varchar(20)"`
	// 密码
	Password *string `gorm:"type:varchar(30)"`
	// 昵称
	NickName *string `gorm:"type:varchar(30)"`
	// 头像
	AvatarUrl *string `gorm:"type:varchar(500)"`
	// 状态 0：待激活 1：已激活
	Status int8 `gorm:"type:tinyint(1)"`
	// 激活时间
	ActivateTime *time.Time `gorm:"type:datetime"`
}
