package serviceimplement

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
	stringlistutils "github.com/swefinal-travel-planner/travel-app-be/internal/utils/string_list_utils"
)

type TripService struct {
	tripRepository  repository.TripRepository
	tripItemService service.TripItemService
}

func NewTripService(
	tripRepository repository.TripRepository,
	tripItemService service.TripItemService,
) service.TripService {
	return &TripService{
		tripRepository:  tripRepository,
		tripItemService: tripItemService,
	}
}

func (service *TripService) CreateTrip(ctx *gin.Context, tripRequest model.TripRequest) (int64, string) {
	trip := &entity.Trip{
		Title:                 tripRequest.Title,
		City:                  tripRequest.City,
		StartDate:             tripRequest.StartDate,
		Days:                  tripRequest.Days,
		Budget:                tripRequest.Budget,
		NumMembers:            tripRequest.NumMembers,
		ViLocationAttributes:  tripRequest.ViLocationAttributes,
		ViFoodAttributes:      tripRequest.ViFoodAttributes,
		ViSpecialRequirements: tripRequest.ViSpecialRequirements,
		ViMedicalConditions:   tripRequest.ViMedicalConditions,
		EnLocationAttributes:  tripRequest.EnLocationAttributes,
		EnFoodAttributes:      tripRequest.EnFoodAttributes,
		EnSpecialRequirements: tripRequest.EnSpecialRequirements,
		EnMedicalConditions:   tripRequest.EnMedicalConditions,
		Status:                "in_progress",
	}
	tripID, err := service.tripRepository.CreateCommand(ctx, trip)
	if err != nil {
		log.Error("TripService.CreateTrip Error: " + err.Error())
		return 0, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	return tripID, ""
}

func validValue(fields ...stringlistutils.StringList) stringlistutils.StringList {
	for _, f := range fields {
		if len(f) > 0 {
			return f
		}
	}
	return nil // All are empty
}

func (service *TripService) CreateTripByAI(ctx *gin.Context, tripRequest model.TripRequest) ([]model.TripItemResponse, string) {
	tripToCoreRequest := model.TripToCoreRequest{
		City:                tripRequest.City,
		Days:                tripRequest.Days,
		LocationAttributes:  validValue(tripRequest.ViLocationAttributes, tripRequest.EnLocationAttributes),
		FoodAttributes:      validValue(tripRequest.ViFoodAttributes, tripRequest.EnFoodAttributes),
		SpecialRequirements: validValue(tripRequest.ViSpecialRequirements, tripRequest.EnSpecialRequirements),
		MedicalConditions:   validValue(tripRequest.ViMedicalConditions, tripRequest.EnMedicalConditions),
	}
	fmt.Println("tripToCoreRequest: ", tripToCoreRequest)

	// send trip data to core service
	//...
	// receive trip item from core service
	tripItemsRespFromCore := []model.TripItemResponse{
		{
			PlaceID:    "place1",
			TripDay:    1,
			OrderInDay: 1,
			Tag:        "food_location",
		},
		{
			PlaceID:    "place2",
			TripDay:    1,
			OrderInDay: 2,
			Tag:        "travel_location",
		},
		{
			PlaceID:    "place3",
			TripDay:    2,
			OrderInDay: 1,
			Tag:        "food_location",
		},
		{
			PlaceID:    "place4",
			TripDay:    2,
			OrderInDay: 2,
			Tag:        "travel_location",
		},
	}

	// save trip to database
	tripID, errCode := service.CreateTrip(ctx, tripRequest)
	if errCode != "" {
		return []model.TripItemResponse{}, errCode
	}

	// add tripID to trip items
	for i := range tripItemsRespFromCore {
		tripItemsRespFromCore[i].TripID = tripID
	}

	// save trip items to database
	for _, tripItemResp := range tripItemsRespFromCore {
		tripItemReq := model.TripItemRequest{
			TripID:     tripItemResp.TripID,
			PlaceID:    tripItemResp.PlaceID,
			TripDay:    tripItemResp.TripDay,
			OrderInDay: tripItemResp.OrderInDay,
			Tag:        tripItemResp.Tag,
		}
		err := service.tripItemService.CreateTripItem(ctx, tripItemReq)
		if err != "" {
			log.Error("TripService.CreateTripByAI Error: " + err)
			return []model.TripItemResponse{}, err
		}
	}

	// return trip items to client
	return tripItemsRespFromCore, ""
}
