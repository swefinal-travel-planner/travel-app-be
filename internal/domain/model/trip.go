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
	Role                  string                        `json:"role"`
}

type CreateTripResponse struct {
	ID int64 `json:"id"`
}

type TripPatchRequest struct {
	Title                 *string                        `json:"title,omitempty"`
	City                  *string                        `json:"city,omitempty"`
	StartDate             *time.Time                     `json:"startDate,omitempty"`
	Days                  *int                           `json:"days,omitempty"`
	Budget                *float64                       `json:"budget,omitempty"`
	NumMembers            *int                           `json:"numMembers,omitempty"`
	ViLocationAttributes  *stringlistutils.SqlListString `json:"viLocationAttributes,omitempty"`
	ViFoodAttributes      *stringlistutils.SqlListString `json:"viFoodAttributes,omitempty"`
	ViSpecialRequirements *stringlistutils.SqlListString `json:"viSpecialRequirements,omitempty"`
	ViMedicalConditions   *stringlistutils.SqlListString `json:"viMedicalConditions,omitempty"`
	EnLocationAttributes  *stringlistutils.SqlListString `json:"enLocationAttributes,omitempty"`
	EnFoodAttributes      *stringlistutils.SqlListString `json:"enFoodAttributes,omitempty"`
	EnSpecialRequirements *stringlistutils.SqlListString `json:"enSpecialRequirements,omitempty"`
	EnMedicalConditions   *stringlistutils.SqlListString `json:"enMedicalConditions,omitempty"`
	Status                *string                        `json:"status,omitempty"`
}
