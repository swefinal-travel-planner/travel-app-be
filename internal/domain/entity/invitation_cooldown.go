package entity

import "time"

type InvitationCooldown struct {
	ID                  int64      `json:"id,omitempty" db:"id"`
	UserID1             int64      `json:"userId1,omitempty" db:"user_id_1"`
	UserID2             int64      `json:"userId2,omitempty" db:"user_id_2"`
	StartCooldownMillis int64      `json:"startCooldownMillis,omitempty" db:"start_cooldown_millis"`
	CooldownDuration    int64      `json:"cooldownDuration,omitempty" db:"cooldown_duration"`
	CreatedAt           time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt           time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt           *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}
