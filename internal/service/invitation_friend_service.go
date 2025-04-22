package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type InvitationFriendService interface {
	AddFriend(ctx *gin.Context, invitation model.InvitationFriendRequest, userId int64) string
	GetAllReceivedInvitations(ctx *gin.Context, userId int64) ([]model.InvitationFriendReceivedResponse, string)
	GetAllRequestedInvitations(ctx *gin.Context, userId int64) ([]model.InvitationFriendRequestedResponse, string)
	AcceptInvitation(ctx *gin.Context, invitationId int64, userId int64) string
	DenyInvitation(ctx *gin.Context, invitationId int64, userId int64) string
	IsInCooldown(ctx *gin.Context, userId1, userId2 int64) bool
	GetCooldownRemainingAsSender(ctx *gin.Context, userId1, userId2 int64) int64
	WithdrawInvitation(ctx *gin.Context, invitationId int64, userId int64) string
}
