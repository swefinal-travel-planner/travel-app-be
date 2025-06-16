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
	friend, err := service.userRepository.GetOneByEmailQuery(ctx, userEmail, nil)
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

	if service.friendRepository.ExistsByUserId1AndUserId2Query(ctx, userId, friend.Id, nil) {
		status = model.FriendStatus.Friend
	} else if invitationFriend, err := service.invitationFriendRepository.GetBySenderAndReceiverIdQuery(ctx, userId, friend.Id, nil); err == nil {
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

func (service *UserService) UpdateNotificationToken(ctx *gin.Context, userId int64, notificationTokenRequest model.UpdateNotificationTokenRequest) string {
	err := service.userRepository.UpdateNotificationTokenCommand(ctx, userId, notificationTokenRequest.NotificationToken, nil)
	if err != nil {
		log.Error("UserService.UpdateNotificationToken Error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	return ""
}

func (service *UserService) UpdateUser(ctx *gin.Context, userId int64, request model.UpdateUserRequest) string {
	// Get existing user
	user, err := service.userRepository.GetOneByIDQuery(ctx, userId, nil)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return error_utils.ErrorCode.FORBIDDEN
		}
		log.Error("UserService.UpdateUser Error getting user: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	// If email is being updated, check if it's already in use
	if request.Email != "" && request.Email != user.Email {
		existingUser, err := service.userRepository.GetOneByEmailQuery(ctx, request.Email, nil)
		if err == nil && existingUser.Id != userId {
			return error_utils.ErrorCode.REGISTER_EMAIL_EXISTED
		}
		if err != nil && err.Error() != error_utils.SystemErrorMessage.SqlxNoRow {
			log.Error("UserService.UpdateUser Error checking email: " + err.Error())
			return error_utils.ErrorCode.DB_DOWN
		}
		user.Email = request.Email
	}

	// Update other fields if provided
	if request.Name != "" {
		user.Name = request.Name
	}
	if request.PhoneNumber != "" {
		user.PhoneNumber = request.PhoneNumber
	}
	if request.PhotoURL != nil {
		user.PhotoURL = request.PhotoURL
	}
	if request.NotificationToken != nil {
		user.NotificationToken = request.NotificationToken
	}

	// Save the updated user
	err = service.userRepository.UpdateCommand(ctx, user, nil)
	if err != nil {
		log.Error("UserService.UpdateUser Error updating user: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	return ""
}
