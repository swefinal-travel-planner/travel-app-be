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

func (service *InvitationFriendService) AddFriend(ctx *gin.Context, invitation model.InvitationFriendRequest) error {
	userId, exists := ctx.Get("userId")
	if !exists {
		return errors.New("user not exists")
	}
	err := service.invitationFriendRepository.CreateCommand(ctx, &entity.InvitationFriend{
		SenderID:   userId.(int64),
		ReceiverID: invitation.ReceiverID,
		Status:     "pending",
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *InvitationFriendService) GetAllInvitations(ctx *gin.Context) ([]model.InvitationFriendResponse, error) {
	userId, exists := ctx.Get("userId")
	if !exists {
		return nil, errors.New("user not exists")
	}
	invitations, err := service.invitationFriendRepository.GetByReceiverIdCommand(ctx, userId.(int64))
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
