package model

type UpdateNotificationTokenRequest struct {
	NotificationToken string `json:"notificationToken" binding:"required"`
}

type UpdateUserRequest struct {
	Email             string  `json:"email,omitempty"`
	Name              string  `json:"name,omitempty"`
	PhoneNumber       string  `json:"phoneNumber,omitempty"`
	PhotoURL          *string `json:"photoURL,omitempty"`
	NotificationToken *string `json:"notificationToken,omitempty"`
}
