package repository

import (
	"context"

	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type FriendRepository interface {
	CreateCommand(ctx context.Context, friend *entity.Friend) error
	GetByUserIdQuery(ctx context.Context, userId int64) ([]*entity.User, error)
	DeleteByUserId1AndUserId2Command(ctx context.Context, userId1 int64, userId2 int64) error
	GetByUserId1AndUserId2Query(ctx context.Context, userId1 int64, userId2 int64) error
}
