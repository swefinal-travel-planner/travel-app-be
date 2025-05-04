package entity

import (
	"database/sql"
	"time"
)

type TripItem struct {
	ID         int64        `json:"id,omitempty" db:"id"`
	TripID     int64        `json:"tripId,omitempty" db:"trip_id"`
	PlaceID    string       `json:"placeId,omitempty" db:"place_id"`
	TripDay    int64        `json:"tripDay,omitempty" db:"trip_day"`
	OrderInDay int64        `json:"orderInDay,omitempty" db:"order_in_day"`
	TimeInDate string       `json:"timeInDate,omitempty" db:"time_in_date"`
	CreatedAt  time.Time    `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt  time.Time    `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt  sql.NullTime `json:"deletedAt,omitempty" db:"deleted_at"`
}
