package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type UserService interface {
	SearchUser(ctx *gin.Context, userEmail string) (*model.FriendResponse, string)
}
