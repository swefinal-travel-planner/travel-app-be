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

type TripService struct {
	tripRepository       repository.TripRepository
	unitOfWork           repository.UnitOfWork
	tripMemberRepository repository.TripMemberRepository
}

func NewTripService(
	tripRepository repository.TripRepository,
	unitOfWork repository.UnitOfWork,
	tripMemberRepository repository.TripMemberRepository,
) service.TripService {
	return &TripService{
		tripRepository:       tripRepository,
		unitOfWork:           unitOfWork,
		tripMemberRepository: tripMemberRepository,
	}
}

func (service *TripService) CreateTrip(ctx *gin.Context, tripRequest model.TripRequest, userId int64) (int64, string) {
	tx, err := service.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("TripService.CreateTrip - BeginTx Error: " + err.Error())
		return 0, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	defer service.unitOfWork.Rollback(tx)

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
		Status:                "not_started",
	}

	// create trip
	tripID, err := service.tripRepository.CreateCommand(ctx, trip, tx)
	if err != nil {
		log.Error("TripService.CreateTrip Error: " + err.Error())
		return 0, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// create trip member as administrator
	member := &entity.TripMember{
		TripID: tripID,
		UserID: userId,
		Role:   model.TripMemberRole.Administrator,
	}
	err = service.tripMemberRepository.CreateCommand(ctx, member, tx)
	if err != nil {
		log.Error("TripService.CreateTrip - CreateTripMember Error: " + err.Error())
		return 0, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// commit transaction
	err = service.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("TripService.CreateTrip - Commit Error: " + err.Error())
		return 0, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	return tripID, ""
}

func (service *TripService) GetAllTripsByUserID(ctx *gin.Context, userId int64) ([]*model.TripResponse, string) {
	trips, err := service.tripRepository.GetAllTripsWithUserRoleByUserIdQuery(ctx, userId, nil)
	if err != nil {
		log.Error("TripService.GetAllTripsByUserID Error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	var tripResponses []*model.TripResponse
	for _, trip := range trips {
		tripResponse := &model.TripResponse{
			Title:                 trip.Title,
			City:                  trip.City,
			StartDate:             trip.StartDate,
			Days:                  trip.Days,
			Budget:                trip.Budget,
			NumMembers:            trip.NumMembers,
			ViLocationAttributes:  trip.ViLocationAttributes,
			ViFoodAttributes:      trip.ViFoodAttributes,
			ViSpecialRequirements: trip.ViSpecialRequirements,
			ViMedicalConditions:   trip.ViMedicalConditions,
			EnLocationAttributes:  trip.EnLocationAttributes,
			EnFoodAttributes:      trip.EnFoodAttributes,
			EnSpecialRequirements: trip.EnSpecialRequirements,
			EnMedicalConditions:   trip.EnMedicalConditions,
			Status:                trip.Status,
			Role:                  trip.Role,
		}
		tripResponses = append(tripResponses, tripResponse)
	}

	return tripResponses, ""
}
