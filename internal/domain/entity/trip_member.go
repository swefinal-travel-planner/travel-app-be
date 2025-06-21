package entity

import (
	"time"
)

type TripMember struct {
	ID        int64      `db:"id"`
	TripID    int64      `db:"trip_id"`
	UserID    int64      `db:"user_id"`
	Role      string     `db:"role"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type TripMemberWithUser struct {
	TripMember
	Name     string  `db:"name"`
	PhotoURL *string `db:"photo_url"`
}
