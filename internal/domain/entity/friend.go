package entity

import "time"

type Friend struct {
	ID        int64      `json:"id,omitempty" db:"id"`
	UserID1   int64      `json:"userId1,omitempty" db:"user_id_1"`
	UserID2   int64      `json:"userId2,omitempty" db:"user_id_2"`
	CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}
