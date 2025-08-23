package resp

type JwtValidateResp struct {
	UserId         int64  `json:"userId"`
	NewAccessToken string `json:"newAccessToken"`
}
