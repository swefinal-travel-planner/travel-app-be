package model

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,min=10,max=255"`
	Name     string `json:"name" binding:"required,min=5,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
	OTP      string `json:"otp" binding:"required,min=6,max=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,min=10,max=255"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}

type GoogleLoginRequest struct {
	DisplayName string  `json:"displayName" binding:"required,min=5,max=255"`
	Email       string  `json:"email" binding:"required,email,min=10,max=255"`
	PhoneNumber string  `json:"phoneNumber" binding:"max=255"`
	PhotoURL    *string `json:"photoURL" binding:"required,min=10,max=255"`
	IDToken     *string `json:"id_token" binding:"required,min=10"`
	Password    string  `json:"password" binding:"required,min=8,max=255"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LoginResponse struct {
	Email        string `json:"email" binding:"required,email,min=10,max=255"`
	Name         string `json:"name" binding:"required,min=5,max=255"`
	UserId       int64  `json:"userId" binding:"required"`
	AccessToken  string `json:"accessToken" binding:"required"`
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email,min=10,max=255"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email,min=10,max=255"`
	OTP   string `json:"otp" binding:"required,min=6,max=6"`
}

type SetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email,min=10,max=255"`
	OTP      string `json:"otp" binding:"required,min=6,max=6"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}
