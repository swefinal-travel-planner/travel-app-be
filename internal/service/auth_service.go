package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type AuthService interface {
	Register(ctx *gin.Context, userRequest model.RegisterRequest) error
}
