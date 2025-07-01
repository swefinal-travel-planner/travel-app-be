package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type InvitationTripRepository interface {
	CreateCommand(ctx context.Context, invitation *entity.InvitationTrip, tx *sqlx.Tx) error
	GetByTripIDQuery(ctx context.Context, tripId int64, tx *sqlx.Tx) ([]entity.InvitationTrip, error)
	GetPendingByTripIDQuery(ctx context.Context, tripId int64, tx *sqlx.Tx) ([]entity.InvitationTrip, error)
	GetByReceiverIDQuery(ctx context.Context, receiverId int64, tx *sqlx.Tx) ([]entity.InvitationTrip, error)
	GetBySenderIDQuery(ctx context.Context, senderId int64, tx *sqlx.Tx) ([]entity.InvitationTrip, error)
	GetOneByIDQuery(ctx context.Context, id int64, tx *sqlx.Tx) (*entity.InvitationTrip, error)
	GetOneByReceiverIdAndTripIDQuery(ctx context.Context, userId int64, tripId int64, tx *sqlx.Tx) (*entity.InvitationTrip, error)
	DeleteByIDCommand(ctx context.Context, id int64, tx *sqlx.Tx) error
}
