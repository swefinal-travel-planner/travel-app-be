package serviceimplement

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
)

type AuthService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) service.AuthService {
	return &AuthService{userRepository: userRepository}
}

func (service *AuthService) Register(ctx *gin.Context, registerRequest model.RegisterRequest) error {
	user := &entity.User{
		Name:     registerRequest.Name,
		Password: registerRequest.Password,
	}
	err := service.userRepository.CreateCommand(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
