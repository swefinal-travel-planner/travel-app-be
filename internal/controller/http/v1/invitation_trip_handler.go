package v1

import (
	"net/http"
	"strconv"

	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/validation"

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

// @Summary Accept trip invitation
// @Description Accept an invitation to join a trip
// @Tags invitation-trips
// @Accept json
// @Produce json
// @Param invitationId path int true "Invitation ID"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/trip-invitations/accept/{invitationId} [put]
func (h *InvitationTripHandler) AcceptInvitation(ctx *gin.Context) {
	invitationId, err := strconv.ParseInt(ctx.Param("invitationId"), 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	userId := middleware.GetUserIdHelper(ctx)
	errCode := h.invitationTripService.AcceptInvitation(ctx, invitationId, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(http.StatusNoContent)
}

// @Summary Deny trip invitation
// @Description Deny an invitation to join a trip
// @Tags invitation-trips
// @Accept json
// @Produce json
// @Param invitationId path int true "Invitation ID"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/trip-invitations/deny/{invitationId} [put]
func (h *InvitationTripHandler) DenyInvitation(ctx *gin.Context) {
	invitationId, err := strconv.ParseInt(ctx.Param("invitationId"), 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	userId := middleware.GetUserIdHelper(ctx)
	errCode := h.invitationTripService.DenyInvitation(ctx, invitationId, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(http.StatusNoContent)
}

// @Summary Withdraw trip invitation (DISABLED)
// @Description Withdraw a sent trip invitation
// @Tags invitation-trips
// @Accept json
// @Produce json
// @Param invitationId path int true "Invitation ID"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/v1/trip-invitations/withdraw/{invitationId} [delete]
func (h *InvitationTripHandler) WithdrawInvitation(ctx *gin.Context) {
	invitationId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	userId := middleware.GetUserIdHelper(ctx)
	errCode := h.invitationTripService.WithdrawInvitation(ctx, invitationId, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(http.StatusNoContent)
}
