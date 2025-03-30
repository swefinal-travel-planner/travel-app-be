package repository

import (
	"context"

	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type FriendRepository interface {
	CreateCommand(ctx context.Context, friend *entity.Friend) error
}
