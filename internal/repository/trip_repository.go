package repository

import (
	"context"

	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type TripRepository interface {
	CreateCommand(ctx context.Context, trip *entity.Trip) (int64, error)
}
