package model

type UpdateNotificationTokenRequest struct {
	NotificationToken string `json:"notificationToken" binding:"required"`
}
