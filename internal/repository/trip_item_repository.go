package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type TripItemRepository interface {
	CreateCommand(ctx context.Context, tripItem *entity.TripItem, tx *sqlx.Tx) error
	DeleteByTripIDCommand(ctx context.Context, tripID int64, tx *sqlx.Tx) error
}
