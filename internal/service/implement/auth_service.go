package serviceimplement

import (
	"fmt"

	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal/bean"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/constants"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/env"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/jwt"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/mail"
	redisHelper "github.com/swefinal-travel-planner/travel-app-be/internal/utils/redis_helper"
)

type AuthService struct {
	userRepository           repository.UserRepository
	authenticationRepository repository.AuthenticationRepository
	passwordEncoder          bean.PasswordEncoder
	redisClient              bean.RedisClient
	mailClient               bean.MailClient
}

func NewAuthService(userRepository repository.UserRepository,
	authenticationRepository repository.AuthenticationRepository,
	passwordEncoder bean.PasswordEncoder,
	redisClient bean.RedisClient,
	mailClient bean.MailClient,
) service.AuthService {
	return &AuthService{
		userRepository:           userRepository,
		authenticationRepository: authenticationRepository,
		passwordEncoder:          passwordEncoder,
		redisClient:              redisClient,
		mailClient:               mailClient,
	}
}

func (service *AuthService) Register(ctx *gin.Context, registerRequest model.RegisterRequest) string {
	// OTP validation
	email := registerRequest.Email
	baseKey := constants.VERIFY_EMAIL_KEY
	key := fmt.Sprintf("%s:%s", baseKey, email)

	// Register if pass OTP validation
	existsCustomer, err := service.userRepository.GetOneByEmailQuery(ctx, registerRequest.Email, nil)
	if err != nil {
		log.Error("AuthService.Register DB is down: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if existsCustomer != nil {
		return error_utils.ErrorCode.REGISTER_EMAIL_EXISTED
	}

	val, err := service.redisClient.Get(ctx, key)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.RedisNil {
			return error_utils.ErrorCode.REGISTER_OTP_NOT_FOUND
		}
		log.Error("AuthService.Register Redis is down when get: " + err.Error())
		return error_utils.ErrorCode.REDIS_DOWN
	}

	if val == registerRequest.OTP {
		err := service.redisClient.Delete(ctx, key)
		if err != nil {
			log.Error("AuthService.Register Redis is down when delete: " + err.Error())
			return error_utils.ErrorCode.REDIS_DOWN
		}
	} else {
		return error_utils.ErrorCode.REGISTER_OTP_INVALID
	}

	hashPW, err := service.passwordEncoder.Encrypt(registerRequest.Password)
	if err != nil {
		log.Error("AuthService.Register Error when encrypt password: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	user := &entity.User{
		Email:       registerRequest.Email,
		Name:        registerRequest.Name,
		PhoneNumber: "",
		Password:    string(hashPW),
	}

	err = service.userRepository.CreateCommand(ctx, user, nil)
	if err != nil {
		log.Error("AuthService.Register Error when create user: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	return ""
}

func (service *AuthService) generateAndStoreTokens(ctx *gin.Context, userId int64) (string, string, string) {
	jwtSecret, err := env.GetEnv("JWT_SECRET")
	if err != nil {
		log.Error("AuthService.generateAndStoreTokens Error when get JWT secret: " + err.Error())
		return "", "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	accessToken, err := jwt.GenerateToken(constants.ACCESS_TOKEN_DURATION, jwtSecret, map[string]interface{}{
		"id": userId,
	})
	if err != nil {
		log.Error("AuthService.generateAndStoreTokens Error when generate access token: " + err.Error())
		return "", "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	refreshToken, err := jwt.GenerateToken(constants.REFRESH_TOKEN_DURATION, jwtSecret, map[string]interface{}{
		"id": userId,
	})
	if err != nil {
		log.Error("AuthService.generateAndStoreTokens Error when generate refresh token: " + err.Error())
		return "", "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// Check if a refresh token already exists
	existingRefreshToken, err := service.authenticationRepository.GetOneByUserIdQuery(ctx, userId, nil)
	if err != nil {
		log.Error("AuthService.generateAndStoreTokens Error when get existing refresh token: " + err.Error())
		return "", "", error_utils.ErrorCode.DB_DOWN
	}

	authData := entity.Authentication{
		UserId:       userId,
		RefreshToken: refreshToken,
	}

	if existingRefreshToken == nil {
		// Create a new refresh token
		err = service.authenticationRepository.CreateCommand(ctx, authData, nil)
	} else {
		// Update the existing refresh token
		err = service.authenticationRepository.UpdateCommand(ctx, authData, nil)
	}

	if err != nil {
		log.Error("AuthService.generateAndStoreTokens Error when create/update user's authentication: " + err.Error())
		return "", "", error_utils.ErrorCode.DB_DOWN
	}

	return accessToken, refreshToken, ""
}

func (service *AuthService) Login(ctx *gin.Context, loginRequest model.LoginRequest) (*model.LoginResponse, string) {
	existsUser, err := service.userRepository.GetOneByEmailQuery(ctx, loginRequest.Email, nil)
	if err != nil {
		log.Error("AuthService.Login Error when get user: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	if existsUser == nil {
		return nil, error_utils.ErrorCode.LOGIN_EMAIL_NOT_FOUND
	}
	checkPw := service.passwordEncoder.Compare(existsUser.Password, loginRequest.Password)
	if !checkPw {
		return nil, error_utils.ErrorCode.LOGIN_INVALID_PASSWORD
	}

	// Generate and store tokens
	accessToken, refreshToken, errCode := service.generateAndStoreTokens(ctx, existsUser.Id)
	if errCode != "" {
		return nil, errCode
	}

	return &model.LoginResponse{
		Name:         existsUser.Name,
		Email:        existsUser.Email,
		UserId:       existsUser.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, ""
}

func (service *AuthService) GoogleLogin(ctx *gin.Context, loginRequest model.GoogleLoginRequest) (*model.LoginResponse, string) {
	existsUser, err := service.userRepository.GetOneByEmailQuery(ctx, loginRequest.Email, nil)
	if err != nil {
		log.Error("AuthService.GoogleLogin Error when get user: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	if existsUser == nil { // email not founded => create new user with googleID
		newUser := &entity.User{
			Email:       loginRequest.Email,
			Name:        loginRequest.DisplayName,
			PhoneNumber: loginRequest.PhoneNumber,
			Password:    "",
			IDToken:     loginRequest.IDToken,
		}
		err = service.userRepository.CreateCommand(ctx, newUser, nil)
		if err != nil {
			log.Error("AuthService.GoogleLogin Error when create user: " + err.Error())
			return nil, error_utils.ErrorCode.DB_DOWN
		}
		existsUser, err = service.userRepository.GetOneByEmailQuery(ctx, loginRequest.Email, nil)
		if err != nil {
			log.Error("AuthService.GoogleLogin Error when get existing user: " + err.Error())
			return nil, error_utils.ErrorCode.DB_DOWN
		}
	}

	// User exists
	// 1. User register internally => Deny
	// 2. User register with google => Process to login
	if existsUser.IDToken == nil {
		return nil, error_utils.ErrorCode.REGISTER_EMAIL_EXISTED
	}

	// Generate and store tokens
	accessToken, refreshToken, errCode := service.generateAndStoreTokens(ctx, existsUser.Id)
	if errCode != "" {
		return nil, errCode
	}
	return &model.LoginResponse{
		Name:         existsUser.Name,
		Email:        existsUser.Email,
		UserId:       existsUser.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, ""
}

func (service *AuthService) SendOTPToEmailForRegister(ctx *gin.Context, sendOTPRequest model.SendOTPRequest) string {
	// check if email exists
	user, err := service.userRepository.GetOneByEmailQuery(ctx, sendOTPRequest.Email, nil)
	if err != nil {
		log.Error("AuthService.SendOTPToEmailForRegister Error when get user by email: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if user != nil {
		return error_utils.ErrorCode.REGISTER_EMAIL_EXISTED
	}

	// generate otp
	otp := mail.GenerateOTP(6)

	// store otp in redis
	email := sendOTPRequest.Email
	baseKey := constants.VERIFY_EMAIL_KEY
	key := fmt.Sprintf("%s:%s", baseKey, email)

	err = service.redisClient.Set(ctx, key, otp)
	if err != nil {
		log.Error("AuthService.SendOTPToEmailForRegister Error when set redis key: " + err.Error())
		return error_utils.ErrorCode.REDIS_DOWN
	}

	// send otp to user email
	emailBody := service.mailClient.GenerateOTPBody(sendOTPRequest.Email, otp, constants.VERIFY_EMAIL, constants.VERIFY_EMAIL_EXP_TIME)
	err = service.mailClient.SendEmail(ctx, sendOTPRequest.Email, "OTP verify email", emailBody)
	if err != nil {
		log.Error("AuthService.SendOTPToEmailForRegister Error when send email: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	log.Info("AuthService.SendOTPToEmailForRegister OTP is: " + otp)

	return ""
}

func (service *AuthService) VerifyOTPForRegister(ctx *gin.Context, verifyOTPRequest model.VerifyOTPRequest) string {
	email := verifyOTPRequest.Email
	baseKey := constants.VERIFY_EMAIL_KEY
	key := fmt.Sprintf("%s:%s", baseKey, email)

	val, err := service.redisClient.Get(ctx, key)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.RedisNil {
			return error_utils.ErrorCode.REGISTER_OTP_NOT_FOUND
		}
		log.Error("AuthService.VerifyOTPForRegister Error when get redis key: " + err.Error())
		return error_utils.ErrorCode.REDIS_DOWN
	}

	if val != verifyOTPRequest.OTP {
		return error_utils.ErrorCode.REGISTER_OTP_INVALID
	}

	return ""
}

func (service *AuthService) SendOTPToEmailForResetPassword(ctx *gin.Context, sendOTPRequest model.SendOTPRequest) string {
	// generate otp
	otp := mail.GenerateOTP(6)

	// store otp in redis
	customerId, err := service.userRepository.GetIdByEmailQuery(ctx, sendOTPRequest.Email, nil)
	if err != nil {
		log.Error("AuthService.SendOTPToEmailForResetPassword Error when get ID by email: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if customerId == 0 {
		return error_utils.ErrorCode.RESET_PASSWORD_EMAIL_NOT_FOUND
	}

	baseKey := constants.RESET_PASSWORD_KEY
	key := redisHelper.Concat(baseKey, customerId)

	err = service.redisClient.Set(ctx, key, otp)
	if err != nil {
		log.Error("AuthService.SendOTPToEmailForResetPassword Error when set redis key: " + err.Error())
		return error_utils.ErrorCode.REDIS_DOWN
	}

	// send otp to user email
	emailBody := service.mailClient.GenerateOTPBody(sendOTPRequest.Email, otp, constants.FORGOT_PASSWORD, constants.RESET_PASSWORD_EXP_TIME)
	err = service.mailClient.SendEmail(ctx, sendOTPRequest.Email, "OTP reset password", emailBody)
	if err != nil {
		log.Error("AuthService.SendOTPToEmailForResetPassword Error when send email: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	log.Info("AuthService.SendOTPToEmailForResetPassword OTP is: " + otp)

	return ""
}

func (service *AuthService) VerifyOTPForResetPassword(ctx *gin.Context, verifyOTPRequest model.VerifyOTPRequest) string {
	customerId, err := service.userRepository.GetIdByEmailQuery(ctx, verifyOTPRequest.Email, nil)
	if err != nil {
		log.Error("AuthService.VerifyOTPForResetPassword Error when get ID by email: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if customerId == 0 {
		return error_utils.ErrorCode.RESET_PASSWORD_EMAIL_NOT_FOUND
	}

	baseKey := constants.RESET_PASSWORD_KEY
	key := redisHelper.Concat(baseKey, customerId)

	val, err := service.redisClient.Get(ctx, key)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.RedisNil {
			return error_utils.ErrorCode.RESET_PASSWORD_EMAIL_NOT_FOUND
		}
		log.Error("AuthService.VerifyOTPForResetPassword Error when get redis key: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	if val != verifyOTPRequest.OTP {
		return error_utils.ErrorCode.RESET_PASSWORD_OTP_INVALID
	}

	return ""
}

func (service *AuthService) SetPassword(ctx *gin.Context, setPasswordRequest model.SetPasswordRequest) string {
	customerId, err := service.userRepository.GetIdByEmailQuery(ctx, setPasswordRequest.Email, nil)
	if err != nil {
		log.Error("AuthService.SetPassword Error when get ID by email: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if customerId == 0 {
		return error_utils.ErrorCode.RESET_PASSWORD_EMAIL_NOT_FOUND
	}

	baseKey := constants.RESET_PASSWORD_KEY
	key := redisHelper.Concat(baseKey, customerId)

	val, err := service.redisClient.Get(ctx, key)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.RedisNil {
			return error_utils.ErrorCode.SET_PASSWORD_OTP_NOT_FOUND
		}
		log.Error("AuthService.SetPassword Error when get redis key: " + err.Error())
		return error_utils.ErrorCode.REDIS_DOWN
	}

	if val == setPasswordRequest.OTP {
		err := service.redisClient.Delete(ctx, key)
		if err != nil {
			log.Error("AuthService.SetPassword Error when delete redis key: " + err.Error())
			return error_utils.ErrorCode.REDIS_DOWN
		}

		hashedPW, err := service.passwordEncoder.Encrypt(setPasswordRequest.Password)
		if err != nil {
			log.Error("AuthService.SetPassword Error when encrypt password: " + err.Error())
			return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
		}

		err = service.userRepository.UpdatePasswordByIdQuery(ctx, customerId, hashedPW, nil)
		if err != nil {
			log.Error("AuthService.SetPassword Error when update password: " + err.Error())
			return error_utils.ErrorCode.DB_DOWN
		}
	} else {
		return error_utils.ErrorCode.SET_PASSWORD_OTP_INVALID
	}

	return ""
}

func (service *AuthService) RefreshToken(ctx *gin.Context, refreshTokenRequest model.RefreshTokenRequest) (string, string) {
	jwtSecret, err := env.GetEnv("JWT_SECRET")
	if err != nil {
		log.Error("AuthService.RefreshToken Error when get JWT secret: " + err.Error())
		return "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	refreshClaims, errRf := jwt.VerifyToken(refreshTokenRequest.RefreshToken, jwtSecret)
	if errRf != nil {
		log.Error("AuthService.RefreshToken Error when verify JWT secret: " + errRf.Error())
		return "", error_utils.ErrorCode.REFRESH_TOKEN_INVALID
	}

	// Extract user Id from refresh token claims
	payload, ok := refreshClaims.Payload.(map[string]interface{})
	if !ok {
		log.Error("AuthService.RefreshToken Error when extracting claims from request token")
		return "", error_utils.ErrorCode.REFRESH_TOKEN_INVALID
	}
	userId := int64(payload["id"].(float64))

	// Check if the refresh token exists in the database
	existingRefreshToken, err := service.authenticationRepository.GetOneByUserIdQuery(ctx, userId, nil)
	if err != nil {
		log.Error("AuthService.RefreshToken Error when get existing refresh token: " + err.Error())
		return "", error_utils.ErrorCode.DB_DOWN
	}
	if existingRefreshToken == nil {
		return "", error_utils.ErrorCode.REFRESH_TOKEN_INVALID
	}
	if existingRefreshToken.RefreshToken != refreshTokenRequest.RefreshToken {
		return "", error_utils.ErrorCode.REFRESH_TOKEN_INVALID
	}

	// Generate a new access token
	newAccessToken, err := jwt.GenerateToken(constants.ACCESS_TOKEN_DURATION, jwtSecret, map[string]interface{}{
		"id": userId,
	})
	if err != nil {
		log.Error("AuthService.RefreshToken Error when generate access token: " + err.Error())
		return "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	return newAccessToken, ""
}
