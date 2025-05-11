package service

import (
	"github.com/gin-gonic/gin"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type NotificationService interface {
	SendTestNotification(ctx *gin.Context, testNotificationRequest model.TestNotification) string
	SendNotification(ctx *gin.Context, notification entity.Notification)
	SaveAndSendNotification(ctx *gin.Context, notification model.SaveNotificationRequest) string
	GeneratePushNotification(pushToken expo.ExponentPushToken, notification entity.Notification) *expo.PushMessage
	GetAllNotification(ctx *gin.Context, userID int64, filters model.GetAllNotificationFilters) ([]model.NotificationResponse, string)
	SeenNotification(ctx *gin.Context, userID int64, notificationID int64) string
	DeleteNotificationWithTypeAndTriggerEntityID(ctx *gin.Context, typeFilter string, triggerEntityID int64) string
}
