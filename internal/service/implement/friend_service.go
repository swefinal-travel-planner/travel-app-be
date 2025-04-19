package serviceimplement

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
)

type FriendService struct {
	userRepository   repository.UserRepository
	friendRepository repository.FriendRepository
}

func NewFriendService(
	friendRepository repository.FriendRepository,
	userRepository repository.UserRepository,
) service.FriendService {
	return &FriendService{
		friendRepository: friendRepository,
		userRepository:   userRepository,
	}
}

func (service *FriendService) GetAllFriends(ctx *gin.Context, userId int64) ([]model.FriendResponse, string) {
	friends, err := service.friendRepository.GetByUserIdQuery(ctx, userId)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, ""
		}
		log.Error("FriendService.GetAllFriends error when get user:", err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	var friendResponses []model.FriendResponse
	for _, friend := range friends {
		friendResponses = append(friendResponses, model.FriendResponse{
			Id:       friend.Id,
			Username: friend.Name,
			ImageURL: friend.PhotoURL,
		})
	}
	return friendResponses, ""
}

func (service *FriendService) RemoveFriend(ctx *gin.Context, userId int64, friendId int64) string {
	// Check if the user is a friend
	isFriend := service.friendRepository.ExistsByUserId1AndUserId2Query(ctx, userId, friendId)
	if !isFriend {
		return error_utils.ErrorCode.REMOVE_FRIEND_NOT_FOUND
	}

	// Only delete the friend if the user is a friend
	err := service.friendRepository.DeleteByUserId1AndUserId2Command(ctx, userId, friendId)
	if err != nil {
		log.Error("FriendService.RemoveFriend error when remove friend:", err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	return ""
}
