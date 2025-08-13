package req

type EmailRegReq struct {
	// 邮箱
	Email string `json:"email"`
	// 用户名
	Username string `json:"username"`
	// 密码
	Password string `json:"password"`
}
