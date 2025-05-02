package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type FriendRepository interface {
	CreateCommand(ctx context.Context, friend *entity.Friend, tx *sqlx.Tx) error
	GetByUserIdQuery(ctx context.Context, userId int64, tx *sqlx.Tx) ([]*entity.User, error)
	DeleteByUserId1AndUserId2Command(ctx context.Context, userId1 int64, userId2 int64, tx *sqlx.Tx) error
	ExistsByUserId1AndUserId2Query(ctx context.Context, userId1 int64, userId2 int64, tx *sqlx.Tx) bool
}
