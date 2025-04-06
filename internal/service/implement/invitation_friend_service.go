package serviceimplement

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
)

type InvitationFriendService struct {
	invitationFriendRepository repository.InvitationFriendRepository
	userRepository             repository.UserRepository
	friendRepository           repository.FriendRepository
}

func NewInvitationFriendService(
	invitationFriendRepository repository.InvitationFriendRepository,
	userRepository repository.UserRepository,
	friendRepository repository.FriendRepository,
) service.InvitationFriendService {
	return &InvitationFriendService{
		invitationFriendRepository: invitationFriendRepository,
		userRepository:             userRepository,
		friendRepository:           friendRepository,
	}
}

func (service *InvitationFriendService) AddFriend(ctx *gin.Context, invitation model.InvitationFriendRequest, userId int64) error {
	if userId == invitation.ReceiverID {
		return errors.New("cannot add yourself")
	}
	_, err := service.invitationFriendRepository.GetBySenderAndReceiverIdQuery(ctx, userId, invitation.ReceiverID)
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
	_, err := service.validateInvitation(ctx, invitationId, userId)
	if err != nil {
		return err
	}
	err = service.invitationFriendRepository.DeleteByIDCommand(ctx, invitationId)
	if err != nil {
		return err
	}
	return nil
}
