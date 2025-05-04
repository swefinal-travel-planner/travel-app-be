package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type TripItemService interface {
	CreateTripItems(ctx *gin.Context, userId int64, tripId int64, tripItemRequests []model.TripItemRequest) string
}
