package entity

import "time"

type InvitationTrip struct {
	ID         int64      `json:"id,omitempty" db:"id"`
	TripID     int64      `json:"tripId,omitempty" db:"trip_id"`
	SenderID   int64      `json:"senderId,omitempty" db:"sender_id"`
	ReceiverID int64      `json:"receiverId,omitempty" db:"receiver_id"`
	Status     string     `json:"status,omitempty" db:"status"`
	CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}
