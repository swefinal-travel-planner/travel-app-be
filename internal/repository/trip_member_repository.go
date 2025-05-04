package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type TripMemberRepository interface {
	CreateCommand(ctx context.Context, tripMember *entity.TripMember, tx *sqlx.Tx) error
}
