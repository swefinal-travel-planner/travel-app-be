package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type AuthService interface {
	Register(ctx *gin.Context, userRequest model.RegisterRequest) string
	RefreshToken(ctx *gin.Context, refreshRequest model.RefreshTokenRequest) (newAccessToken string, errCode string)
	Login(ctx *gin.Context, userRequest model.LoginRequest) (*model.LoginResponse, string)
	GoogleLogin(ctx *gin.Context, userRequest model.GoogleLoginRequest) (*model.LoginResponse, string)
	Logout(ctx *gin.Context, userId int64) string

	SendOTPToEmailForRegister(ctx *gin.Context, sendOTPRequest model.SendOTPRequest) string
	VerifyOTPForRegister(ctx *gin.Context, verifyOTPRequest model.VerifyOTPRequest) string
	SendOTPToEmailForResetPassword(ctx *gin.Context, sendOTPRequest model.SendOTPRequest) string
	VerifyOTPForResetPassword(ctx *gin.Context, verifyOTPRequest model.VerifyOTPRequest) string
	SetPassword(ctx *gin.Context, setPasswordRequest model.SetPasswordRequest) string
}
