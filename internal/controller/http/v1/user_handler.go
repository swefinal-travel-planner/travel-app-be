package v1

import (
	"net/http"

	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/validation"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// @Summary Search 1 friend
// @Description Search 1 friend by email
// @Tags Users
// @Accept json
// @Produce  json
// @Router /users/{userEmail} [get]
// @Param userEmail query string true "Email of the friend"
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[model.FriendResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *UserHandler) SearchUser(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	userEmail := ctx.Query("userEmail")
	if userEmail == "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.BAD_REQUEST, "userEmail")
		ctx.JSON(statusCode, errResponse)
		return
	}
	user, errCode := handler.userService.SearchUser(ctx, userId, userEmail)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[model.FriendResponse](user))
}

// @Summary Update notification token
// @Description Update user's notification token
// @Tags Users
// @Accept json
// @Produce json
// @Router /users/notification-token [put]
// @Param request body model.UpdateNotificationTokenRequest true "Update notification token request"
// @Param Authorization header string true "Authorization: Bearer"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 401 {object} httpcommon.HttpResponse[any]
// @Failure 404 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *UserHandler) UpdateNotificationToken(ctx *gin.Context) {
	var request model.UpdateNotificationTokenRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	userId := middleware.GetUserIdHelper(ctx)
	errCode := handler.userService.UpdateNotificationToken(ctx, userId, request)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(204)
}
