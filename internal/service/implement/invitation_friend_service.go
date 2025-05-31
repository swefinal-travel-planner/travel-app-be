package serviceimplement

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"

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
	notificationService          service.NotificationService
}

func NewInvitationFriendService(
	invitationFriendRepository repository.InvitationFriendRepository,
	userRepository repository.UserRepository,
	friendRepository repository.FriendRepository,
	invitationCooldownRepository repository.InvitationCooldownRepository,
	notificationService service.NotificationService,
) service.InvitationFriendService {
	return &InvitationFriendService{
		invitationFriendRepository:   invitationFriendRepository,
		userRepository:               userRepository,
		friendRepository:             friendRepository,
		invitationCooldownRepository: invitationCooldownRepository,
		notificationService:          notificationService,
	}
}

func (service *InvitationFriendService) AddFriend(ctx *gin.Context, invitation model.InvitationFriendRequest, userId int64) string {
	friend, err := service.userRepository.GetOneByEmailQuery(ctx, invitation.ReceiverEmail, nil)
	if err != nil {
		log.Error("InvitationFriendService.AddFriend GetOneByEmailQuery error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if friend == nil {
		return error_utils.ErrorCode.ADD_FRIEND_RECEIVER_NOT_FOUND
	}
	if userId == friend.Id {
		return error_utils.ErrorCode.ADD_FRIEND_RECEIVER_NOT_FOUND
	}

	// Check if users are in cooldown period
	inCooldown, _ := service.IsInCoolDownAndGetRemainingTime(ctx, userId, friend.Id)
	if inCooldown {
		return error_utils.ErrorCode.ADD_FRIEND_IN_COOLDOWN
	}

	_, err = service.invitationFriendRepository.GetBySenderAndReceiverIdQuery(ctx, userId, friend.Id, nil)
	if err == nil {
		return error_utils.ErrorCode.ADD_FRIEND_INVITATION_ALREADY_EXISTS
	}

	isFriend := service.friendRepository.ExistsByUserId1AndUserId2Query(ctx, userId, friend.Id, nil)
	if isFriend {
		return error_utils.ErrorCode.ADD_FRIEND_ALREADY_FRIEND
	}

	invitationEntity := entity.InvitationFriend{
		SenderID:   userId,
		ReceiverID: friend.Id,
	}

	err = service.invitationFriendRepository.CreateCommand(ctx, &invitationEntity, nil)

	if err != nil {
		log.Error("InvitationFriendService.AddFriend CreateCommand error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	errCode := service.notificationService.SaveAndSendNotification(ctx, model.SaveNotificationRequest{
		Type:                entity.NotificationType.FriendRequestReceived,
		ReceiverUserID:      friend.Id,
		TriggerEntityType:   entity.NotificationTriggerType.User,
		TriggerEntityID:     &userId,
		ReferenceEntityType: entity.NotificationReferenceType.FriendInvitation,
		ReferenceEntityID:   &invitationEntity.ID,
	})

	return errCode
}

// PROBLEM: N+1
func (service *InvitationFriendService) GetAllReceivedInvitations(ctx *gin.Context, userId int64) ([]model.InvitationFriendReceivedResponse, string) {
	invitations, err := service.invitationFriendRepository.GetByReceiverIdQuery(ctx, userId, nil)
	if err != nil {
		log.Error("InvitationFriendService.GetAllReceivedInvitations GetByReceiverIdQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	if len(invitations) == 0 {
		return make([]model.InvitationFriendReceivedResponse, 0), ""
	}

	var invitationResponses []model.InvitationFriendReceivedResponse
	for _, invitation := range invitations {
		user, err := service.userRepository.GetOneByIDQuery(ctx, invitation.SenderID, nil)
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
	invitations, err := service.invitationFriendRepository.GetBySenderIdQuery(ctx, userId, nil)
	if err != nil {
		log.Error("InvitationFriendService.GetAllRequestedInvitations GetBySenderIdQuery error: ", err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	if len(invitations) == 0 {
		return make([]model.InvitationFriendRequestedResponse, 0), ""
	}

	var invitationResponses []model.InvitationFriendRequestedResponse
	for _, invitation := range invitations {
		user, err := service.userRepository.GetOneByIDQuery(ctx, invitation.SenderID, nil)
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
	invitation, err := service.invitationFriendRepository.GetOneByIDQuery(ctx, invitationId, nil)
	if err != nil {
		log.Error("InvitationFriendService.validateInvitation GetOneByIDQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	if invitation == nil {
		return nil, error_utils.ErrorCode.FRIEND_INVITATION_NOT_FOUND
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
	}, nil)
	if err != nil {
		log.Error("InvitationFriendService.AcceptInvitation create friend error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	err = service.invitationFriendRepository.DeleteByIDCommand(ctx, invitationId, nil)
	if err != nil {
		log.Error("InvitationFriendService.AcceptInvitation delete invitation error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	// PROBLEM: NOT SAFE TRANSACTION
	service.notificationService.SaveAndSendNotification(ctx, model.SaveNotificationRequest{
		Type:                entity.NotificationType.FriendRequestAccepted,
		ReceiverUserID:      invitation.SenderID,
		TriggerEntityType:   entity.NotificationTriggerType.User,
		TriggerEntityID:     &userId,
		ReferenceEntityType: entity.NotificationReferenceType.FriendInvitation,
		ReferenceEntityID:   &invitationId,
	})

	service.notificationService.DeleteFriendInvitation(ctx, invitation.ReceiverID, entity.NotificationType.FriendRequestReceived, invitation.SenderID)

	return errCode
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
	}, nil)
	if err != nil {
		log.Error("InvitationFriendService.DenyInvitation create cooldown error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	// Delete the invitation
	err = service.invitationFriendRepository.DeleteByIDCommand(ctx, invitationId, nil)
	if err != nil {
		log.Error("InvitationFriendService.DenyInvitation delete invitation error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	service.notificationService.DeleteFriendInvitation(ctx, invitation.ReceiverID, entity.NotificationType.FriendRequestReceived, invitation.SenderID)

	return ""
}

func (service *InvitationFriendService) IsInCoolDownAndGetRemainingTime(ctx *gin.Context, userId1, userId2 int64) (bool, int64) {
	cooldown, err := service.invitationCooldownRepository.GetLatestCooldownBetweenUsersQuery(ctx, userId1, userId2, nil)
	if err != nil {
		// If no cooldown record exists, return false
		return false, 0
	}

	currentTime := time.Now().UnixMilli()
	cooldownEndTime := cooldown.StartCooldownMillis + cooldown.CooldownDuration

	// Check if cooldown period has ended
	if currentTime > cooldownEndTime {
		return false, 0
	}

	return true, cooldownEndTime - currentTime
}

func (service *InvitationFriendService) WithdrawInvitation(ctx *gin.Context, invitationId int64, userId int64) string {
	// Get the invitation
	invitation, err := service.invitationFriendRepository.GetOneByIDQuery(ctx, invitationId, nil)
	if err != nil {
		log.Error("InvitationFriendService.WithdrawInvitation get invitation error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if invitation == nil {
		return error_utils.ErrorCode.FRIEND_INVITATION_NOT_FOUND
	}

	// Check if the current user is the sender
	if userId != invitation.SenderID {
		return error_utils.ErrorCode.FORBIDDEN
	}

	// Delete the invitation
	err = service.invitationFriendRepository.DeleteByIDCommand(ctx, invitationId, nil)
	if err != nil {
		log.Error("InvitationFriendService.WithdrawInvitation delete invitation error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	service.notificationService.DeleteFriendInvitation(ctx, invitation.ReceiverID, entity.NotificationType.FriendRequestReceived, invitation.SenderID)

	return ""
}
