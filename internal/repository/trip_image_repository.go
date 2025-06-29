package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type TripImageRepository interface {
	CreateCommand(ctx context.Context, tripImage *entity.TripImage, tx *sqlx.Tx) error
	GetAllQuery(ctx context.Context, tripID int64, tx *sqlx.Tx) ([]entity.TripImage, error)
	GetAllWithUserInfoQuery(ctx context.Context, tripID int64, tx *sqlx.Tx) ([]entity.TripImageWithUserInfo, error)
	DeleteOneByIDCommand(ctx context.Context, id int64, tx *sqlx.Tx) error
}
