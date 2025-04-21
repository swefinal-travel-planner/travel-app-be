package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type TripService interface {
	CreateTrip(ctx *gin.Context, tripRequest model.TripRequest) (int64, string)
	CreateTripByAI(ctx *gin.Context, tripRequest model.TripRequest) ([]model.TripItemResponse, string)
}
