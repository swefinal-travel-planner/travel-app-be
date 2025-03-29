package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type InvitationFriendService interface {
	AddFriend(ctx *gin.Context, invitation model.InvitationFriendRequest, userId int64) error
	GetAllInvitations(ctx *gin.Context, userId int64) ([]model.InvitationFriendResponse, error)
}
