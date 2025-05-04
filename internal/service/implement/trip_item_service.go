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
	tripItemRepository repository.TripItemRepository
	tripRepository     repository.TripRepository
	unitOfWork         repository.UnitOfWork
}

func NewTripItemService(
	tripItemRepository repository.TripItemRepository,
	tripRepository repository.TripRepository,
	unitOfWork repository.UnitOfWork,
) service.TripItemService {
	return &TripItemService{
		tripItemRepository: tripItemRepository,
		tripRepository:     tripRepository,
		unitOfWork:         unitOfWork,
	}
}

func (service *TripItemService) CreateTripItems(ctx *gin.Context, tripItemRequests []model.TripItemRequest) string {
	// begin transaction
	tx, err := service.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("TripItemService.CreateTripItems Begin error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	defer service.unitOfWork.Rollback(tx)

	// check to see if trip exists
	_, err = service.tripRepository.GetOneByIDQuery(ctx, tripItemRequests[0].TripID, tx)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return error_utils.ErrorCode.TRIP_NOT_FOUND
		}
		log.Error("TripItemService.CreateTripItems GetOneByIDQuery error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	// lock trip row before update trip items
	tripID := tripItemRequests[0].TripID
	_, err = service.tripRepository.LockTripRowByIDCommand(ctx, tripID, tx)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
			return error_utils.ErrorCode.TRIP_NOT_FOUND
		}
		log.Error("TripItemService.CreateTripItems LockTripRowByIDCommand error: " + err.Error())
		return error_utils.ErrorCode.DB_DOWN
	}

	// delete existing trip items
	err = service.tripItemRepository.DeleteByTripIDCommand(ctx, tripID, tx)
	if err != nil {
		log.Error("TripItemService.CreateTripItems DeleteByTripIDCommand error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// insert new trip items
	for _, tripItemRequest := range tripItemRequests {
		tripItem := &entity.TripItem{
			TripID:     tripItemRequest.TripID,
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
