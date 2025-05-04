package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type TripService interface {
	CreateTrip(ctx *gin.Context, tripRequest model.TripRequest, userId int64) (int64, string)
	GetTrip(ctx *gin.Context, tripId int64) (*model.TripResponse, string)
	GetTripsByUserId(ctx *gin.Context, userId int64) ([]model.TripResponse, string)
}
