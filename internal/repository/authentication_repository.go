package repository

import (
	"context"

	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type AuthenticationRepository interface {
	CreateCommand(ctx context.Context, refreshToken entity.Authentication) error
	UpdateCommand(ctx context.Context, refreshToken entity.Authentication) error
	GetOneByUserIdQuery(ctx context.Context, userId int64) (*entity.Authentication, error)
	DeleteByRefreshToken(ctx context.Context, refreshToken string) error
}
