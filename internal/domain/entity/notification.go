package entity

import "time"

type Notification struct {
	ID                  int64      `json:"id" db:"id"`
	UserID              int64      `json:"userId" db:"user_id"`
	Type                string     `json:"type" db:"type"`
	IsSeen              bool       `json:"isSeen" db:"is_seen"`
	TriggerEntityType   string     `json:"triggerEntity.type" db:"trigger_entity_type"`
	TriggerEntityAvatar *string    `json:"triggerEntity.avatar" db:"trigger_entity_avatar"`
	TriggerEntityName   string     `json:"triggerEntity.name" db:"trigger_entity_name"`
	TriggerEntityID     *int64     `json:"triggerEntity.id" db:"trigger_entity_id"`
	ReferenceEntityType string     `json:"referenceEntity.type" db:"reference_entity_type"`
	ReferenceEntityID   *int64     `json:"referenceEntity.id" db:"reference_entity_id"`
	CreatedAt           time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt           *time.Time `json:"deletedAt" db:"deleted_at"`
}

type notificationType struct {
	FriendRequestReceived  string
	FriendRequestAccepted  string
	TripInvitationReceived string
	TripGenerated          string
}

var NotificationType = notificationType{
	FriendRequestReceived:  "friendRequestReceived",
	FriendRequestAccepted:  "friendRequestAccepted",
	TripInvitationReceived: "tripInvitationReceived",
	TripGenerated:          "tripGenerated",
}

type notificationReferenceType struct {
	FriendInvitation string
	TripInvitation   string
	TripGeneration   string
}

var NotificationReferenceType = notificationReferenceType{
	FriendInvitation: "friendInvitation",
	TripInvitation:   "tripInvitation",
	TripGeneration:   "tripGeneration",
}

type notificationTriggerType struct {
	User   string
	System string
}

var NotificationTriggerType = notificationTriggerType{
	User:   "user",
	System: "system",
}
