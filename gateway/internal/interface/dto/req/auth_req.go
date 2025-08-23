package req

type AuthReq struct {
	AccessToken string `json:"accessToken"`
	Path        string `json:"path"`
}
