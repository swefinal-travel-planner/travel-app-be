package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type UserService interface {
	SearchUser(ctx *gin.Context, userId int64, userEmail string) (*model.FriendResponse, string)
	UpdateNotificationToken(ctx *gin.Context, userId int64, notificationTokenRequest model.UpdateNotificationTokenRequest) string
}
