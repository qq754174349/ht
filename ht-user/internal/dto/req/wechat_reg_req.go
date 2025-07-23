package req

type WechatRegReq struct {
	// 登录凭证
	Code string `json:"code"`
	// 微信昵称
	Nickname string `json:"nickname"`
	// 微信头像
	AvatarUrl string `json:"avatar_url"`
}
