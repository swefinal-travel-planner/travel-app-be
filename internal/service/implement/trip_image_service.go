package serviceimplement

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
)

type TripImageService struct {
	tripImageRepository  repository.TripImageRepository
	tripRepository       repository.TripRepository
	tripMemberRepository repository.TripMemberRepository
	tripItemRepository   repository.TripItemRepository
	unitOfWork           repository.UnitOfWork
}

func NewTripImageService(
	tripImageRepository repository.TripImageRepository,
	tripRepository repository.TripRepository,
	tripMemberRepository repository.TripMemberRepository,
	tripItemRepository repository.TripItemRepository,
	unitOfWork repository.UnitOfWork,
) service.TripImageService {
	return &TripImageService{
		tripImageRepository:  tripImageRepository,
		tripRepository:       tripRepository,
		tripMemberRepository: tripMemberRepository,
		tripItemRepository:   tripItemRepository,
		unitOfWork:           unitOfWork,
	}
}

func (service *TripImageService) CreateTripImage(ctx *gin.Context, userId int64, tripId int64, tripImageRequest model.TripImageRequest) string {
	// begin transaction
	tx, err := service.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("TripImageService.CreateTripImage Begin error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	defer service.unitOfWork.Rollback(tx)

	// check if trip exists
	trip, err := service.tripRepository.GetOneByIDQuery(ctx, tripId, tx)
	if err != nil {
		log.Error("TripImageService.CreateTripImage GetOneByIDQuery error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if trip == nil {
		return error_utils.ErrorCode.TRIP_NOT_FOUND
	}

	// check if user is member of the trip
	isMember, err := service.tripMemberRepository.IsUserInTripQuery(ctx, tripId, userId, tx)
	if err != nil {
		log.Error("TripImageService.CreateTripImage IsUserInTripQuery error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if !isMember {
		return error_utils.ErrorCode.FORBIDDEN
	}

	if tripImageRequest.TripItemID != nil {
		isTripItemExists, err := service.tripItemRepository.ExistsByTripIDAndTripItemIDCommand(ctx, tripId, *tripImageRequest.TripItemID, tx)
		if err != nil {
			log.Error("TripImageService.CreateTripImage GetOneByIDQuery error: " + err.Error())
			return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
		}
		if !isTripItemExists {
			return error_utils.ErrorCode.BAD_REQUEST
		}
	}

	// create trip image
	tripImage := &entity.TripImage{
		TripID:     tripId,
		ImageURL:   tripImageRequest.ImageURL,
		UserID:     userId,
		TripItemID: tripImageRequest.TripItemID,
	}
	err = service.tripImageRepository.CreateCommand(ctx, tripImage, tx)
	if err != nil {
		log.Error("TripImageService.CreateTripImage CreateCommand error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// commit transaction
	err = service.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("TripImageService.CreateTripImage Commit error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	return ""
}

func (service *TripImageService) GetTripImages(ctx *gin.Context, userId int64, tripId int64) ([]model.TripImageResponse, string) {
	// check if trip exists
	trip, err := service.tripRepository.GetOneByIDQuery(ctx, tripId, nil)
	if err != nil {
		log.Error("TripImageService.GetTripImages GetOneByIDQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	if trip == nil {
		return nil, error_utils.ErrorCode.TRIP_NOT_FOUND
	}

	// check if user is member of the trip
	isMember, err := service.tripMemberRepository.IsUserInTripQuery(ctx, tripId, userId, nil)
	if err != nil {
		log.Error("TripImageService.GetTripImages IsUserInTripQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if !isMember {
		return nil, error_utils.ErrorCode.FORBIDDEN
	}

	// get trip images
	tripImages, err := service.tripImageRepository.GetAllQuery(ctx, tripId, nil)
	if err != nil {
		log.Error("TripImageService.GetTripImages GetAllQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// convert to response model
	var tripImageResponses []model.TripImageResponse
	for _, image := range tripImages {
		tripImageResponse := model.TripImageResponse{
			ID:         image.ID,
			TripID:     image.TripID,
			TripItemID: image.TripItemID,
			ImageURL:   image.ImageURL,
			CreatedAt:  image.CreatedAt,
		}
		tripImageResponses = append(tripImageResponses, tripImageResponse)
	}

	return tripImageResponses, ""
}

func (service *TripImageService) DeleteTripImage(ctx *gin.Context, userId int64, tripId int64, imageId int64) string {
	// begin transaction
	tx, err := service.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("TripImageService.DeleteTripImage Begin error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	defer service.unitOfWork.Rollback(tx)

	// check if trip exists
	trip, err := service.tripRepository.GetOneByIDQuery(ctx, tripId, tx)
	if err != nil {
		log.Error("TripImageService.DeleteTripImage GetOneByIDQuery error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}
	if trip == nil {
		return error_utils.ErrorCode.TRIP_NOT_FOUND
	}

	// check if user is admin of the trip
	isAdmin, err := service.tripMemberRepository.IsUserTripAdminQuery(ctx, tripId, userId, tx)
	if err != nil {
		log.Error("TripImageService.DeleteTripImage IsUserTripAdminQuery error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if !isAdmin {
		return error_utils.ErrorCode.FORBIDDEN
	}

	// delete trip image
	err = service.tripImageRepository.DeleteOneByIDCommand(ctx, imageId, tx)
	if err != nil {
		log.Error("TripImageService.DeleteTripImage DeleteOneByIDCommand error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// commit transaction
	err = service.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("TripImageService.DeleteTripImage Commit error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	return ""
}

func (service *TripImageService) GetTripImagesWithUserInfo(ctx *gin.Context, userId int64, tripId int64) ([]model.TripImageWithUserInfoResponse, string) {
	// check if trip exists
	trip, err := service.tripRepository.GetOneByIDQuery(ctx, tripId, nil)
	if err != nil {
		log.Error("TripImageService.GetTripImagesWithUserInfo GetOneByIDQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	if trip == nil {
		return nil, error_utils.ErrorCode.TRIP_NOT_FOUND
	}

	// check if user is member of the trip
	isMember, err := service.tripMemberRepository.IsUserInTripQuery(ctx, tripId, userId, nil)
	if err != nil {
		log.Error("TripImageService.GetTripImagesWithUserInfo IsUserInTripQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if !isMember {
		return nil, error_utils.ErrorCode.FORBIDDEN
	}

	// get trip images with user info
	tripImagesWithUserInfo, err := service.tripImageRepository.GetAllWithUserInfoQuery(ctx, tripId, nil)
	if err != nil {
		log.Error("TripImageService.GetTripImagesWithUserInfo GetAllWithUserInfoQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// convert to response model
	var tripImageResponses []model.TripImageWithUserInfoResponse
	for _, image := range tripImagesWithUserInfo {
		tripImageResponse := model.TripImageWithUserInfoResponse{
			ID:         image.ID,
			TripID:     image.TripID,
			TripItemID: image.TripItemID,
			ImageURL:   image.ImageURL,
			CreatedAt:  image.CreatedAt,
			Author: model.UserInfo{
				ID:       image.UserID,
				Name:     *image.UserName,
				PhotoURL: image.UserPhotoUrl,
			},
		}
		tripImageResponses = append(tripImageResponses, tripImageResponse)
	}

	return tripImageResponses, ""
}

func (service *TripImageService) GetAllByTripIDAndTripItemID(ctx *gin.Context, userId int64, tripID int64, tripItemID int64) ([]model.TripImageResponse, string) {
	// check if trip exists
	trip, err := service.tripRepository.GetOneByIDQuery(ctx, tripID, nil)
	if err != nil {
		log.Error("TripImageService.GetAllByTripIDAndTripItemID GetOneByIDQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.DB_DOWN
	}
	if trip == nil {
		return nil, error_utils.ErrorCode.TRIP_NOT_FOUND
	}

	// check if user is member of the trip
	isMember, err := service.tripMemberRepository.IsUserInTripQuery(ctx, tripID, userId, nil)
	if err != nil {
		log.Error("TripImageService.GetAllByTripIDAndTripItemID IsUserInTripQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if !isMember {
		return nil, error_utils.ErrorCode.FORBIDDEN
	}

	isTripItemExists, err := service.tripItemRepository.ExistsByTripIDAndTripItemIDCommand(ctx, tripID, tripItemID, nil)
	if err != nil {
		log.Error("TripImageService.GetAllByTripIDAndTripItemID ExistsByTripIDAndTripItemIDCommand error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if !isTripItemExists {
		return nil, error_utils.ErrorCode.BAD_REQUEST
	}

	// get trip images
	tripImages, err := service.tripImageRepository.GetAllByTripIDAndTripItemIDQuery(ctx, tripID, tripItemID, nil)
	if err != nil {
		log.Error("TripImageService.GetAllByTripIDAndTripItemID GetAllByTripIDAndTripItemIDQuery error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// convert to response model
	var tripImageResponses []model.TripImageResponse
	for _, image := range tripImages {
		tripImageResponse := model.TripImageResponse{
			ID:         image.ID,
			TripID:     image.TripID,
			TripItemID: image.TripItemID,
			ImageURL:   image.ImageURL,
			CreatedAt:  image.CreatedAt,
		}
		tripImageResponses = append(tripImageResponses, tripImageResponse)
	}

	return tripImageResponses, ""
}
