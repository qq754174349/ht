package req

type SessionCreateReq struct {
	Keyword     string `json:"keyword" binding:"required,min=6,max=30"`
	Password    string `json:"password" binding:"required,min=12,max=128"`
	AutoRenewal bool   `json:"autoRenewal"`
}
