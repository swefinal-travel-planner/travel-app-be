package serviceimplement

import (
	"os"

	"github.com/gin-gonic/gin"
	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
)

type ExpoNotificationService struct {
	expoClient             *expo.PushClient
	notificationRepository repository.NotificationRepository
	userRepository         repository.UserRepository
}

// NewNotificationService creates a new instance of NotificationService
func NewExpoNotificationService(notificationRepository repository.NotificationRepository, userRepository repository.UserRepository) service.NotificationService {
	expoConfig := expo.ClientConfig{
		AccessToken: os.Getenv("NOTIFICATION_ACCESS_TOKEN"),
	}
	expoClient := expo.NewPushClient(&expoConfig)
	return &ExpoNotificationService{
		expoClient:             expoClient,
		notificationRepository: notificationRepository,
		userRepository:         userRepository,
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

func (n *ExpoNotificationService) GeneratePushNotification(pushToken expo.ExponentPushToken, notification entity.Notification) *expo.PushMessage {
	var title string
	var body string

	switch notification.Type {
	case entity.NotificationType.FriendRequestReceived:
		title = "You have a new friend invitation"
		body = notification.TriggerEntityName + " invited you to be friends"
	case entity.NotificationType.FriendRequestAccepted:
		title = "Your friend request has been accepted"
		body = notification.TriggerEntityName + " accepted your friend request"
	}

	return &expo.PushMessage{
		To:    []expo.ExponentPushToken{pushToken},
		Body:  body,
		Sound: "default",
		Title: title,
	}
}

func (n *ExpoNotificationService) SendNotification(ctx *gin.Context, notification entity.Notification) {
	notificationToken, err := n.userRepository.GetNotificationTokenByIDQuery(ctx, notification.UserID, nil)
	if err != nil {
		log.Error("ExpoNotificationService.SendNotification err: ", err)
	}

	if notificationToken != nil {
		pushToken, err := expo.NewExponentPushToken(*notificationToken)
		if err != nil {
			log.Error("ExpoNotificationService.SendNotification NewExponentPushToken err: ", err)
		}

		pushNotification := n.GeneratePushNotification(pushToken, notification)
		n.expoClient.Publish(pushNotification)
	}
}

func (n *ExpoNotificationService) SaveAndSendNotification(ctx *gin.Context, notification model.SaveNotificationRequest) string {
	var notificationEntity entity.Notification
	if notification.TriggerEntityType == entity.NotificationTriggerType.User {
		user, err := n.userRepository.GetOneByIDQuery(ctx, *notification.TriggerEntityID, nil)
		if err != nil {
			log.Error("ExpoNotificationService.SaveAndSendNotification err: ", err)
			return err.Error()
		}

		notificationEntity.UserID = notification.ReceiverUserID
		notificationEntity.TriggerEntityAvatar = user.PhotoURL
		notificationEntity.TriggerEntityName = user.Name
		notificationEntity.TriggerEntityID = notification.TriggerEntityID
		notificationEntity.TriggerEntityType = entity.NotificationTriggerType.User
		notificationEntity.ReferenceEntityType = notification.ReferenceEntityType
		notificationEntity.ReferenceEntityID = notification.ReferenceEntityID
		notificationEntity.Type = notification.Type
	}

	err := n.notificationRepository.CreateCommand(ctx, &notificationEntity, nil)
	if err != nil {
		log.Error("ExpoNotificationService.SaveAndSendNotification err: ", err)
		return err.Error()
	}
	n.SendNotification(ctx, notificationEntity)
	return ""
}

func (n *ExpoNotificationService) GetAllNotification(ctx *gin.Context, userID int64, filters model.GetAllNotificationFilters) ([]model.NotificationResponse, string) {
	notifications, err := n.notificationRepository.GetAllByUserIDQuery(ctx, userID, filters.Type, nil)
	if err != nil {
		log.Error("ExpoNotificationService.GetAllNotification err: ", err)
		return nil, err.Error()
	}

	var notificationsResponse []model.NotificationResponse

	for _, notification := range notifications {
		notificationsResponse = append(notificationsResponse, model.NotificationResponse{
			Type:      notification.Type,
			IsSeen:    notification.IsSeen,
			CreatedAt: notification.CreatedAt,
			TriggerEntity: model.NotificationTriggerEntity{
				Type:   notification.TriggerEntityType,
				Avatar: notification.TriggerEntityAvatar,
				Name:   notification.TriggerEntityName,
				ID:     notification.TriggerEntityID,
			},
			ReferenceEntity: model.NotificationReferenceEntity{
				Type: notification.ReferenceEntityType,
				ID:   notification.ReferenceEntityID,
			},
		})
	}

	return notificationsResponse, ""
}

func (n *ExpoNotificationService) SeenNotification(ctx *gin.Context, userID int64, notificationID int64) string {
	err := n.notificationRepository.SeenNotificationCommand(ctx, userID, notificationID, nil)
	if err != nil {
		log.Error("ExpoNotificationService.SeenNotification err: ", err)
		return err.Error()
	}

	return ""
}

func (n *ExpoNotificationService) DeleteFriendInvitation(ctx *gin.Context, userId int64, typeFilter string, triggerEntityID int64) string {
	notification, err := n.notificationRepository.GetOneByUserIdAndTypeAndTriggerEntityIDQuery(ctx, userId, typeFilter, triggerEntityID, nil)
	if err != nil {
		log.Error("ExpoNotificationService.DeleteFriendInvitation err: ", err)
		return err.Error()
	}

	err = n.notificationRepository.DeleteNotificationCommand(ctx, notification.ID, nil)
	if err != nil {
		log.Error("ExpoNotificationService.DeleteFriendInvitation err: ", err)
		return err.Error()
	}

	return ""
}

func (n *ExpoNotificationService) DeleteTripInvitation(ctx *gin.Context, userId int64, tripId int64, triggerEntityID int64) string {
	err := n.notificationRepository.DeleteTripNotificationCommand(ctx, userId, triggerEntityID, tripId, nil)
	if err != nil {
		log.Error("ExpoNotificationService.DeleteFriendInvitation err: ", err)
		return err.Error()
	}

	return ""
}
