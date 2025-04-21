package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type TripItemService interface {
	CreateTripItem(ctx *gin.Context, tripItemRequest model.TripItemRequest) string
}
