package repository

import (
	"context"

	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type UserRepository interface {
	CreateCommand(ctx context.Context, user *entity.User) error
	GetOneByEmailQuery(ctx context.Context, email string) (*entity.User, error)
	GetIdByEmailQuery(ctx context.Context, email string) (int64, error)
	UpdatePasswordByIdQuery(ctx context.Context, id int64, password string) error
	GetOneByIDQuery(ctx context.Context, id int64) (*entity.User, error)
}
