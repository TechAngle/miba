package domain

type LoginResponse struct {
	URL   string `json:"url"`
	Token string `json:"token"`
	Code  int    `json:"code"`
}
