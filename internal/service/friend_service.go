package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type FriendService interface {
	GetAllFriends(ctx *gin.Context, userId int64) ([]model.FriendResponse, error)
	RemoveFriend(ctx *gin.Context, userId int64, friendId int64) error
}
