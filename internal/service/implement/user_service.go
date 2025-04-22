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
	userRepository             repository.UserRepository
	friendRepository           repository.FriendRepository
	invitationFriendRepository repository.InvitationFriendRepository
	invitationFriendService    service.InvitationFriendService
}

func NewUserService(
	userRepository repository.UserRepository,
	friendRepository repository.FriendRepository,
	invitationFriendRepository repository.InvitationFriendRepository,
	invitationFriendService service.InvitationFriendService,
) service.UserService {
	return &UserService{
		userRepository:             userRepository,
		friendRepository:           friendRepository,
		invitationFriendRepository: invitationFriendRepository,
		invitationFriendService:    invitationFriendService,
	}
}

func (service *UserService) SearchUser(ctx *gin.Context, userId int64, userEmail string) (*model.FriendResponse, string) {
	friend, err := service.userRepository.GetOneByEmailQuery(ctx, userEmail)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, ""
		}
		log.Error("UserService.SearchUser Error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// check to see if the user is a friend or not
	// if not a friend, check to see if the user is the one who sent the invitation or not
	// if not sent the invitation, check to see if the user is the one who received the invitation or not
	// if not invited, check to see if the user is rejcted or not (in cooldown table, user must be a sender)
	// if not match any cases above, user is stranger
	var status string
	var timeRemaining *int64 = nil

	if service.friendRepository.ExistsByUserId1AndUserId2Query(ctx, userId, friend.Id) {
		status = model.FriendStatus.Friend
	} else if invitationFriend, err := service.invitationFriendRepository.GetBySenderAndReceiverIdQuery(ctx, userId, friend.Id); err == nil {
		if invitationFriend.SenderID == userId { // user is sender
			status = model.FriendStatus.Sent
		} else if invitationFriend.ReceiverID == userId { // user is receiver
			status = model.FriendStatus.Received
		}
	} else if _, cooldownRemaining := service.invitationFriendService.IsInCoolDownAndGetRemainingTime(ctx, userId, friend.Id); cooldownRemaining > 0 {
		status = model.FriendStatus.Restricted
		timeRemaining = &cooldownRemaining
	} else {
		status = model.FriendStatus.Stranger
	}

	friendResponse := &model.FriendResponse{
		Id:            friend.Id,
		Email:         friend.Email,
		Username:      friend.Name,
		ImageURL:      friend.PhotoURL,
		Status:        status,
		TimeRemaining: timeRemaining,
	}
	return friendResponse, ""
}
