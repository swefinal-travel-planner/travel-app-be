package v1

import (
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/validation"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
)

type TripHandler struct {
	tripService service.TripService
}

func NewTripHandler(tripService service.TripService) *TripHandler {
	return &TripHandler{tripService: tripService}
}

func (handler *TripHandler) CreateTripByAI(ctx *gin.Context) {
	var tripRequest model.TripRequest

	if err := validation.BindJsonAndValidate(ctx, &tripRequest); err != nil {
		return
	}

	tripItemRespList, errCode := handler.tripService.CreateTripByAI(ctx, tripRequest)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(200, httpcommon.NewSuccessResponse(&tripItemRespList))
}
