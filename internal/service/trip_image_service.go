package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type TripImageService interface {
	CreateTripImage(ctx *gin.Context, userId int64, tripId int64, tripImageRequest model.TripImageRequest) string
	GetTripImages(ctx *gin.Context, userId int64, tripId int64) ([]model.TripImageResponse, string)
	GetTripImagesWithUserInfo(ctx *gin.Context, userId int64, tripId int64) ([]model.TripImageWithUserInfoResponse, string)
	DeleteTripImage(ctx *gin.Context, userId int64, tripId int64, imageId int64) string
}
