package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type AuthService interface {
	Register(ctx *gin.Context, userRequest model.RegisterRequest) error
	Login(ctx *gin.Context, userRequest model.LoginRequest) (*model.LoginResponse, error)
	GoogleLogin(ctx *gin.Context, userRequest model.GoogleLoginRequest) (*model.LoginResponse, error)
	ValidateRefreshToken(ctx *gin.Context, userId int64) (*entity.Authentication, error)

	SendOTPToEmailForRegister(ctx *gin.Context, sendOTPRequest model.SendOTPRequest) error
	VerifyOTPForRegister(ctx *gin.Context, verifyOTPRequest model.VerifyOTPRequest) error
	SendOTPToEmailForResetPassword(ctx *gin.Context, sendOTPRequest model.SendOTPRequest) error
	VerifyOTPForResetPassword(ctx *gin.Context, verifyOTPRequest model.VerifyOTPRequest) error
	SetPassword(ctx *gin.Context, setPasswordRequest model.SetPasswordRequest) error
}
