package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type TripRepository interface {
	CreateCommand(ctx context.Context, trip *entity.Trip, tx *sqlx.Tx) (int64, error)
	GetOneByIDQuery(ctx context.Context, id int64, tx *sqlx.Tx) (*entity.Trip, error)
	GetAllByUserIDQuery(ctx context.Context, userId int64, tx *sqlx.Tx) ([]*entity.Trip, error)
	GetAllWithUserRoleByUserIdQuery(ctx context.Context, userId int64, tx *sqlx.Tx) ([]*entity.TripWithRole, error)
	GetOneWithUserRoleByIDQuery(ctx context.Context, tripId int64, userId int64, tx *sqlx.Tx) (*entity.TripWithRole, error)
	SelectForUpdateById(ctx context.Context, id int64, tx *sqlx.Tx) (*entity.Trip, error)
	SelectForShareById(ctx context.Context, id int64, tx *sqlx.Tx) (*entity.Trip, error)
	UpdateCommand(ctx context.Context, trip *entity.Trip, tx *sqlx.Tx) error
}
