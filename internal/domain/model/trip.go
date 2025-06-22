package model

import (
	"time"

	stringlistutils "github.com/swefinal-travel-planner/travel-app-be/internal/utils/string_list_utils"
)

type CreateTripManuallyRequest struct {
	Title                 string                        `json:"title" binding:"required,min=1"`
	City                  string                        `json:"city" binding:"required"`
	StartDate             time.Time                     `json:"startDate" binding:"required"`
	Days                  int                           `json:"days" binding:"required,min=1,max=7"`
	ViLocationAttributes  stringlistutils.SqlListString `json:"-"`
	ViFoodAttributes      stringlistutils.SqlListString `json:"-"`
	ViSpecialRequirements stringlistutils.SqlListString `json:"-"`
	ViMedicalConditions   stringlistutils.SqlListString `json:"-"`
	EnLocationAttributes  stringlistutils.SqlListString `json:"-"`
	EnFoodAttributes      stringlistutils.SqlListString `json:"-"`
	EnSpecialRequirements stringlistutils.SqlListString `json:"-"`
	EnMedicalConditions   stringlistutils.SqlListString `json:"-"`
	Status                string                        `json:"-"`
	ReferenceID           *string                       `json:"-"`
}

type CreateTripByAIRequest struct {
	Title                 string                        `json:"title" binding:"required,min=1"`
	City                  string                        `json:"city" binding:"required"`
	StartDate             time.Time                     `json:"startDate" binding:"required"`
	Days                  int                           `json:"days" binding:"required,min=1,max=7"`
	ViLocationAttributes  stringlistutils.SqlListString `json:"viLocationAttributes"`
	ViFoodAttributes      stringlistutils.SqlListString `json:"viFoodAttributes"`
	ViSpecialRequirements stringlistutils.SqlListString `json:"viSpecialRequirements"`
	ViMedicalConditions   stringlistutils.SqlListString `json:"viMedicalConditions"`
	EnLocationAttributes  stringlistutils.SqlListString `json:"enLocationAttributes"`
	EnFoodAttributes      stringlistutils.SqlListString `json:"enFoodAttributes"`
	EnSpecialRequirements stringlistutils.SqlListString `json:"enSpecialRequirements"`
	EnMedicalConditions   stringlistutils.SqlListString `json:"enMedicalConditions"`
	LocationsPerDay       int                           `json:"locationsPerDay" binding:"required,min=5,max=9"`
	LocationPreference    string                        `json:"locationPreference"`
	ReferenceID           string                        `json:"referenceId,omitempty"`
}

type TripToCoreRequest struct {
	City                string                        `json:"city" binding:"required"`
	Days                int                           `json:"days" binding:"required,min=1"`
	LocationAttributes  stringlistutils.SqlListString `json:"location_attributes"`
	FoodAttributes      stringlistutils.SqlListString `json:"food_attributes"`
	SpecialRequirements stringlistutils.SqlListString `json:"special_requirements"`
	MedicalConditions   stringlistutils.SqlListString `json:"medical_conditions"`
	LocationsPerDay     int                           `json:"locationsPerDay" binding:"required,min=1"`
	LocationPreference  string                        `json:"locationPreference"`
}

type TripResponse struct {
	ID                    int64                         `json:"id"`
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
	ReferenceID           *string                        `json:"referenceId,omitempty"`
}

type TokenRequest struct {
	SecretKey string `json:"secret_key"`
}
type TokenResponse struct {
	Token string `json:"token"`
}

type TripAIResponseWrapper struct {
	Data struct {
		ReferenceID string                   `json:"reference_id"`
		TripItems   []TripItemFromAIResponse `json:"trip_items"`
	} `json:"data"`
	Status int `json:"status"`
}

type tripStatus struct {
	NotStarted   string
	InProgress   string
	Completed    string
	Received     string
	AIGenerating string
	Failed       string
}

var TripStatus = tripStatus{
	NotStarted:   "not_started",
	InProgress:   "in_progress",
	Completed:    "completed",
	Received:     "cancel",
	AIGenerating: "ai_generating",
	Failed:       "failed",
}
