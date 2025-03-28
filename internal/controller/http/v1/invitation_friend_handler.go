package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

// @Summary Register user
// @Description User register
// @Tags User
// @Accept json
// @Param request body model.RegisterRequest true "Auth payload"
// @Produce  json
// @Router /auth/register [post]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *InvitationFriendHandler) AddFriend(ctx *gin.Context) {
	var invitationFriendRequest model.InvitationFriendRequest

	if err := validation.BindJsonAndValidate(ctx, &invitationFriendRequest); err != nil {
		return
	}

	err := handler.invitationFriendService.AddFriend(ctx, invitationFriendRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.AbortWithStatus(204)
}

// @Summary Register user
// @Description User register
// @Tags User
// @Accept json
// @Param request body model.RegisterRequest true "Auth payload"
// @Produce  json
// @Router /auth/register [post]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *InvitationFriendHandler) GetAllInvitations(ctx *gin.Context) {
	invitationFriends, err := handler.invitationFriendService.GetAllInvitations(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[[]model.InvitationFriendResponse](&invitationFriends))
}
