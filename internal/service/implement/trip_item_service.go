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

type TripItemService struct {
	tripItemRepository   repository.TripItemRepository
	tripRepository       repository.TripRepository
	tripMemberRepository repository.TripMemberRepository
	unitOfWork           repository.UnitOfWork
}

func NewTripItemService(
	tripItemRepository repository.TripItemRepository,
	tripRepository repository.TripRepository,
	tripMemberRepository repository.TripMemberRepository,
	unitOfWork repository.UnitOfWork,
) service.TripItemService {
	return &TripItemService{
		tripItemRepository:   tripItemRepository,
		tripRepository:       tripRepository,
		tripMemberRepository: tripMemberRepository,
		unitOfWork:           unitOfWork,
	}
}

func (service *TripItemService) CreateTripItems(ctx *gin.Context, userId int64, tripId int64, tripItemRequests []model.TripItemRequest) string {
	// begin transaction
	tx, err := service.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("TripItemService.CreateTripItems Begin error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	defer service.unitOfWork.Rollback(tx)

	// lock trip row before update trip items
	_, err = service.tripRepository.SelectForUpdateById(ctx, tripId, tx)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return error_utils.ErrorCode.TRIP_NOT_FOUND
		}
		log.Error("TripItemService.CreateTripItems LockTripRowByIDCommand error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	// check to see if trip exists
	_, err = service.tripRepository.GetOneByIDQuery(ctx, tripId, tx)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return error_utils.ErrorCode.TRIP_NOT_FOUND
		}
		log.Error("TripItemService.CreateTripItems GetOneByIDQuery error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	// check if user is admin or staff
	isAdminOrStaff, err := service.tripMemberRepository.IsUserTripAdminOrStaffQuery(ctx, tripId, userId, tx)
	if err != nil {
		log.Error("TripItemService.CreateTripItems IsUserTripAdminOrStaffQuery error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if !isAdminOrStaff {
		return error_utils.ErrorCode.FORBIDDEN
	}

	// delete existing trip items
	err = service.tripItemRepository.DeleteByTripIDCommand(ctx, tripId, tx)
	if err != nil {
		log.Error("TripItemService.CreateTripItems DeleteByTripIDCommand error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// insert new trip items
	for _, tripItemRequest := range tripItemRequests {
		tripItem := &entity.TripItem{
			TripID:     tripId,
			PlaceID:    tripItemRequest.PlaceID,
			TripDay:    tripItemRequest.TripDay,
			OrderInDay: tripItemRequest.OrderInDay,
			TimeInDate: tripItemRequest.TimeInDate,
		}
		err := service.tripItemRepository.CreateCommand(ctx, tripItem, tx)
		if err != nil {
			log.Error("TripItemService.CreateTripItems Error: " + err.Error())
			return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
		}
	}

	// commit transaction
	err = service.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("TripItemService.CreateTripItems Commit error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	return ""
}
