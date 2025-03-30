package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
)

type FriendHandler struct {
	friendService service.FriendService
}

func NewFriendHandler(friendService service.FriendService) *FriendHandler {
	return &FriendHandler{friendService: friendService}
}

// @Summary View friends
// @Description View friends
// @Tags Friend
// @Accept json
// @Produce  json
// @Router /friends [get]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[[]model.FriendResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *FriendHandler) ViewFriends(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	friends, err := handler.friendService.GetAllFriends(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[[]model.FriendResponse](&friends))
}
