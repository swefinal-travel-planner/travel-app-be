package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type AuthenticationRepository interface {
	CreateCommand(ctx context.Context, refreshToken entity.Authentication, tx *sqlx.Tx) error
	UpdateCommand(ctx context.Context, refreshToken entity.Authentication, tx *sqlx.Tx) error
	GetOneByUserIdQuery(ctx context.Context, userId int64, tx *sqlx.Tx) (*entity.Authentication, error)
	DeleteByRefreshToken(ctx context.Context, refreshToken string, tx *sqlx.Tx) error
}
