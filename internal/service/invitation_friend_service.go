package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type InvitationFriendService interface {
	AddFriend(ctx *gin.Context, invitation model.InvitationFriendRequest, userId int64) error
	GetAllReceivedInvitations(ctx *gin.Context, userId int64) ([]model.InvitationFriendReceivedResponse, error)
	GetAllRequestedInvitations(ctx *gin.Context, userId int64) ([]model.InvitationFriendRequestedResponse, error)
	AcceptInvitation(ctx *gin.Context, invitationId int64, userId int64) error
	DenyInvitation(ctx *gin.Context, invitationId int64, userId int64) error
	IsInCooldown(ctx *gin.Context, userId1, userId2 int64) (bool, error)
	WithdrawInvitation(ctx *gin.Context, invitationId int64, userId int64) error
}
