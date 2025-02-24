package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/validation"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
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
func (handler *AuthHandler) Register(ctx *gin.Context) {
	var registerRequest model.RegisterRequest

	if err := validation.BindJsonAndValidate(ctx, &registerRequest); err != nil {
		return
	}

	err := handler.authService.Register(ctx, registerRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.AbortWithStatus(204)
}

// @Summary Login
// @Description Login to account
// @Tags Auths
// @Accept json
// @Param request body model.LoginRequest true "Auth payload"
// @Produce  json
// @Router /auth/login [post]
// @Success 200 {object} httpcommon.HttpResponse[entity.User]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) Login(ctx *gin.Context) {
	var loginRequest model.LoginRequest

	if err := validation.BindJsonAndValidate(ctx, &loginRequest); err != nil {
		return
	}

	user, err := handler.authService.Login(ctx, loginRequest)
	if err != nil || user == nil {
		if user == nil && err == nil {
			err = errors.New(httpcommon.ErrorMessage.SqlxNoRow)
		}
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
		return
	}

	ctx.JSON(200, httpcommon.NewSuccessResponse(user))
}

func (handler *AuthHandler) Test(ctx *gin.Context) {
	message := "hello world"
	ctx.JSON(200, httpcommon.NewSuccessResponse(&message))
}
