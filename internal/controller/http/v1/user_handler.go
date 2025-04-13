package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// @Summary Search 1 friend
// @Description Search 1 friend by email
// @Tags User
// @Accept json
// @Produce  json
// @Router /users/{userEmail} [get]
// @Param userEmail query string true "Email of the friend"
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[model.FriendResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *UserHandler) SearchUser(ctx *gin.Context) {
	userEmail := ctx.Query("userEmail")
	if userEmail == "" {
		ctx.JSON(http.StatusBadRequest, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "userEmail is required", Field: "", Code: httpcommon.ErrorResponseCode.InvalidRequest,
		}))
	}
	user, err := handler.userService.SearchUser(ctx, userEmail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[model.FriendResponse](user))
}
