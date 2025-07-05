package v1

import (
	"strconv"

	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/validation"

	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
)

type TripHandler struct {
	tripService         service.TripService
	tripItemService     service.TripItemService
	notificationService service.NotificationService
}

func NewTripHandler(tripService service.TripService, tripItemService service.TripItemService, notificationService service.NotificationService) *TripHandler {
	return &TripHandler{
		tripService:         tripService,
		tripItemService:     tripItemService,
		notificationService: notificationService,
	}
}

// @Summary Create trip manually
// @Description Create trip manually
// @Tags Trips
// @Accept json
// @Param request body model.CreateTripManuallyRequest true "Trip payload"
// @Param  Authorization header string true "Authorization: Bearer"
// @Produce  json
// @Router /trips [post]
// @Success 200 {object} httpcommon.HttpResponse[model.CreateTripResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *TripHandler) CreateTripManually(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	var tripRequest model.CreateTripManuallyRequest

	if err := validation.BindJsonAndValidate(ctx, &tripRequest); err != nil {
		return
	}

	tripID, errCode := handler.tripService.CreateTrip(ctx, tripRequest, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	response := model.CreateTripResponse{
		ID: tripID,
	}
	ctx.JSON(200, httpcommon.NewSuccessResponse(&response))
}

// @Summary Create/Update trip items
// @Description Create/Update trip items
// @Tags Trips
// @Accept json
// @Param request body []model.TripItemRequest true "TripItem payload"
// @Param tripId path int true "Trip ID"
// @Param  Authorization header string true "Authorization: Bearer"
// @Produce json
// @Router /trips/{tripId}/trip-items [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *TripHandler) CreateTripItems(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	tripId := ctx.Param("tripId")
	if tripId == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "friendId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	tripIdInt, err := strconv.ParseInt(tripId, 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "friendId")
		ctx.JSON(statusCode, errResponse)
		return
	}
	var tripItemRequests []model.TripItemRequest

	if err := validation.BindJsonAndValidate(ctx, &tripItemRequests); err != nil {
		return
	}

	errCode := handler.tripItemService.CreateTripItems(ctx, userId, tripIdInt, tripItemRequests)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(204)
}

// @Summary Get all trips for a user
// @Description Get all trips that the user is a member of
// @Tags Trips
// @Param  Authorization header string true "Authorization: Bearer"
// @Produce json
// @Router /trips [get]
// @Success 200 {object} httpcommon.HttpResponse[[]model.TripResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *TripHandler) GetAllTrips(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	trips, errCode := handler.tripService.GetAllTripsByUserID(ctx, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(200, httpcommon.NewSuccessResponse(&trips))
}

// @Summary Get trip by ID
// @Description Get a specific trip with user's role
// @Tags Trips
// @Param tripId path int true "Trip ID"
// @Param  Authorization header string true "Authorization: Bearer"
// @Produce json
// @Router /trips/{tripId} [get]
// @Success 200 {object} httpcommon.HttpResponse[model.TripResponse]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *TripHandler) GetTrip(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	tripId := ctx.Param("tripId")
	if tripId == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "tripId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	tripIdInt, err := strconv.ParseInt(tripId, 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "tripId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	trip, errCode := handler.tripService.GetTripByID(ctx, tripIdInt, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(200, httpcommon.NewSuccessResponse(trip))
}

// @Summary Get trip items
// @Description Get trip items by trip ID
// @Tags Trips
// @Param tripId path int true "Trip ID"
// @Param  Authorization header string true "Authorization: Bearer"
// @Param language query string false "Language for place info (vi or en)" Enums(vi,en) default(vi)
// @Produce json
// @Router /trips/{tripId}/trip-items [get]
// @Success 200 {object} httpcommon.HttpResponse[[]model.TripItemResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *TripHandler) GetTripItems(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	tripId := ctx.Param("tripId")
	if tripId == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "tripId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	tripIdInt, err := strconv.ParseInt(tripId, 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "tripId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	// Get language from query, only allow 'vi' and 'en', default to 'vi'
	lang := ctx.DefaultQuery("language", "vi")
	if lang != "vi" && lang != "en" {
		lang = "vi"
	}

	tripItems, errCode := handler.tripItemService.GetTripItemsByTripID(ctx, userId, tripIdInt)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(200, httpcommon.NewSuccessResponse(&tripItems))
}

