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

// @Summary Search users
// @Description Search users by email search term
// @Tags Users
// @Accept json
// @Produce  json
// @Router /users [get]
// @Query searchTerm query string true "Search term"
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[model.FriendResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *UserHandler) SearchUser(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)

	searchTerm := ctx.Query("searchTerm")

	users, errCode := handler.userService.SearchUser(ctx, userId, searchTerm)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[[]model.FriendResponse](&users))
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

// @Summary Update user profile
// @Description Update current user's profile information
// @Tags Users
// @Accept json
// @Produce json
// @Router /users/me [patch]
// @Param request body model.UpdateUserRequest true "Update user profile request"
// @Param Authorization header string true "Authorization: Bearer"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 401 {object} httpcommon.HttpResponse[any]
// @Failure 403 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *UserHandler) UpdateProfile(ctx *gin.Context) {
	var request model.UpdateUserRequest
	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	userId := middleware.GetUserIdHelper(ctx)
	errCode := handler.userService.UpdateUser(ctx, userId, request)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(204)
}

// @Summary Get user info
// @Description Get current user's info
// @Tags Users
// @Accept json
// @Produce json
// @Router /whoami [get]
// @Param Authorization header string true "Authorization: Bearer"
func (handler *UserHandler) WhoAmI(ctx *gin.Context) {
	userId := middleware.GetUserIdHelper(ctx)
	user, errCode := handler.userService.GetUserInfo(ctx, userId)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[model.UserInfoResponse](user))
}
