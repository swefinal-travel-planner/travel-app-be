package serviceimplement

import (
	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/constants"
)

type InvitationFriendService struct {
	invitationFriendRepository   repository.InvitationFriendRepository
	userRepository               repository.UserRepository
	friendRepository             repository.FriendRepository
	invitationCooldownRepository repository.InvitationCooldownRepository
}

func NewInvitationFriendService(
	invitationFriendRepository repository.InvitationFriendRepository,
	userRepository repository.UserRepository,
	friendRepository repository.FriendRepository,
	invitationCooldownRepository repository.InvitationCooldownRepository,
) service.InvitationFriendService {
	return &InvitationFriendService{
		invitationFriendRepository:   invitationFriendRepository,
		userRepository:               userRepository,
		friendRepository:             friendRepository,
		invitationCooldownRepository: invitationCooldownRepository,
	}
}

func (service *InvitationFriendService) AddFriend(ctx *gin.Context, invitation model.InvitationFriendRequest, userId int64) string {
	friend, err := service.userRepository.GetOneByEmailQuery(ctx, invitation.ReceiverEmail)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return error_utils.ErrorCode.ADD_FRIEND_RECEIVER_NOT_FOUND
		}
		log.Error("InvitationFriendService.AddFriend GetOneByEmailQuery error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if userId == friend.Id {
		return error_utils.ErrorCode.ADD_FRIEND_RECEIVER_NOT_FOUND
	}

	// Check if users are in cooldown period
	inCooldown := service.IsInCooldown(ctx, userId, friend.Id)
	if inCooldown {
		return error_utils.ErrorCode.ADD_FRIEND_IN_COOLDOWN
	}

	_, err = service.invitationFriendRepository.GetBySenderAndReceiverIdQuery(ctx, userId, friend.Id)
	if err == nil {
		return error_utils.ErrorCode.ADD_FRIEND_INVITATION_ALREADY_EXISTS
	}

	isFriend := service.friendRepository.ExistsByUserId1AndUserId2Query(ctx, userId, friend.Id)
	if isFriend {
		return error_utils.ErrorCode.ADD_FRIEND_ALREADY_FRIEND
	}

	err = service.invitationFriendRepository.CreateCommand(ctx, &entity.InvitationFriend{
		SenderID:   userId,
		ReceiverID: friend.Id,
	})
	if err != nil {
		log.Error("InvitationFriendService.AddFriend CreateCommand error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	return ""
}

// PROBLEM: N+1
func (service *InvitationFriendService) GetAllReceivedInvitations(ctx *gin.Context, userId int64) ([]model.InvitationFriendReceivedResponse, string) {
	invitations, err := service.invitationFriendRepository.GetByReceiverIdQuery(ctx, userId)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return make([]model.InvitationFriendReceivedResponse, 0), ""
		}
		log.Error("InvitationFriendService.GetAllReceivedInvitations GetByReceiverIdQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	var invitationResponses []model.InvitationFriendReceivedResponse
	for _, invitation := range invitations {
		user, err := service.userRepository.GetOneByIDQuery(ctx, invitation.SenderID)
		if err != nil {
			log.Error("InvitationFriendService.GetAllReceivedInvitations GetOneByIDQuery error: " + err.Error())
			return nil, error_utils.ErrorCode.DB_DOWN
		}
		invitationResponses = append(invitationResponses, model.InvitationFriendReceivedResponse{
			Id:             invitation.ID,
			SenderUsername: user.Name,
			SenderImageURL: user.PhotoURL,
		})
	}
	return invitationResponses, ""
}

// PROBLEM: N+1
func (service *InvitationFriendService) GetAllRequestedInvitations(ctx *gin.Context, userId int64) ([]model.InvitationFriendRequestedResponse, string) {
	invitations, err := service.invitationFriendRepository.GetBySenderIdQuery(ctx, userId)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return make([]model.InvitationFriendRequestedResponse, 0), ""
		}
		log.Error("InvitationFriendService.GetAllRequestedInvitations GetBySenderIdQuery error: ", err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	var invitationResponses []model.InvitationFriendRequestedResponse
	for _, invitation := range invitations {
		user, err := service.userRepository.GetOneByIDQuery(ctx, invitation.SenderID)
		if err != nil {
			log.Error("InvitationFriendService.GetAllRequestedInvitations GetOneByIDQuery error: " + err.Error())
			return nil, error_utils.ErrorCode.DB_DOWN
		}
		invitationResponses = append(invitationResponses, model.InvitationFriendRequestedResponse{
			Id:               invitation.ID,
			ReceiverUsername: user.Name,
			ReceiverImageURL: user.PhotoURL,
		})
	}
	return invitationResponses, ""
}

func (service *InvitationFriendService) validateInvitation(ctx *gin.Context, invitationId int64, userId int64) (*entity.InvitationFriend, string) {
	invitation, err := service.invitationFriendRepository.GetOneByIDQuery(ctx, invitationId)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return nil, error_utils.ErrorCode.FRIEND_INVITATION_NOT_FOUND
		}
		log.Error("InvitationFriendService.validateInvitation GetOneByIDQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	if userId == invitation.SenderID {
		return nil, error_utils.ErrorCode.FRIEND_INVITATION_CANNOT_ACCEPT_AS_SENDER
	}
	if userId != invitation.ReceiverID {
		return nil, error_utils.ErrorCode.FRIEND_INVITATION_NOT_FOUND
	}
	return invitation, ""
}

