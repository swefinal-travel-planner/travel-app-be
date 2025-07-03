package model

import "time"

type TestNotification struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	PushToken string `json:"pushToken"`
}

type GetAllNotificationFilters struct {
	Type string
}

type NotificationTriggerEntity struct {
	Type   string  `json:"type"`
	Avatar *string `json:"avatar"`
	Name   string  `json:"name"`
	ID     *int64  `json:"id,omitempty"`
}

type NotificationReferenceEntity struct {
	Type string `json:"type"`
	ID   *int64 `json:"id,omitempty"`
}

type NotificationResponse struct {
	ID              int64                       `json:"id"`
	Type            string                      `json:"type"`
	IsSeen          bool                        `json:"isSeen"`
	CreatedAt       time.Time                   `json:"createdAt"`
	TriggerEntity   NotificationTriggerEntity   `json:"triggerEntity"`
	ReferenceEntity NotificationReferenceEntity `json:"referenceEntity"`
	ReferenceData   interface{}                 `json:"referenceData,omitempty"`
}

type SaveNotificationRequest struct {
	ReceiverUserID      int64  `json:"receiverUserID"`
	TriggerEntityType   string `json:"triggerEntityType"`
	TriggerEntityID     *int64 `json:"triggerEntityID,omitempty"`
	ReferenceEntityType string `json:"referenceEntityType"`
	ReferenceEntityID   *int64 `json:"referenceEntityID,omitempty"`
	Type                string `json:"type"`
}
