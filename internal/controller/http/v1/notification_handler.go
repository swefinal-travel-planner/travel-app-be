package v1

import (
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
