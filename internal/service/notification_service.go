package service

import (
	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type NotificationService interface {
	SendTestNotification(ctx *gin.Context, testNotificationRequest model.TestNotification) string
}
