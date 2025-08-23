package req

type EmailRegReq struct {
	// 邮箱
	Email string `json:"email" binding:"required,email" label:"邮箱"`
	// 用户名
	Username string `json:"username" binding:"required,min=6,max=20" label:"用户名"`
	// 密码
	Password string `json:"password" binding:"required,min=12,max=128" label:"密码"`
}
