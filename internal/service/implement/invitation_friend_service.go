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
}

func NewInvitationFriendService(
	invitationFriendRepository repository.InvitationFriendRepository,
	userRepository repository.UserRepository,
) service.InvitationFriendService {
	return &InvitationFriendService{
		invitationFriendRepository: invitationFriendRepository,
		userRepository:             userRepository,
	}
}

func (service *InvitationFriendService) AddFriend(ctx *gin.Context, invitation model.InvitationFriendRequest, userId int64) error {
	if userId == invitation.ReceiverID {
		return errors.New("cannot add yourself")
	}
	invitations, err := service.invitationFriendRepository.GetBySenderAndReceiverIdQuery(ctx, userId, invitation.ReceiverID)
	if err != nil {
		return err
	}
	if invitations.Status == "pending" {
		return errors.New("invitation already sent")
	}
	err = service.invitationFriendRepository.CreateCommand(ctx, &entity.InvitationFriend{
		SenderID:   userId,
		ReceiverID: invitation.ReceiverID,
		Status:     "pending",
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *InvitationFriendService) GetAllInvitations(ctx *gin.Context, userId int64) ([]model.InvitationFriendResponse, error) {
	invitations, err := service.invitationFriendRepository.GetByReceiverIdCommand(ctx, userId)
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
