package repository

import (
	"context"

	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type UserRepository interface {
	CreateCommand(ctx context.Context, user *entity.User) error
}
