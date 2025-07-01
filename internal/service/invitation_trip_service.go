package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type InvitationTripService interface {
	SendInvitation(ctx *gin.Context, invitation model.InvitationTripRequest, senderId int64) string
	GetAllReceivedInvitations(ctx *gin.Context, userId int64) ([]model.InvitationTripReceivedResponse, string)
	GetAllSentInvitations(ctx *gin.Context, userId int64) ([]model.InvitationTripSentResponse, string)
	GetPendingInvitationsByTripID(ctx *gin.Context, tripId int64, userId int64) ([]model.InvitationTripPendingResponse, string)
	AcceptInvitation(ctx *gin.Context, invitationId int64, userId int64) string
	DenyInvitation(ctx *gin.Context, invitationId int64, userId int64) string
	WithdrawInvitation(ctx *gin.Context, invitationId int64, userId int64) string
}
