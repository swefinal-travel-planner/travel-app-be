package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
)

type TripMemberHandler struct {
	tripMemberService service.TripMemberService
}

func NewTripMemberHandler(tripMemberService service.TripMemberService) *TripMemberHandler {
	return &TripMemberHandler{
		tripMemberService: tripMemberService,
	}
}

// @Summary Get trip members
// @Description Get all members of a trip if the user is a member
// @Tags trips
// @Accept json
// @Produce json
// @Param tripId path int true "Trip ID"
// @Success 200 {array} model.TripMemberResponse
// @Failure 403 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
// @Router /trips/{tripId}/members [get]
func (h *TripMemberHandler) GetTripMembers(c *gin.Context) {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "tripId")
		c.JSON(statusCode, errResponse)
		return
	}

	userID := middleware.GetUserIdHelper(c)
	members, errCode := h.tripMemberService.GetTripMembersIfUserInTrip(c.Request.Context(), tripID, userID)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		c.JSON(statusCode, errResponse)
		return
	}

	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&members))
}
