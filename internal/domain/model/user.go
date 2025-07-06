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

type UserInfoResponse struct {
	ID                int64   `json:"id"`
	Email             string  `json:"email"`
	Name              string  `json:"name"`
	PhotoURL          *string `json:"photoURL"`
	PhoneNumber       string  `json:"phoneNumber"`
	IDToken           *string `json:"idToken"`
	NotificationToken *string `json:"notificationToken"`
}
