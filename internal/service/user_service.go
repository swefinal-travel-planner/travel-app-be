package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type UserService interface {
	SearchUser(ctx *gin.Context, userId int64, searchTerm string) ([]model.FriendResponse, string)
	UpdateNotificationToken(ctx *gin.Context, userId int64, notificationTokenRequest model.UpdateNotificationTokenRequest) string
	UpdateUser(ctx *gin.Context, userId int64, request model.UpdateUserRequest) string
	GetUserInfo(ctx *gin.Context, userId int64) (*model.UserInfoResponse, string)
}
