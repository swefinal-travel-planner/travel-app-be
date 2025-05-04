package entity

import (
	"database/sql"
	"time"
)

type TripMember struct {
	ID        int64        `json:"id,omitempty" db:"id"`
	TripID    int64        `json:"tripId,omitempty" db:"trip_id"`
	UserID    int64        `json:"userId,omitempty" db:"user_id"`
	Role      string       `json:"role,omitempty" db:"role"`
	CreatedAt time.Time    `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt time.Time    `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt sql.NullTime `json:"deletedAt,omitempty" db:"deleted_at"`
}
