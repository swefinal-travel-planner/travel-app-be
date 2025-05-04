package entity

import (
	"database/sql"
	"time"

	stringlistutils "github.com/swefinal-travel-planner/travel-app-be/internal/utils/string_list_utils"
)

type Trip struct {
	ID                    int64                      `json:"id,omitempty" db:"id"`
	Title                 string                     `json:"title,omitempty" db:"title"`
	City                  string                     `json:"city,omitempty" db:"city"`
	StartDate             time.Time                  `json:"startDate,omitempty" db:"start_date"`
	Days                  int                        `json:"days,omitempty" db:"days"`
	Budget                float64                    `json:"budget,omitempty" db:"budget"`
	NumMembers            int                        `json:"numMembers,omitempty" db:"num_members"`
	ViLocationAttributes  stringlistutils.StringList `json:"viLocationAttributes,omitempty" db:"vi_location_attributes"`
	ViFoodAttributes      stringlistutils.StringList `json:"viFoodAttributes,omitempty" db:"vi_food_attributes"`
	ViSpecialRequirements stringlistutils.StringList `json:"viSpecialRequirements,omitempty" db:"vi_special_requirements"`
	ViMedicalConditions   stringlistutils.StringList `json:"viMedicalConditions,omitempty" db:"vi_medical_conditions"`
	EnLocationAttributes  stringlistutils.StringList `json:"enLocationAttributes,omitempty" db:"en_location_attributes"`
	EnFoodAttributes      stringlistutils.StringList `json:"enFoodAttributes,omitempty" db:"en_food_attributes"`
	EnSpecialRequirements stringlistutils.StringList `json:"enSpecialRequirements,omitempty" db:"en_special_requirements"`
	EnMedicalConditions   stringlistutils.StringList `json:"enMedicalConditions,omitempty" db:"en_medical_conditions"`
	Status                string                     `json:"status,omitempty" db:"status"`
	CreatedAt             time.Time                  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt             time.Time                  `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt             sql.NullTime               `json:"deletedAt,omitempty" db:"deleted_at"`
}
