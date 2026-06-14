package dto

type RegisterRequestDto struct {
	Pseudo   string `json:"pseudo"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
