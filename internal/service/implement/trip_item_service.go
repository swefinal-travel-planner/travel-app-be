package serviceimplement

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/env"
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

	// check if user is admin
	isAdmin, err := service.tripMemberRepository.IsUserTripAdminQuery(ctx, tripId, userId, tx)
	if err != nil {
		log.Error("TripItemService.CreateTripItems IsUserTripAdminQuery error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if !isAdmin {
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

func (service *TripItemService) GetTripItemsByTripID(ctx *gin.Context, userId int64, tripId int64) ([]model.TripItemResponse, string) {
	lang := ctx.DefaultQuery("language", "vi")

	// Get trip items with membership check
	tripItems, err := service.tripItemRepository.GetTripItemsByTripIDCommand(ctx, tripId, userId, nil)
	if err != nil {
		if err.Error() == error_utils.SystemErrorMessage.NotMemberOfTrip {
			return nil, error_utils.ErrorCode.FORBIDDEN
		}
		log.Error("TripItemService.GetTripItemsByTripID Error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// Convert to response model
	var tripItemResponses []model.TripItemResponse
	for _, item := range tripItems {
		tripItemResponse := model.TripItemResponse{
			ID:         item.ID,
			TripID:     item.TripID,
			PlaceID:    item.PlaceID,
			TripDay:    item.TripDay,
			OrderInDay: item.OrderInDay,
			TimeInDate: item.TimeInDate,
		}

		// Fetch place info from external API
		placeInfo, err := fetchPlaceInfo(item.PlaceID, lang)
		if err == nil {
			tripItemResponse.PlaceInfo = placeInfo
		} else {
			tripItemResponse.PlaceInfo = nil
		}

		tripItemResponses = append(tripItemResponses, tripItemResponse)
	}

	return tripItemResponses, ""
}

// fetchPlaceInfo calls the external API and returns *model.PlaceInfo or error
func fetchPlaceInfo(placeID string, lang string) (*model.PlaceInfo, error) {
	apiRoute, err := env.GetEnv("PLACE_INFO_URL")
	if err != nil {
		log.Error("fetchPlaceInfo - Get PLACE_INFO_URL Error: " + err.Error())
		return nil, err
	}
	url := fmt.Sprintf("%s/%s?language=%s", apiRoute, placeID, lang)

	secretKey, err := env.GetEnv("CORE_SECRET_KEY")
	if err != nil {
		log.Error("fetchPlaceInfo - Get CORE_SECRET_KEY Error: " + err.Error())
		return nil, err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+secretKey)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var apiResp struct {
		Data   model.PlaceInfo `json:"data"`
		Status int             `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return &apiResp.Data, nil
}
