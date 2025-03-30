package serviceimplement

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
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

func (service *FriendService) GetAllFriends(ctx *gin.Context, userId int64) ([]model.FriendResponse, error) {
	friends, err := service.friendRepository.GetByUserIdCommand(ctx, userId)
	if err != nil {
		return nil, err
	}
	var friendResponses []model.FriendResponse
	for _, friendRecord := range friends {
		user, err := service.userRepository.GetOneByIDQuery(ctx, userId)
		if err != nil {
			return nil, err
		}

		// get user from friend record
		var friend *entity.User
		if user.Id == friendRecord.UserID1 {
			friend, err = service.userRepository.GetOneByIDQuery(ctx, friendRecord.UserID2)
			if err != nil {
				return nil, err
			}
		} else {
			friend, err = service.userRepository.GetOneByIDQuery(ctx, friendRecord.UserID1)
			if err != nil {
				return nil, err
			}
		}

		friendResponses = append(friendResponses, model.FriendResponse{
			Username: friend.Name,
			ImageURL: friend.PhotoURL,
		})
	}
	return friendResponses, nil
}
