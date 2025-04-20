package entity

import (
	"database/sql"
	"time"
)

type Trip struct {
	ID                    int64          `json:"id,omitempty" db:"id"`
	Title                 string         `json:"title,omitempty" db:"title"`
	City                  string         `json:"city,omitempty" db:"city"`
	StartDate             time.Time      `json:"startDate,omitempty" db:"start_date"`
	Days                  int            `json:"days,omitempty" db:"days"`
	Budget                float64        `json:"budget,omitempty" db:"budget"`
	NumMembers            int            `json:"numMembers,omitempty" db:"num_members"`
	ViLocationAttributes  sql.NullString `json:"viLocationAttributes,omitempty" db:"vi_location_attributes"`
	ViFoodAttributes      sql.NullString `json:"viFoodAttributes,omitempty" db:"vi_food_attributes"`
	ViSpecialRequirements sql.NullString `json:"viSpecialRequirements,omitempty" db:"vi_special_requirements"`
	ViMedicalConditions   sql.NullString `json:"viMedicalConditions,omitempty" db:"vi_medical_conditions"`
	EnLocationAttributes  sql.NullString `json:"enLocationAttributes,omitempty" db:"en_location_attributes"`
	EnFoodAttributes      sql.NullString `json:"enFoodAttributes,omitempty" db:"en_food_attributes"`
	EnSpecialRequirements sql.NullString `json:"enSpecialRequirements,omitempty" db:"en_special_requirements"`
	EnMedicalConditions   sql.NullString `json:"enMedicalConditions,omitempty" db:"en_medical_conditions"`
	Status                string         `json:"status,omitempty" db:"status"`
	CreatedAt             time.Time      `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt             time.Time      `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt             sql.NullTime   `json:"deletedAt,omitempty" db:"deleted_at"`
}
