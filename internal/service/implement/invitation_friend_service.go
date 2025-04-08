package serviceimplement

import (
	"errors"
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

func (service *InvitationFriendService) AddFriend(ctx *gin.Context, invitation model.InvitationFriendRequest, userId int64) error {
	if userId == invitation.ReceiverID {
		return errors.New("cannot add yourself")
	}

	// Check if users are in cooldown period
	inCooldown, err := service.IsInCooldown(ctx, userId, invitation.ReceiverID)
	if err != nil {
		return err
	}
	if inCooldown {
		return errors.New("cannot send friend request during cooldown period")
	}

	_, err = service.invitationFriendRepository.GetBySenderAndReceiverIdQuery(ctx, userId, invitation.ReceiverID)
	if err == nil {
		return errors.New("invitation already exists")
	}

	err = service.friendRepository.GetByUserId1AndUserId2Query(ctx, userId, invitation.ReceiverID)
	if err == nil {
		return errors.New("already friends")
	}

	err = service.invitationFriendRepository.CreateCommand(ctx, &entity.InvitationFriend{
		SenderID:   userId,
		ReceiverID: invitation.ReceiverID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *InvitationFriendService) GetAllInvitations(ctx *gin.Context, userId int64) ([]model.InvitationFriendResponse, error) {
	invitations, err := service.invitationFriendRepository.GetByReceiverIdQuery(ctx, userId)
	if err != nil {
		return nil, err
	}
	var invitationResponses []model.InvitationFriendResponse
	for _, invitation := range invitations {
		user, err := service.userRepository.GetOneByIDQuery(ctx, invitation.ReceiverID)
		if err != nil {
			return nil, err
		}
		invitationResponses = append(invitationResponses, model.InvitationFriendResponse{
			Id:               invitation.ID,
			ReceiverUsername: user.Name,
			ReceiverImageURL: user.PhotoURL,
		})
	}
	return invitationResponses, nil
}

func (service *InvitationFriendService) validateInvitation(ctx *gin.Context, invitationId int64, userId int64) (*entity.InvitationFriend, error) {
	invitation, err := service.invitationFriendRepository.GetOneByIDQuery(ctx, invitationId)
	if err != nil {
		return nil, err
	}
	if userId == invitation.SenderID {
		return nil, errors.New("cannot process your own invitation")
	}
	if userId != invitation.ReceiverID {
		return nil, errors.New("not your invitation")
	}
	return invitation, nil
}

func (service *InvitationFriendService) AcceptInvitation(ctx *gin.Context, invitationId int64, userId int64) error {
	invitation, err := service.validateInvitation(ctx, invitationId, userId)
	if err != nil {
		return err
	}
	err = service.friendRepository.CreateCommand(ctx, &entity.Friend{
		UserID1: invitation.ReceiverID,
		UserID2: invitation.SenderID,
	})
	if err != nil {
		return err
	}
	err = service.invitationFriendRepository.DeleteByIDCommand(ctx, invitationId)
	if err != nil {
		return err
	}
	return nil
}

func (service *InvitationFriendService) DenyInvitation(ctx *gin.Context, invitationId int64, userId int64) error {
	invitation, err := service.validateInvitation(ctx, invitationId, userId)
	if err != nil {
		return err
	}

	// Create cooldown between users
	currentTime := time.Now().UnixMilli()

	err = service.invitationCooldownRepository.CreateCommand(ctx, &entity.InvitationCooldown{
		UserID1:             invitation.SenderID,
		UserID2:             invitation.ReceiverID,
		StartCooldownMillis: currentTime,
		CooldownDuration:    constants.InvitationCooldownDuration,
	})
	if err != nil {
		return err
	}

	// Delete the invitation
	err = service.invitationFriendRepository.DeleteByIDCommand(ctx, invitationId)
	if err != nil {
		return err
	}

	return nil
}

func (service *InvitationFriendService) IsInCooldown(ctx *gin.Context, userId1, userId2 int64) (bool, error) {
	cooldown, err := service.invitationCooldownRepository.GetLatestCooldownBetweenUsersQuery(ctx, userId1, userId2)
	if err != nil {
		// If no cooldown record exists, return false
		return false, nil
	}

	currentTime := time.Now().UnixMilli()
	cooldownEndTime := cooldown.StartCooldownMillis + cooldown.CooldownDuration

	// Check if cooldown period has ended
	if currentTime >= cooldownEndTime {
		return false, nil
	}

	return true, nil
}
