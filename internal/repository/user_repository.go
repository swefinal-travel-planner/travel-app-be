package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type UserRepository interface {
	CreateCommand(ctx context.Context, user *entity.User, tx *sqlx.Tx) error
	GetOneByEmailQuery(ctx context.Context, email string, tx *sqlx.Tx) (*entity.User, error)
	GetIdByEmailQuery(ctx context.Context, email string, tx *sqlx.Tx) (int64, error)
	UpdatePasswordByIdQuery(ctx context.Context, id int64, password string, tx *sqlx.Tx) error
	GetOneByIDQuery(ctx context.Context, id int64, tx *sqlx.Tx) (*entity.User, error)
	UpdateNotificationTokenCommand(ctx context.Context, id int64, token string, tx *sqlx.Tx) error
}
