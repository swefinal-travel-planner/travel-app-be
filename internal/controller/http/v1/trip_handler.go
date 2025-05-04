package v1

import (
	"strconv"

	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/validation"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
)

type TripHandler struct {
	tripService     service.TripService
	tripItemService service.TripItemService
}

func NewTripHandler(tripService service.TripService, tripItemService service.TripItemService) *TripHandler {
	return &TripHandler{
		tripService:     tripService,
		tripItemService: tripItemService,
	}
}

// @Summary Create trip manually
// @Description Create trip manually
// @Tags Trips
// @Accept json
// @Param request body model.TripRequest true "Trip payload"
// @Param  Authorization header string true "Authorization: Bearer"
// @Produce  json
// @Router /trips [post]
// @Success 200 {object} httpcommon.HttpResponse[int64]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *TripHandler) CreateTripManually(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	var tripRequest model.TripRequest

	if err := validation.BindJsonAndValidate(ctx, &tripRequest); err != nil {
		return
	}

	tripID, errCode := handler.tripService.CreateTrip(ctx, tripRequest, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(200, httpcommon.NewSuccessResponse(&tripID))
}

// @Summary Create/Update trip items
// @Description Create/Update trip items
// @Tags Trips
// @Accept json
// @Param request body []model.TripItemRequest true "TripItem payload"
// @Param tripId path int true "Trip ID"
// @Param  Authorization header string true "Authorization: Bearer"
// @Produce json
// @Router /trips/:tripId/trip-items [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *TripHandler) CreateTripItems(ctx *gin.Context) {
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

	for i := range tripItemRequests {
		tripItemRequests[i].TripID = tripIdInt
	}

	errCode := handler.tripItemService.CreateTripItems(ctx, tripItemRequests)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(204)
}
