package serviceimplement

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(
	userRepository repository.UserRepository,
) service.UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (service *UserService) SearchUser(ctx *gin.Context, userEmail string) (*model.FriendResponse, string) {
	friend, err := service.userRepository.GetOneByEmailQuery(ctx, userEmail)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, ""
		}
		log.Error("UserService.SearchUser Error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	friendResponse := &model.FriendResponse{
		Id:       friend.Id,
		Username: friend.Name,
		ImageURL: friend.PhotoURL,
	}
	return friendResponse, ""
}
