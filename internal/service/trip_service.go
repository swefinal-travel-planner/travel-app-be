package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type TripService interface {
	CreateTrip(ctx *gin.Context, tripRequest model.TripRequest, userId int64) (int64, string)
	GetAllTripsByUserID(ctx *gin.Context, userId int64) ([]*model.TripResponse, string)
	GetTripByID(ctx *gin.Context, tripId int64, userId int64) (*model.TripResponse, string)
	UpdateTrip(ctx *gin.Context, tripId int64, userId int64, tripRequest model.TripPatchRequest) string
	CreateTripByAI(ctx *gin.Context, tripRequest model.TripRequest, userID int64) ([]model.TripItemFromAIResponse, string)
}
