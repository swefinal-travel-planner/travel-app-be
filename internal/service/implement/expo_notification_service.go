package serviceimplement

import (
	"github.com/gin-gonic/gin"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"os"
)

type ExpoNotificationService struct {
	expoClient *expo.PushClient
}

// NewNotificationService creates a new instance of NotificationService
func NewExpoNotificationService() service.NotificationService {
	expoConfig := expo.ClientConfig{
		AccessToken: os.Getenv("NOTIFICATION_ACCESS_TOKEN"),
	}
	expoClient := expo.NewPushClient(&expoConfig)
	return &ExpoNotificationService{
		expoClient: expoClient,
	}
}

func (n *ExpoNotificationService) SendTestNotification(ctx *gin.Context, testNotificationRequest model.TestNotification) string {
	pushToken, err := expo.NewExponentPushToken(testNotificationRequest.PushToken)
	if err != nil {
		log.Error("ExpoNotificationService.SendTestNotification NewExponentPushToken err: ", err)
		return err.Error()
	}

	// Publish message
	response, err := n.expoClient.Publish(
		&expo.PushMessage{
			To:       []expo.ExponentPushToken{pushToken},
			Body:     testNotificationRequest.Body,
			Sound:    "default",
			Title:    testNotificationRequest.Title,
			Priority: expo.DefaultPriority,
		},
	)

	if err != nil {
		log.Error("ExpoNotificationService.SendTestNotification err: ", err)
		return err.Error()
	}

	// Validate responses
	if response.ValidateResponse() != nil {
		log.Error("ExpoNotificationService.SendTestNotification failed: ", err)
		return response.ValidateResponse().Error()
	}

	return ""
}
