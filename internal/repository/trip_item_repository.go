package repository

import (
	"context"

	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type TripItemRepository interface {
	CreateCommand(ctx context.Context, tripItem *entity.TripItem) error
}