func (service *InvitationFriendService) AcceptInvitation(ctx *gin.Context, invitationId int64, userId int64) string {
	invitation, errCode := service.validateInvitation(ctx, invitationId, userId)
	if errCode != "" {
		return errCode
	}
	err := service.friendRepository.CreateCommand(ctx, &entity.Friend{
		UserID1: invitation.ReceiverID,
		UserID2: invitation.SenderID,
	})
	if err != nil {
		log.Error("InvitationFriendService.AcceptInvitation create friend error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	err = service.invitationFriendRepository.DeleteByIDCommand(ctx, invitationId)
	if err != nil {
		log.Error("InvitationFriendService.AcceptInvitation delete invitation error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	return ""
}

func (service *InvitationFriendService) DenyInvitation(ctx *gin.Context, invitationId int64, userId int64) string {
	invitation, errCode := service.validateInvitation(ctx, invitationId, userId)
	if errCode != "" {
		return errCode
	}

	// Create cooldown between users
	currentTime := time.Now().UnixMilli()

	err := service.invitationCooldownRepository.CreateCommand(ctx, &entity.InvitationCooldown{
		UserID1:             invitation.SenderID,
		UserID2:             invitation.ReceiverID,
		StartCooldownMillis: currentTime,
		CooldownDuration:    constants.InvitationCooldownDuration,
	})
	if err != nil {
		log.Error("InvitationFriendService.DenyInvitation create cooldown error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	// Delete the invitation
	err = service.invitationFriendRepository.DeleteByIDCommand(ctx, invitationId)
	if err != nil {
		log.Error("InvitationFriendService.DenyInvitation delete invitation error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	return ""
}

func (service *InvitationFriendService) IsInCooldown(ctx *gin.Context, userId1, userId2 int64) bool {
	cooldown, err := service.invitationCooldownRepository.GetLatestCooldownBetweenUsersQuery(ctx, userId1, userId2)
	if err != nil {
		// If no cooldown record exists, return false
		return false
	}

	currentTime := time.Now().UnixMilli()
	cooldownEndTime := cooldown.StartCooldownMillis + cooldown.CooldownDuration

	// Check if cooldown period has ended
	if currentTime >= cooldownEndTime {
		return false
	}

	return true
}

func (service *InvitationFriendService) WithdrawInvitation(ctx *gin.Context, invitationId int64, userId int64) string {
	// Get the invitation
	invitation, err := service.invitationFriendRepository.GetOneByIDQuery(ctx, invitationId)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return error_utils.ErrorCode.FRIEND_INVITATION_NOT_FOUND
		}
		log.Error("InvitationFriendService.WithdrawInvitation get invitation error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	// Check if the current user is the sender
	if userId != invitation.SenderID {
		return error_utils.ErrorCode.FRIEND_INVITATION_ONLY_SENDER_CAN_WITHDRAW
	}

	// Delete the invitation
	err = service.invitationFriendRepository.DeleteByIDCommand(ctx, invitationId)
	if err != nil {
		log.Error("InvitationFriendService.WithdrawInvitation delete invitation error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	return ""
}
