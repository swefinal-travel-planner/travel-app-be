package model

import (
	"time"

	stringlistutils "github.com/swefinal-travel-planner/travel-app-be/internal/utils/string_list_utils"
)

type TripRequest struct {
	Title                 string                        `json:"title" binding:"required,min=1"`
	City                  string                        `json:"city" binding:"required"`
	StartDate             time.Time                     `json:"startDate" binding:"required"`
	Days                  int                           `json:"days" binding:"required,min=1"`
	Budget                float64                       `json:"budget"`
	NumMembers            int                           `json:"numMembers"`
	ViLocationAttributes  stringlistutils.SqlListString `json:"viLocationAttributes"`
	ViFoodAttributes      stringlistutils.SqlListString `json:"viFoodAttributes"`
	ViSpecialRequirements stringlistutils.SqlListString `json:"viSpecialRequirements"`
	ViMedicalConditions   stringlistutils.SqlListString `json:"viMedicalConditions"`
	EnLocationAttributes  stringlistutils.SqlListString `json:"enLocationAttributes"`
	EnFoodAttributes      stringlistutils.SqlListString `json:"enFoodAttributes"`
	EnSpecialRequirements stringlistutils.SqlListString `json:"enSpecialRequirements"`
	EnMedicalConditions   stringlistutils.SqlListString `json:"enMedicalConditions"`
}

type TripToCoreRequest struct {
	City                string                        `json:"city" binding:"required"`
	Days                int                           `json:"days" binding:"required,min=1"`
	LocationAttributes  stringlistutils.SqlListString `json:"locationAttributes"`
	FoodAttributes      stringlistutils.SqlListString `json:"foodAttributes"`
	SpecialRequirements stringlistutils.SqlListString `json:"specialRequirements"`
	MedicalConditions   stringlistutils.SqlListString `json:"medicalConditions"`
}

type TripResponse struct {
	Title                 string                        `json:"title" binding:"required,min=1"`
	City                  string                        `json:"city" binding:"required"`
	StartDate             time.Time                     `json:"startDate" binding:"required"`
	Days                  int                           `json:"days" binding:"required,min=1"`
	Budget                float64                       `json:"budget"`
	NumMembers            int                           `json:"numMembers"`
	ViLocationAttributes  stringlistutils.SqlListString `json:"viLocationAttributes"`
	ViFoodAttributes      stringlistutils.SqlListString `json:"viFoodAttributes"`
	ViSpecialRequirements stringlistutils.SqlListString `json:"viSpecialRequirements"`
	ViMedicalConditions   stringlistutils.SqlListString `json:"viMedicalConditions"`
	EnLocationAttributes  stringlistutils.SqlListString `json:"enLocationAttributes"`
	EnFoodAttributes      stringlistutils.SqlListString `json:"enFoodAttributes"`
	EnSpecialRequirements stringlistutils.SqlListString `json:"enSpecialRequirements"`
	EnMedicalConditions   stringlistutils.SqlListString `json:"enMedicalConditions"`
	Status                string                        `json:"status"`
}

type CreateTripResponse struct {
	ID int64 `json:"id"`
}
