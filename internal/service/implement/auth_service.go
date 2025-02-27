package serviceimplement

import (
	"database/sql"
	"errors"

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
		PhoneNumber: registerRequest.PhoneNumber,
		Password:    string(hashPW),
	}

	err = service.userRepository.CreateCommand(ctx, user)
	if err != nil {
		return err
	}
	return nil
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

	jwtSecret, err := env.GetEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}
	accessToken, err := jwt.GenerateToken(constants.ACCESS_TOKEN_DURATION, jwtSecret, map[string]interface{}{
		"id": existsUser.Id,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateToken(constants.REFRESH_TOKEN_DURATION, jwtSecret, map[string]interface{}{
		"id": existsUser.Id,
	})
	if err != nil {
		return nil, err
	}

	// Check if a refresh token already exists
	existingRefreshToken, err := service.authenticationRepository.GetOneByUserIdQuery(ctx, existsUser.Id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existingRefreshToken == nil {
		// Create a new refresh token
		err = service.authenticationRepository.CreateCommand(ctx, entity.Authentication{
			UserId:       existsUser.Id,
			RefreshToken: refreshToken,
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Update the existing refresh token
		err = service.authenticationRepository.UpdateCommand(ctx, entity.Authentication{
			UserId:       existsUser.Id,
			RefreshToken: refreshToken,
		})
		if err != nil {
			return nil, err
		}
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

func (service *AuthService) SendOTPToEmail(ctx *gin.Context, sendOTPRequest model.SendOTPRequest) error {
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

func (service *AuthService) VerifyOTP(ctx *gin.Context, verifyOTPRequest model.VerifyOTPRequest) error {
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
		return errors.New("invalid OTP")
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
