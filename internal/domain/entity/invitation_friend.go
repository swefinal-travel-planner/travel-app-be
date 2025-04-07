package entity

import "time"

type InvitationFriend struct {
	ID         int64      `json:"id,omitempty" db:"id"`
	SenderID   int64      `json:"senderId,omitempty" db:"sender_id"`
	ReceiverID int64      `json:"receiverId,omitempty" db:"receiver_id"`
	CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}
