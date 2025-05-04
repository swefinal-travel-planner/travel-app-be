package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type TripMemberRepository interface {
	CreateCommand(ctx context.Context, tripMember *entity.TripMember, tx *sqlx.Tx) error
	IsUserInTripQuery(ctx context.Context, tripID int64, userID int64, tx *sqlx.Tx) (bool, error)
	IsUserTripAdminOrStaffQuery(ctx context.Context, tripID int64, userID int64, tx *sqlx.Tx) (bool, error)
}
