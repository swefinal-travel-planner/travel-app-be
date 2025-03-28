package v1

import (
	"net/http"

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

	err := handler.invitationFriendService.AddFriend(ctx, invitationFriendRequest, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.AbortWithStatus(204)
}

// @Summary Get all invitations
// @Description Get all invitations of current user
// @Tags InvitationFriend
// @Accept json
// @Produce  json
// @Router /invitation-friends [get]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[[]model.InvitationFriendResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *InvitationFriendHandler) GetAllInvitations(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	invitationFriends, err := handler.invitationFriendService.GetAllInvitations(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[[]model.InvitationFriendResponse](&invitationFriends))
}
