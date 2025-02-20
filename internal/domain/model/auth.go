package model

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email,min=10,max=255"`
	Name        string `json:"name" binding:"required,min=5,max=255"`
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10,max=255"`
	Password    string `json:"password" binding:"required,min=8,max=255"`
}

type LoginRequest struct {
	Email          string `json:"email" binding:"required,email,min=10,max=255"`
	Password       string `json:"password" binding:"required,min=8,max=255"`
	RecaptchaToken string `json:"recaptchaToken" binding:"required"`
}

type LoginResponse struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	UserId       int64  `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