// @Summary Update trip
// @Description Update a trip's details
// @Tags Trips
// @Accept json
// @Param tripId path int true "Trip ID"
// @Param request body model.TripPatchRequest true "Trip update payload"
// @Param  Authorization header string true "Authorization: Bearer"
// @Produce json
// @Router /trips/{tripId} [patch]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 403 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *TripHandler) UpdateTrip(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	tripId := ctx.Param("tripId")
	if tripId == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "tripId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	tripIdInt, err := strconv.ParseInt(tripId, 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "tripId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	var tripRequest model.TripPatchRequest
	if err := validation.BindJsonAndValidate(ctx, &tripRequest); err != nil {
		return
	}

	errCode := handler.tripService.UpdateTrip(ctx, tripIdInt, userId, tripRequest)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(204)
}

// @Summary Create trip by AI
// @Description Create trip by AI
// @Tags Trips
// @Accept json
// @Param request body model.CreateTripByAIRequest true "Trip payload"
// @Param  Authorization header string true "Authorization: Bearer"
// @Produce  json
// @Router /trips/ai [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *TripHandler) CreateTripByAI(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	var tripRequest model.CreateTripByAIRequest
	if err := validation.BindJsonAndValidate(ctx, &tripRequest); err != nil {
		return
	}

	ctx.AbortWithStatus(204)

	// process the remaining tasks in the background
	go func(tripRequest model.CreateTripByAIRequest, userId int64) {
		defer func() {
			recover()
		}()

		routineCtx := ctx.Copy()

		tripItemRespList, tripId, errCode := handler.tripService.CreateTripByAI(routineCtx, tripRequest, userId)
		if errCode != "" {
			handler.tripService.UpdateTrip(routineCtx, tripId, userId, model.TripPatchRequest{
				Status: &model.TripStatus.Failed,
			})
			// send notification to current user
			notiErr := handler.notificationService.SaveAndSendNotification(ctx, model.SaveNotificationRequest{
				Type:                entity.NotificationType.TripGeneratedFailed,
				ReceiverUserID:      userId,
				TriggerEntityType:   entity.NotificationTriggerType.System,
				TriggerEntityID:     nil,
				ReferenceEntityType: entity.NotificationReferenceType.Trip,
				ReferenceEntityID:   &tripItemRespList[0].TripID,
			})
			if notiErr != "" {
				return
			}
		} else {
			// send notification to current user
			notiErr := handler.notificationService.SaveAndSendNotification(ctx, model.SaveNotificationRequest{
				Type:                entity.NotificationType.TripGenerated,
				ReceiverUserID:      userId,
				TriggerEntityType:   entity.NotificationTriggerType.System,
				TriggerEntityID:     nil,
				ReferenceEntityType: entity.NotificationReferenceType.Trip,
				ReferenceEntityID:   &tripItemRespList[0].TripID,
			})
			if notiErr != "" {
				return
			}
		}
	}(tripRequest, userId)
}

// @Summary Delete trip
// @Description Delete a trip (admin only)
// @Tags Trips
// @Param tripId path int true "Trip ID"
// @Param  Authorization header string true "Authorization: Bearer"
// @Produce json
// @Router /trips/{tripId} [delete]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 403 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *TripHandler) DeleteTrip(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	tripId := ctx.Param("tripId")
	if tripId == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "tripId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	tripIdInt, err := strconv.ParseInt(tripId, 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "tripId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	errCode := handler.tripService.DeleteTrip(ctx, tripIdInt, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(204)
}
