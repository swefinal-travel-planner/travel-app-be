package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/constants"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/env"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/jwt"
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

	err := handler.authService.Register(ctx, registerRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
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

	user, err := handler.authService.Login(ctx, loginRequest)
	if err != nil || user == nil {
		if user == nil && err == nil {
			err = errors.New(httpcommon.ErrorMessage.SqlxNoRow)
		}
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
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

	jwtSecret, err := env.GetEnv("JWT_SECRET")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: httpcommon.ErrorMessage.BadCredential,
				Code:    httpcommon.ErrorResponseCode.Unauthorized,
			},
		))
		return
	}

	refreshClaims, errRf := jwt.VerifyToken(refreshTokenRequest.RefreshToken, jwtSecret)
	if errRf != nil {
		// If the refresh token is invalid or expired, abort with unauthorized
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: httpcommon.ErrorMessage.BadCredential,
				Code:    httpcommon.ErrorResponseCode.Unauthorized,
			},
		))
		return
	}

	// Extract user Id from refresh token claims
	payload, ok := refreshClaims.Payload.(map[string]interface{})
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: httpcommon.ErrorMessage.BadCredential,
				Code:    httpcommon.ErrorResponseCode.Unauthorized,
			},
		))
		return
	}
	userId := int64(payload["id"].(float64))

	// Check if the refresh token exists and is still valid in the database
	refreshTokenEntity, err := handler.authService.ValidateRefreshToken(ctx, userId)
	if err != nil || refreshTokenEntity == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: httpcommon.ErrorMessage.BadCredential,
				Code:    httpcommon.ErrorResponseCode.Unauthorized,
			},
		))
		return
	}

	// Generate a new access token
	newAccessToken, err := jwt.GenerateToken(constants.ACCESS_TOKEN_DURATION, jwtSecret, map[string]interface{}{
		"id": userId,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			},
		))
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
// @Router /auth/send-otp/register [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) SendOTPToEmailForRegister(ctx *gin.Context) {
	var sendOTPRequest model.SendOTPRequest

	if err := validation.BindJsonAndValidate(ctx, &sendOTPRequest); err != nil {
		return
	}

	err := handler.authService.SendOTPToEmailForRegister(ctx, sendOTPRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
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
// @Router /auth/verify-otp/register [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) VerifyOTPForRegister(ctx *gin.Context) {
	var verifyOTPRequest model.VerifyOTPRequest

	if err := validation.BindJsonAndValidate(ctx, &verifyOTPRequest); err != nil {
		return
	}

	err := handler.authService.VerifyOTPForRegister(ctx, verifyOTPRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
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
// @Router /auth/send-otp/reset-password [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) SendOTPToEmailForResetPassword(ctx *gin.Context) {
	var sendOTPRequest model.SendOTPRequest

	if err := validation.BindJsonAndValidate(ctx, &sendOTPRequest); err != nil {
		return
	}

	err := handler.authService.SendOTPToEmailForResetPassword(ctx, sendOTPRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
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
// @Router /auth/verify-otp/reset-password [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) VerifyOTPForResetPassword(ctx *gin.Context) {
	var verifyOTPRequest model.VerifyOTPRequest

	if err := validation.BindJsonAndValidate(ctx, &verifyOTPRequest); err != nil {
		return
	}

	err := handler.authService.VerifyOTPForResetPassword(ctx, verifyOTPRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
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

	err := handler.authService.SetPassword(ctx, setPasswordRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
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

	if err := ctx.ShouldBindJSON(&googleLoginReq); err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
		return
	}

	user, err := handler.authService.GoogleLogin(ctx, googleLoginReq)
	if err != nil || user == nil {
		if user == nil && err == nil {
			err = errors.New(httpcommon.ErrorMessage.SqlxNoRow)
		}
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
		return
	}

	ctx.JSON(200, httpcommon.NewSuccessResponse(user))
}
