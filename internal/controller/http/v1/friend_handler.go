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

	friends, errCode := handler.friendService.GetAllFriends(ctx, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[[]model.FriendResponse](&friends))
}

// @Summary Remove friend
// @Description Remove friend
// @Tags Friend
// @Accept json
// @Param friendId path int true "friend ID"
// @Produce  json
// @Router /friends/{friendId} [put]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *FriendHandler) RemoveFriend(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	friendId := ctx.Param("friendId")
	if friendId == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "friendId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	friendIdInt, err := strconv.ParseInt(friendId, 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "friendId")
		ctx.JSON(statusCode, errResponse)
		return
	}

	errCode := handler.friendService.RemoveFriend(ctx, userId, friendIdInt)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}
	ctx.AbortWithStatus(204)
}
