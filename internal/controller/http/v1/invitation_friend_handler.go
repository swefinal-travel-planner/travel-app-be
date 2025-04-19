package v1

import (
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/validation"
)

type InvitationFriendHandler struct {
	invitationFriendService service.InvitationFriendService
}

func NewInvitationFriendHandler(invitationFriendService service.InvitationFriendService) *InvitationFriendHandler {
	return &InvitationFriendHandler{invitationFriendService: invitationFriendService}
}

// @Summary Add friend
// @Description Add friend
// @Tags InvitationFriend
// @Accept json
// @Param request body model.InvitationFriendRequest true "InvitationFriend payload"
// @Produce  json
// @Router /invitation-friends [post]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *InvitationFriendHandler) AddFriend(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	var invitationFriendRequest model.InvitationFriendRequest

	if err := validation.BindJsonAndValidate(ctx, &invitationFriendRequest); err != nil {
		return
	}

	errCode := handler.invitationFriendService.AddFriend(ctx, invitationFriendRequest, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}
	ctx.AbortWithStatus(204)
}

// @Summary Get all received invitations
// @Description Get all received invitations of current user
// @Tags InvitationFriend
// @Accept json
// @Produce  json
// @Router /invitation-friends/received [get]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[[]model.InvitationFriendReceivedResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *InvitationFriendHandler) GetAllReceivedInvitations(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	invitationFriends, errCode := handler.invitationFriendService.GetAllReceivedInvitations(ctx, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[[]model.InvitationFriendReceivedResponse](&invitationFriends))
}

// @Summary Get all requested invitations
// @Description Get all requested invitations of current user
// @Tags InvitationFriend
// @Accept json
// @Produce  json
// @Router /invitation-friends/requested [get]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[[]model.InvitationFriendRequestedResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *InvitationFriendHandler) GetAllRequestedInvitations(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	invitationFriends, errCode := handler.invitationFriendService.GetAllRequestedInvitations(ctx, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[[]model.InvitationFriendRequestedResponse](&invitationFriends))
}

// @Summary Accept friend invitation
// @Description Accept friend invitation
// @Tags InvitationFriend
// @Accept json
// @Param invitationId path int true "Invitation ID"
// @Produce  json
// @Router /invitation-friends/accept/{invitationId} [put]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *InvitationFriendHandler) AcceptInvitation(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	invitationId := ctx.Param("invitationId")
	if invitationId == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "invitationId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	invitationIdInt, err := strconv.ParseInt(invitationId, 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "invitationId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	errCode := handler.invitationFriendService.AcceptInvitation(ctx, invitationIdInt, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}
	ctx.AbortWithStatus(204)
}

// @Summary Deny friend invitation
// @Description Deny friend invitation
// @Tags InvitationFriend
// @Accept json
// @Param invitationId path int true "Invitation ID"
// @Produce  json
// @Router /invitation-friends/deny/{invitationId} [put]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *InvitationFriendHandler) DenyInvitation(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	invitationId := ctx.Param("invitationId")
	if invitationId == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "invitationId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	invitationIdInt, err := strconv.ParseInt(invitationId, 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "invitationId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	errCode := handler.invitationFriendService.DenyInvitation(ctx, invitationIdInt, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}
	ctx.AbortWithStatus(204)
}

// @Summary Withdraw friend invitation
// @Description Withdraw friend invitation (only allowed for the sender)
// @Tags InvitationFriend
// @Accept json
// @Param invitationId path int true "Invitation ID"
// @Produce  json
// @Router /invitation-friends/{invitationId} [delete]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *InvitationFriendHandler) WithdrawInvitation(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	invitationId := ctx.Param("invitationId")
	if invitationId == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "invitationId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	invitationIdInt, err := strconv.ParseInt(invitationId, 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "invitationId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	errCode := handler.invitationFriendService.WithdrawInvitation(ctx, invitationIdInt, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}
	ctx.AbortWithStatus(204)
}
