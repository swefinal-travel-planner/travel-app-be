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
}

func NewTripItemService(
	tripItemRepository repository.TripItemRepository,
) service.TripItemService {
	return &TripItemService{
		tripItemRepository: tripItemRepository,
	}
}

func (service *TripItemService) CreateTripItem(ctx *gin.Context, tripItemRequest model.TripItemRequest) string {
	tripItem := &entity.TripItem{
		TripID:     tripItemRequest.TripID,
		PlaceID:    tripItemRequest.PlaceID,
		TripDay:    tripItemRequest.TripDay,
		OrderInDay: tripItemRequest.OrderInDay,
		Tag:        tripItemRequest.Tag,
	}
	err := service.tripItemRepository.CreateCommand(ctx, tripItem)
	if err != nil {
		log.Error("TripItemService.CreateTripItem Error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	return ""
}
