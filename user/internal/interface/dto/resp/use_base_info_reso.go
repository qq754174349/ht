package resp

type UseBaseInfoResp struct {
	// 用户id
	UserId int64
	// 昵称
	NickName *string `gorm:"type:varchar(30)"`
	// 头像
	AvatarUrl *string `gorm:"type:varchar(500)"`
}
