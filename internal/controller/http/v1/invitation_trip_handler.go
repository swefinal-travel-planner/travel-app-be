package v1

import (
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
)

type InvitationTripHandler struct {
	invitationTripService service.InvitationTripService
}

func NewInvitationTripHandler(invitationTripService service.InvitationTripService) *InvitationTripHandler {
	return &InvitationTripHandler{
		invitationTripService: invitationTripService,
	}
}

// @Summary Send trip invitation
// @Description Send an invitation to join a trip
// @Tags invitation-trips
// @Accept json
// @Produce json
// @Param request body model.InvitationTripRequest true "Send Invitation Request"
// @Security BearerAuth
// @Success 200 {object} model.InvitationTripSentResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/trip-invitations [post]
func (h *InvitationTripHandler) SendInvitation(ctx *gin.Context) {
	var request model.InvitationTripRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	userId := middleware.GetUserIdHelper(ctx)
	errCode := h.invitationTripService.SendInvitation(ctx, request, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(http.StatusNoContent)
}
