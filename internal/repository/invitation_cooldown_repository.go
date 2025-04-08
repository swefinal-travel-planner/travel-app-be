package repository

import (
	"context"

	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type InvitationCooldownRepository interface {
	CreateCommand(ctx context.Context, cooldown *entity.InvitationCooldown) error
	GetLatestCooldownBetweenUsersQuery(ctx context.Context, userID1, userID2 int64) (*entity.InvitationCooldown, error)
}
