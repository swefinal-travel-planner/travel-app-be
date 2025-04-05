package serviceimplement

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/bean"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/constants"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/env"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/jwt"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/mail"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/redis"
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

func (service *AuthService) Register(ctx *gin.Context, registerRequest model.RegisterRequest) error {
	existsCustomer, err := service.userRepository.GetOneByEmailQuery(ctx, registerRequest.Email)
	if err != nil && err.Error() != httpcommon.ErrorMessage.SqlxNoRow {
		return err
	}
	if existsCustomer != nil {
		return errors.New("email have already registered")
	}

	hashPW, err := service.passwordEncoder.Encrypt(registerRequest.Password)
	if err != nil {
		return err
	}

	user := &entity.User{
		Email:       registerRequest.Email,
		Name:        registerRequest.Name,
		PhoneNumber: "",
		Password:    string(hashPW),
	}

	err = service.userRepository.CreateCommand(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (service *AuthService) generateAndStoreTokens(ctx *gin.Context, userId int64) (string, string, error) {
	jwtSecret, err := env.GetEnv("JWT_SECRET")
	if err != nil {
		return "", "", err
	}

	accessToken, err := jwt.GenerateToken(constants.ACCESS_TOKEN_DURATION, jwtSecret, map[string]interface{}{
		"id": userId,
	})
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.GenerateToken(constants.REFRESH_TOKEN_DURATION, jwtSecret, map[string]interface{}{
		"id": userId,
	})
	if err != nil {
		return "", "", err
	}

	// Check if a refresh token already exists
	existingRefreshToken, err := service.authenticationRepository.GetOneByUserIdQuery(ctx, userId)
	if err != nil && err != sql.ErrNoRows {
		return "", "", err
	}

	authData := entity.Authentication{
		UserId:       userId,
		RefreshToken: refreshToken,
	}

	if existingRefreshToken == nil {
		// Create a new refresh token
		err = service.authenticationRepository.CreateCommand(ctx, authData)
	} else {
		// Update the existing refresh token
		err = service.authenticationRepository.UpdateCommand(ctx, authData)
	}

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (service *AuthService) Login(ctx *gin.Context, loginRequest model.LoginRequest) (*model.LoginResponse, error) {
	existsUser, err := service.userRepository.GetOneByEmailQuery(ctx, loginRequest.Email)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, errors.New("email not found")
		}
		return nil, err
	}
	checkPw := service.passwordEncoder.Compare(existsUser.Password, loginRequest.Password)
	if !checkPw {
		return nil, errors.New("invalid password")
	}

	// Generate and store tokens
	accessToken, refreshToken, err := service.generateAndStoreTokens(ctx, existsUser.Id)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Name:         existsUser.Name,
		Email:        existsUser.Email,
		UserId:       existsUser.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (service *AuthService) GoogleLogin(ctx *gin.Context, loginRequest model.GoogleLoginRequest) (*model.LoginResponse, error) {
	existsUser, err := service.userRepository.GetOneByEmailQuery(ctx, loginRequest.Email)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow { // email not founded => create new user with googleID
			newUser := &entity.User{
				Email:       loginRequest.Email,
				Name:        loginRequest.DisplayName,
				PhoneNumber: loginRequest.PhoneNumber,
				Password:    "",
				IDToken:     loginRequest.IDToken,
			}
			err = service.userRepository.CreateCommand(ctx, newUser)
			if err != nil {
				return nil, err
			}
			existsUser, err = service.userRepository.GetOneByEmailQuery(ctx, loginRequest.Email)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	// User exists
	// 1. User register internally => Deny
	// 2. User register with google => Process to login
	if existsUser.IDToken == nil {
		return nil, errors.New("this email is already registered")
	}

	// Generate and store tokens
	accessToken, refreshToken, err := service.generateAndStoreTokens(ctx, existsUser.Id)
	if err != nil {
		return nil, err
	}
	return &model.LoginResponse{
		Name:         existsUser.Name,
		Email:        existsUser.Email,
		UserId:       existsUser.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (service *AuthService) ValidateRefreshToken(ctx *gin.Context, userId int64) (*entity.Authentication, error) {
	refreshToken, err := service.authenticationRepository.GetOneByUserIdQuery(ctx, userId)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}

func (service *AuthService) SendOTPToEmailForRegister(ctx *gin.Context, sendOTPRequest model.SendOTPRequest) error {
	// check if email exists
	_, err := service.userRepository.GetOneByEmailQuery(ctx, sendOTPRequest.Email)
	if err == nil {
		return errors.New("email already exists")
	}

	// generate otp
	otp := mail.GenerateOTP(6)

	// store otp in redis
	email := sendOTPRequest.Email
	baseKey := constants.VERIFY_EMAIL_KEY
	key := fmt.Sprintf("%s:%s", baseKey, email)

	err = service.redisClient.Set(ctx, key, otp)
	if err != nil {
		return err
	}

	// send otp to user email
	emailBody := service.mailClient.GenerateOTPBody(sendOTPRequest.Email, otp, constants.VERIFY_EMAIL, constants.VERIFY_EMAIL_EXP_TIME)
	err = service.mailClient.SendEmail(ctx, sendOTPRequest.Email, "OTP verify email", emailBody)
	if err != nil {
		return err
	}

	return nil
}

func (service *AuthService) VerifyOTPForRegister(ctx *gin.Context, verifyOTPRequest model.VerifyOTPRequest) error {
	email := verifyOTPRequest.Email
	baseKey := constants.VERIFY_EMAIL_KEY
	key := fmt.Sprintf("%s:%s", baseKey, email)

	val, err := service.redisClient.Get(ctx, key)
	if err != nil {
		return err
	}

	if val != verifyOTPRequest.OTP {
		return errors.New("invalid OTP")
	}

	return nil
}

func (service *AuthService) SendOTPToEmailForResetPassword(ctx *gin.Context, sendOTPRequest model.SendOTPRequest) error {
	// generate otp
	otp := mail.GenerateOTP(6)

	// store otp in redis
	customerId, err := service.userRepository.GetIdByEmailQuery(ctx, sendOTPRequest.Email)
	if err != nil {
		return err
	}
	baseKey := constants.RESET_PASSWORD_KEY
	key := redis.Concat(baseKey, customerId)

	err = service.redisClient.Set(ctx, key, otp)
	if err != nil {
		return err
	}

	// send otp to user email
	emailBody := service.mailClient.GenerateOTPBody(sendOTPRequest.Email, otp, constants.FORGOT_PASSWORD, constants.RESET_PASSWORD_EXP_TIME)
	err = service.mailClient.SendEmail(ctx, sendOTPRequest.Email, "OTP reset password", emailBody)
	if err != nil {
		return err
	}

	return nil
}

func (service *AuthService) VerifyOTPForResetPassword(ctx *gin.Context, verifyOTPRequest model.VerifyOTPRequest) error {
	customerId, err := service.userRepository.GetIdByEmailQuery(ctx, verifyOTPRequest.Email)
	if err != nil {
		return err
	}

	baseKey := constants.RESET_PASSWORD_KEY
	key := redis.Concat(baseKey, customerId)

	val, err := service.redisClient.Get(ctx, key)
	if err != nil {
		return err
	}

	if val != verifyOTPRequest.OTP {
		return errors.New("Invalid OTP")
	}

	return nil
}

func (service *AuthService) SetPassword(ctx *gin.Context, setPasswordRequest model.SetPasswordRequest) error {
	customerId, err := service.userRepository.GetIdByEmailQuery(ctx, setPasswordRequest.Email)
	if err != nil {
		return err
	}

	baseKey := constants.RESET_PASSWORD_KEY
	key := redis.Concat(baseKey, customerId)

	val, err := service.redisClient.Get(ctx, key)
	if err != nil {
		return err
	}

	if val == setPasswordRequest.OTP {
		service.redisClient.Delete(ctx, key)

		hashedPW, err := service.passwordEncoder.Encrypt(setPasswordRequest.Password)
		if err != nil {
			return err
		}

		err = service.userRepository.UpdatePasswordByIdQuery(ctx, customerId, hashedPW)
		if err != nil {
			return err
		}
	} else {
		return errors.New("invalid OTP")
	}

	return nil
}
