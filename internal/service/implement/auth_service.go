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
)

type AuthService struct {
	userRepository           repository.UserRepository
	authenticationRepository repository.AuthenticationRepository
	passwordEncoder          bean.PasswordEncoder
}

func NewAuthService(userRepository repository.UserRepository,
	authenticationRepository repository.AuthenticationRepository,
	passwordEncoder bean.PasswordEncoder,
) service.AuthService {
	return &AuthService{
		userRepository:           userRepository,
		authenticationRepository: authenticationRepository,
		passwordEncoder:          passwordEncoder,
	}
}

func (service *AuthService) Register(ctx *gin.Context, registerRequest model.RegisterRequest) error {
	existsCustomer, err := service.userRepository.GetOneByEmailQuery(ctx, registerRequest.Email)
	if err != nil && err.Error() != httpcommon.ErrorMessage.SqlxNoRow {
		return err
	}
	if existsCustomer != nil {
		return errors.New("Email have already registered")
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
			return nil, errors.New("Email not found")
		}
		return nil, err
	}
	checkPw := service.passwordEncoder.Compare(existsUser.Password, loginRequest.Password)
	if checkPw == false {
		return nil, errors.New("invalid password")
	}

	jwtSecret, err := env.GetEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}
	accessToken, err := jwt.GenerateToken(constants.ACCESS_TOKEN_DURATION, jwtSecret, map[string]interface{}{
		"id": existsUser.Id,
	})

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
