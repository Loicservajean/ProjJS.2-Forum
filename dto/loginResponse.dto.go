package dto

type LoginResponseDto struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
