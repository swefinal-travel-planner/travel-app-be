package v1

import (
	"github.com/gin-gonic/gin"
	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/validation"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// @Summary Register user
// @Description User register
// @Tags Auths
// @Accept json
// @Param request body model.RegisterRequest true "Auth payload"
// @Produce  json
// @Router /auth/register [post]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) Register(ctx *gin.Context) {
	var registerRequest model.RegisterRequest

	if err := validation.BindJsonAndValidate(ctx, &registerRequest); err != nil {
		return
	}

	errCode := handler.authService.Register(ctx, registerRequest)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}
	ctx.AbortWithStatus(204)
}

// @Summary Login
// @Description Login to account
// @Tags Auths
// @Accept json
// @Param request body model.LoginRequest true "Auth payload"
// @Produce  json
// @Router /auth/login [post]
// @Success 200 {object} httpcommon.HttpResponse[entity.User]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) Login(ctx *gin.Context) {
	var loginRequest model.LoginRequest

	if err := validation.BindJsonAndValidate(ctx, &loginRequest); err != nil {
		return
	}

	user, errCode := handler.authService.Login(ctx, loginRequest)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(200, httpcommon.NewSuccessResponse(user))
}

// @Summary Refresh
// @Description Refresh token
// @Tags Auths
// @Accept json
// @Param request body model.RefreshTokenRequest true "Auth payload"
// @Produce  json
// @Router /auth/refresh [post]
// @Success 200 {object} httpcommon.HttpResponse[entity.User]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) Refresh(ctx *gin.Context) {
	var refreshTokenRequest model.RefreshTokenRequest

	if err := validation.BindJsonAndValidate(ctx, &refreshTokenRequest); err != nil {
		return
	}

	newAccessToken, errCode := handler.authService.RefreshToken(ctx, refreshTokenRequest)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}
	ctx.JSON(200, httpcommon.NewSuccessResponse(&newAccessToken))
}

// @Summary Send OTP to Mail for register
// @Description Send OTP to user email when registering
// @Tags Auths
// @Accept json
// @Param request body model.SendOTPRequest true "Send OTP payload"
// @Produce json
// @Router /auth/register/send-otp [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) SendOTPToEmailForRegister(ctx *gin.Context) {
	var sendOTPRequest model.SendOTPRequest

	if err := validation.BindJsonAndValidate(ctx, &sendOTPRequest); err != nil {
		return
	}

	errCode := handler.authService.SendOTPToEmailForRegister(ctx, sendOTPRequest)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(204)
}

// @Summary Verify OTP for register
// @Description Verify OTP with email and otp when registering
// @Tags Auths
// @Accept json
// @Param request body model.VerifyOTPRequest true "Verify OTP payload"
// @Produce json
// @Router /auth/register/verify-otp [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) VerifyOTPForRegister(ctx *gin.Context) {
	var verifyOTPRequest model.VerifyOTPRequest

	if err := validation.BindJsonAndValidate(ctx, &verifyOTPRequest); err != nil {
		return
	}

	errCode := handler.authService.VerifyOTPForRegister(ctx, verifyOTPRequest)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(204)
}

// @Summary Send OTP to Mail for reset password
// @Description Send OTP to user email when resetting password
// @Tags Auths
// @Accept json
// @Param request body model.SendOTPRequest true "Send OTP payload"
// @Produce json
// @Router /auth/reset-password/send-otp [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) SendOTPToEmailForResetPassword(ctx *gin.Context) {
	var sendOTPRequest model.SendOTPRequest

	if err := validation.BindJsonAndValidate(ctx, &sendOTPRequest); err != nil {
		return
	}

	errCode := handler.authService.SendOTPToEmailForResetPassword(ctx, sendOTPRequest)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(204)
}

// @Summary Verify OTP for reset password
// @Description Verify OTP with email and otp when resetting password
// @Tags Auths
// @Accept json
// @Param request body model.VerifyOTPRequest true "Verify OTP payload"
// @Produce json
// @Router /auth/reset-password/verify-otp [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) VerifyOTPForResetPassword(ctx *gin.Context) {
	var verifyOTPRequest model.VerifyOTPRequest

	if err := validation.BindJsonAndValidate(ctx, &verifyOTPRequest); err != nil {
		return
	}

	errCode := handler.authService.VerifyOTPForResetPassword(ctx, verifyOTPRequest)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(204)
}

// @Summary Set Password
// @Description Set a new password after OTP verification
// @Tags Auths
// @Accept json
// @Param request body model.SetPasswordRequest true "Set Password payload"
// @Produce json
// @Router /auth/forgot-password [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) SetPassword(ctx *gin.Context) {
	var setPasswordRequest model.SetPasswordRequest

	if err := validation.BindJsonAndValidate(ctx, &setPasswordRequest); err != nil {
		return
	}

	errCode := handler.authService.SetPassword(ctx, setPasswordRequest)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(204)
}

func (handler *AuthHandler) Test(ctx *gin.Context) {
	message := "hello world"
	ctx.JSON(200, httpcommon.NewSuccessResponse(&message))
}

// FirebaseLogin handles user login via Firebase (Google OAuth)
// @Summary Login with Google
// @Description Authenticate a user using their Google token via Firebase
// @Tags Auths
// @Accept json
// @Produce json
// @Param request body model.GoogleLoginRequest true "Google Login Request"
// @Success 200 {object} httpcommon.HttpResponse[entity.User]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /auth/google-login [post]
func (handler *AuthHandler) FirebaseLogin(ctx *gin.Context) {
	var googleLoginReq model.GoogleLoginRequest

	if err := validation.BindJsonAndValidate(ctx, &googleLoginReq); err != nil {
		return
	}

	user, errCode := handler.authService.GoogleLogin(ctx, googleLoginReq)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(200, httpcommon.NewSuccessResponse(user))
}
