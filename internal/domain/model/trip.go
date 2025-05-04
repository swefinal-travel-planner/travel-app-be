package model

import (
	"time"

	stringlistutils "github.com/swefinal-travel-planner/travel-app-be/internal/utils/string_list_utils"
)

type TripRequest struct {
	Title                 string                     `json:"title" binding:"required,min=1"`
	City                  string                     `json:"city" binding:"required"`
	StartDate             time.Time                  `json:"startDate" binding:"required"`
	Days                  int                        `json:"days" binding:"required,min=1"`
	Budget                float64                    `json:"budget"`
	NumMembers            int                        `json:"numMembers"`
	ViLocationAttributes  stringlistutils.StringList `json:"viLocationAttributes"`
	ViFoodAttributes      stringlistutils.StringList `json:"viFoodAttributes"`
	ViSpecialRequirements stringlistutils.StringList `json:"viSpecialRequirements"`
	ViMedicalConditions   stringlistutils.StringList `json:"viMedicalConditions"`
	EnLocationAttributes  stringlistutils.StringList `json:"enLocationAttributes"`
	EnFoodAttributes      stringlistutils.StringList `json:"enFoodAttributes"`
	EnSpecialRequirements stringlistutils.StringList `json:"enSpecialRequirements"`
	EnMedicalConditions   stringlistutils.StringList `json:"enMedicalConditions"`
}

type TripToCoreRequest struct {
	City                string                     `json:"city" binding:"required"`
	Days                int                        `json:"days" binding:"required,min=1"`
	LocationAttributes  stringlistutils.StringList `json:"locationAttributes"`
	FoodAttributes      stringlistutils.StringList `json:"foodAttributes"`
	SpecialRequirements stringlistutils.StringList `json:"specialRequirements"`
	MedicalConditions   stringlistutils.StringList `json:"medicalConditions"`
}
