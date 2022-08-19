package models

type LoginDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDto struct {
	AccessToken string `json:"accessToken"`
}
