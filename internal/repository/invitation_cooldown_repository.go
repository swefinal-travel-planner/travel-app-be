package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type InvitationCooldownRepository interface {
	CreateCommand(ctx context.Context, cooldown *entity.InvitationCooldown, tx *sqlx.Tx) error
	GetLatestCooldownBetweenUsersQuery(ctx context.Context, userID1, userID2 int64, tx *sqlx.Tx) (*entity.InvitationCooldown, error)
}
