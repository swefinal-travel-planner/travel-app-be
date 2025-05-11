package v1

import (
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/validation"
)

type NotificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// @Summary Test Notification Route
// @Description Test route to verify NotificationService functionality
// @Tags Notification
// @Accept json
// @Produce json
// @Param request body model.TestNotification true "Test Notification payload"
// @Router /notifications/test [post]
// @Success 204 "No Content"
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *NotificationHandler) TestNotification(ctx *gin.Context) {
	var testNotificationRequest model.TestNotification

	if err := validation.BindJsonAndValidate(ctx, &testNotificationRequest); err != nil {
		return
	}

	errCode := handler.notificationService.SendTestNotification(ctx, testNotificationRequest)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(204)
}

// @Summary Get All Notification Route
// @Description Get all notification for a user
// @Tags Notification
// @Accept json
// @Produce json
// @Param type query string false "Type"
// @Router /notifications [get]
// @Success 200 {object} httpcommon.HttpResponse[[]model.NotificationResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *NotificationHandler) GetAllNotification(ctx *gin.Context) {
	userID := ctx.GetInt64("userID")

	notifications, errCode := handler.notificationService.GetAllNotification(ctx, userID, model.GetAllNotificationFilters{
		Type: ctx.Query("type"),
	})

	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.JSON(200, notifications)
}

// @Summary Seen Notification Route
// @Description Mark a notification as seen
// @Tags Notification
// @Accept json
// @Produce json
// @Param id path int true "Notification ID"
// @Router /notifications/{id}/seen [post]
// @Success 204 "No Content"
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *NotificationHandler) SeenNotification(ctx *gin.Context) {
	userID := middleware.GetUserIdHelper(ctx)
	notificationID := ctx.Param("id")

	notificationIDInt, err := strconv.ParseInt(notificationID, 10, 64)
	if err != nil {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(err.Error(), "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	errCode := handler.notificationService.SeenNotification(ctx, userID, notificationIDInt)
	if errCode != "" {
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(errCode, "")
		ctx.JSON(statusCode, errResponse)
		return
	}

	ctx.AbortWithStatus(204)
}
