package entity

import "time"

type User struct {
	Id                int64      `db:"id" json:"id,omitempty"`
	Email             string     `db:"email" json:"email,omitempty"`
	Name              string     `db:"name" json:"name,omitempty"`
	PhoneNumber       string     `db:"phone_number" json:"phoneNumber,omitempty"`
	Password          string     `db:"password" json:"password,omitempty"`
	PhotoURL          *string    `db:"photo_url" json:"photoURL,omitempty"`
	IDToken           *string    `db:"id_token" json:"idToken,omitempty"`
	NotificationToken *string    `db:"notification_token" json:"notificationToken,omitempty"`
	CreatedAt         *time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt         *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
	DeletedAt         *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}
