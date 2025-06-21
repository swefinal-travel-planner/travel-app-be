package serviceimplement

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/env"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
)

type TripService struct {
	tripRepository       repository.TripRepository
	unitOfWork           repository.UnitOfWork
	tripMemberRepository repository.TripMemberRepository
	tripItemService      service.TripItemService
}

func NewTripService(
	tripRepository repository.TripRepository,
	unitOfWork repository.UnitOfWork,
	tripMemberRepository repository.TripMemberRepository,
	tripItemService service.TripItemService,
) service.TripService {
	return &TripService{
		tripRepository:       tripRepository,
		unitOfWork:           unitOfWork,
		tripMemberRepository: tripMemberRepository,
		tripItemService:      tripItemService,
	}
}

func (service *TripService) CreateTrip(ctx *gin.Context, tripRequest model.CreateTripManuallyRequest, userId int64) (int64, string) {
	tx, err := service.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("TripService.CreateTrip - BeginTx Error: " + err.Error())
		return 0, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	defer service.unitOfWork.Rollback(tx)

	if tripRequest.Status == "" {
		tripRequest.Status = model.TripStatus.NotStarted
	}
	trip := &entity.Trip{
		Title:                 tripRequest.Title,
		City:                  tripRequest.City,
		StartDate:             tripRequest.StartDate,
		Days:                  tripRequest.Days,
		ViLocationAttributes:  tripRequest.ViLocationAttributes,
		ViFoodAttributes:      tripRequest.ViFoodAttributes,
		ViSpecialRequirements: tripRequest.ViSpecialRequirements,
		ViMedicalConditions:   tripRequest.ViMedicalConditions,
		EnLocationAttributes:  tripRequest.EnLocationAttributes,
		EnFoodAttributes:      tripRequest.EnFoodAttributes,
		EnSpecialRequirements: tripRequest.EnSpecialRequirements,
		EnMedicalConditions:   tripRequest.EnMedicalConditions,
		Status:                tripRequest.Status,
		ReferenceID:           tripRequest.ReferenceID,
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
	trips, err := service.tripRepository.GetAllWithUserRoleByUserIdQuery(ctx, userId, nil)
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

func (service *TripService) GetTripByID(ctx *gin.Context, tripId int64, userId int64) (*model.TripResponse, string) {
	// First check if user is a member of the trip
	isMember, err := service.tripMemberRepository.IsUserInTripQuery(ctx, tripId, userId, nil)
	if err != nil {
		log.Error("TripService.GetTripByID - Check membership Error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if !isMember {
		return nil, error_utils.ErrorCode.FORBIDDEN
	}

	trip, err := service.tripRepository.GetOneWithUserRoleByIDQuery(ctx, tripId, userId, nil)
	if err != nil {
		log.Error("TripService.GetTripByID Error: " + err.Error())
		return nil, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if trip == nil {
		return nil, error_utils.ErrorCode.FORBIDDEN
	}

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

	return tripResponse, ""
}

func (service *TripService) updatedTripHelper(ctx *gin.Context, tripId int64, tripRequest model.TripPatchRequest, tx *sqlx.Tx) string {
	// Get existing trip
	existingTrip, err := service.tripRepository.GetOneByIDQuery(ctx, tripId, tx)
	if err != nil {
		log.Error("TripService.UpdateTrip - Get trip Error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if existingTrip == nil {
		return error_utils.ErrorCode.TRIP_NOT_FOUND
	}

	// Update fields if provided
	if tripRequest.Title != nil {
		existingTrip.Title = *tripRequest.Title
	}
	if tripRequest.City != nil {
		existingTrip.City = *tripRequest.City
	}
	if tripRequest.StartDate != nil {
		existingTrip.StartDate = *tripRequest.StartDate
	}
	if tripRequest.Days != nil {
		existingTrip.Days = *tripRequest.Days
	}
	if tripRequest.Budget != nil {
		existingTrip.Budget = *tripRequest.Budget
	}
	if tripRequest.NumMembers != nil {
		existingTrip.NumMembers = *tripRequest.NumMembers
	}
	if tripRequest.ViLocationAttributes != nil {
		existingTrip.ViLocationAttributes = *tripRequest.ViLocationAttributes
	}
	if tripRequest.ViFoodAttributes != nil {
		existingTrip.ViFoodAttributes = *tripRequest.ViFoodAttributes
	}
	if tripRequest.ViSpecialRequirements != nil {
		existingTrip.ViSpecialRequirements = *tripRequest.ViSpecialRequirements
	}
	if tripRequest.ViMedicalConditions != nil {
		existingTrip.ViMedicalConditions = *tripRequest.ViMedicalConditions
	}
	if tripRequest.EnLocationAttributes != nil {
		existingTrip.EnLocationAttributes = *tripRequest.EnLocationAttributes
	}
	if tripRequest.EnFoodAttributes != nil {
		existingTrip.EnFoodAttributes = *tripRequest.EnFoodAttributes
	}
	if tripRequest.EnSpecialRequirements != nil {
		existingTrip.EnSpecialRequirements = *tripRequest.EnSpecialRequirements
	}
	if tripRequest.EnMedicalConditions != nil {
		existingTrip.EnMedicalConditions = *tripRequest.EnMedicalConditions
	}
	if tripRequest.Status != nil {
		existingTrip.Status = *tripRequest.Status
	}
	if tripRequest.ReferenceID != nil {
		existingTrip.ReferenceID = *tripRequest.ReferenceID
	}

	// Update trip
	err = service.tripRepository.UpdateCommand(ctx, existingTrip, tx)
	if err != nil {
		log.Error("TripService.updatedTripHelper - Update trip Error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	return ""
}

func (service *TripService) UpdateTrip(ctx *gin.Context, tripId int64, userId int64, tripRequest model.TripPatchRequest) string {
	tx, err := service.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("TripService.UpdateTrip - BeginTx Error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	defer service.unitOfWork.Rollback(tx)

	// Check if user is admin
	isAdmin, err := service.tripMemberRepository.IsUserTripAdminQuery(ctx, tripId, userId, tx)
	if err != nil {
		log.Error("TripService.UpdateTrip - Check admin Error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if !isAdmin {
		return error_utils.ErrorCode.FORBIDDEN
	}

	errCode := service.updatedTripHelper(ctx, tripId, tripRequest, tx)
	if errCode != "" {
		return errCode
	}

	// Commit transaction
	err = service.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("TripService.UpdateTrip - Commit Error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	return ""
}

func (service *TripService) genToken(secretKey string, genTokenURL string) (string, string) {
	// call gen token URL to get token
	tokenReq := model.TokenRequest{
		SecretKey: secretKey,
	}
	tokenReqBody, marshalTokenReqBodyErr := json.Marshal(tokenReq)
	if marshalTokenReqBodyErr != nil {
		log.Error("TripService.genToken - Marshal token request Error: " + marshalTokenReqBodyErr.Error())
		return "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	resp, genTokenErr := http.Post(genTokenURL, "application/json", bytes.NewBuffer(tokenReqBody))
	if genTokenErr != nil {
		log.Error("TripService.genToken - Generate token Error: " + genTokenErr.Error())
		return "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("TripService.genToken - Generate token failed with status: " + resp.Status)
		return "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	var tokenResp model.TokenResponse
	decodeTokenRespErr := json.NewDecoder(resp.Body).Decode(&tokenResp)
	if decodeTokenRespErr != nil {
		log.Error("TripService.genToken - Decode token response Error: " + decodeTokenRespErr.Error())
		return "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	return tokenResp.Token, ""
}

func (service *TripService) createTripItems(createTourURL string, token string, tripToCoreRequest model.TripToCoreRequest) ([]model.TripItemFromAIResponse, string, string) {
	tripReqBody, marshalTripReqBodyErr := json.Marshal(tripToCoreRequest)
	if marshalTripReqBodyErr != nil {
		log.Error("TripService.createTripItemsRequest - Marshal trip request Error: " + marshalTripReqBodyErr.Error())
		return nil, "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	req, createRequestErr := http.NewRequest("POST", createTourURL, bytes.NewBuffer(tripReqBody))
	if createRequestErr != nil {
		log.Error("TripService.createTripItemsRequest - Create request Error: " + createRequestErr.Error())
		return nil, "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, sendRequestErr := client.Do(req)
	if sendRequestErr != nil {
		log.Error("TripService.createTripItemsRequest - Send request Error: " + sendRequestErr.Error())
		return nil, "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	body, readBodyErr := io.ReadAll(resp.Body)
	if readBodyErr != nil {
		log.Error("TripService.createTripItemsRequest - Read response body Error: " + readBodyErr.Error())
		return nil, "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("TripService.createTripItemsRequest - Create trip failed with status: " + resp.Status)
		return nil, "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	var tripItemsResp model.TripAIResponseWrapper
	if unmarshalTripItemErr := json.Unmarshal(body, &tripItemsResp); unmarshalTripItemErr != nil {
		log.Error("TripService.createTripItemsRequest - Unmarshal response Error: " + unmarshalTripItemErr.Error())
		return nil, "", error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	return tripItemsResp.Data.TripItems, tripItemsResp.Data.ReferenceID, ""
}

func (service *TripService) CreateTripByAI(ctx *gin.Context, tripRequest model.CreateTripByAIRequest, userID int64) ([]model.TripItemFromAIResponse, string) {
	tripToCoreRequest := model.TripToCoreRequest{
		City:                tripRequest.City,
		Days:                tripRequest.Days,
		LocationAttributes:  tripRequest.EnLocationAttributes,
		FoodAttributes:      tripRequest.EnFoodAttributes,
		SpecialRequirements: tripRequest.EnSpecialRequirements,
		MedicalConditions:   tripRequest.EnMedicalConditions,
		LocationsPerDay:     tripRequest.LocationsPerDay,
		LocationPreference:  tripRequest.LocationPreference,
	}

	// save trip to database
	createTripManuallyRequest := model.CreateTripManuallyRequest{
		Title:                 tripRequest.Title,
		City:                  tripRequest.City,
		StartDate:             tripRequest.StartDate,
		Days:                  tripRequest.Days,
		ViLocationAttributes:  tripRequest.ViLocationAttributes,
		ViFoodAttributes:      tripRequest.ViFoodAttributes,
		ViSpecialRequirements: tripRequest.ViSpecialRequirements,
		ViMedicalConditions:   tripRequest.ViMedicalConditions,
		EnLocationAttributes:  tripRequest.EnLocationAttributes,
		EnFoodAttributes:      tripRequest.EnFoodAttributes,
		EnSpecialRequirements: tripRequest.EnSpecialRequirements,
		EnMedicalConditions:   tripRequest.EnMedicalConditions,
		Status:                model.TripStatus.AIGenerating,
	}
	tripID, errCode := service.CreateTrip(ctx, createTripManuallyRequest, userID)
	if errCode != "" {
		return []model.TripItemFromAIResponse{}, errCode
	}
	var tripStatus string

	// get secret key & generate token URL
	secretKey, getSecretKeyErr := env.GetEnv("CORE_SECRET_KEY")
	if getSecretKeyErr != nil {
		tripStatus = model.TripStatus.Failed
		log.Error("TripService.CreateTripByAI - Get CORE_SECRET_KEY Error: " + getSecretKeyErr.Error())
		return []model.TripItemFromAIResponse{}, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	genTokenURL, getGenTokenURLErr := env.GetEnv("GEN_TOKEN_URL")
	if getGenTokenURLErr != nil {
		tripStatus = model.TripStatus.Failed
		log.Error("TripService.CreateTripByAI - Get GEN_TOKEN_URL Error: " + getGenTokenURLErr.Error())
		return []model.TripItemFromAIResponse{}, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// call gen token URL to get token
	token, genTokenErr := service.genToken(secretKey, genTokenURL)
	if genTokenErr != "" {
		tripStatus = model.TripStatus.Failed
		log.Error("TripService.CreateTripByAI - Generate token Error: " + genTokenErr)
		return []model.TripItemFromAIResponse{}, genTokenErr
	}

	// get create tour URL
	createTourURL, createTourURLErr := env.GetEnv("CREATE_TOUR_URL")
	if createTourURLErr != nil {
		tripStatus = model.TripStatus.Failed
		log.Error("TripService.CreateTripByAI - Get CREATE_TOUR_URL Error: " + createTourURLErr.Error())
		return []model.TripItemFromAIResponse{}, error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// send trip data to core service to get trip items
	tripItemsRespFromCore, referenceID, createTripItemsError := service.createTripItems(createTourURL, token, tripToCoreRequest)
	if createTripItemsError != "" {
		tripStatus = model.TripStatus.Failed
		log.Error("TripService.CreateTripByAI - Create trip items Error: " + createTripItemsError)
		return []model.TripItemFromAIResponse{}, createTripItemsError
	}

	// update reference ID and status
	defer func() {
		if tripID != 0 {
			tx, txErr := service.unitOfWork.Begin(ctx)
			if txErr != nil {
				log.Error("TripService.CreateTripByAI - BeginTx Error: " + txErr.Error())
				return
			}
			defer service.unitOfWork.Rollback(tx)

			tripStatus = model.TripStatus.NotStarted
			tripRequestForUpdate := model.TripPatchRequest{
				Status:      &tripStatus,
				ReferenceID: &referenceID,
			}

			updateErrCode := service.updatedTripHelper(ctx, tripID, tripRequestForUpdate, tx)
			if updateErrCode != "" {
				log.Error("TripService.CreateTripByAI - updateTrip Error: " + updateErrCode)
				return
			}

			commitErr := service.unitOfWork.Commit(tx)
			if commitErr != nil {
				log.Error("TripService.CreateTripByAI - Commit Error: " + commitErr.Error())
			}
		}
	}()

	// add tripID to trip items
	for i := range tripItemsRespFromCore {
		tripItemsRespFromCore[i].TripID = tripID
	}

	// convert trip items from core to request model
	var tripItemsReqs []model.TripItemRequest
	for _, item := range tripItemsRespFromCore {
		req := model.TripItemRequest{
			PlaceID:    item.PlaceID,
			TripDay:    item.TripDay,
			OrderInDay: item.OrderInDay,
			TimeInDate: item.TimeInDay,
		}
		tripItemsReqs = append(tripItemsReqs, req)
	}

	// save trip items to database
	err := service.tripItemService.CreateTripItems(ctx, userID, tripID, tripItemsReqs)
	if err != "" {
		log.Error("TripService.CreateTripByAI Error: " + err)
		return []model.TripItemFromAIResponse{}, err
	}

	// return trip items to noti
	return tripItemsRespFromCore, ""
}

func (service *TripService) DeleteTrip(ctx *gin.Context, tripId int64, userId int64) string {
	tx, err := service.unitOfWork.Begin(ctx)
	if err != nil {
		log.Error("TripService.DeleteTrip - BeginTx Error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	defer service.unitOfWork.Rollback(tx)

	// Check if user is admin
	isAdmin, err := service.tripMemberRepository.IsUserTripAdminQuery(ctx, tripId, userId, tx)
	if err != nil {
		log.Error("TripService.DeleteTrip - Check admin Error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}
	if !isAdmin {
		return error_utils.ErrorCode.FORBIDDEN
	}

	// Delete trip
	err = service.tripRepository.DeleteByIDCommand(ctx, tripId, tx)
	if err != nil {
		log.Error("TripService.DeleteTrip - Delete Error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	// Commit transaction
	err = service.unitOfWork.Commit(tx)
	if err != nil {
		log.Error("TripService.DeleteTrip - Commit Error: " + err.Error())
		return error_utils.ErrorCode.INTERNAL_SERVER_ERROR
	}

	return ""
}
