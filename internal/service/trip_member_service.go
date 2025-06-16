package service

import (
	"context"

	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/model"
)

type TripMemberService interface {
	GetTripMembersIfUserInTrip(ctx context.Context, tripID int64, userID int64) ([]model.TripMemberResponse, string)
}
