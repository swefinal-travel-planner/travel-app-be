package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type NotificationRepository interface {
	CreateCommand(ctx context.Context, notification *entity.Notification, tx *sqlx.Tx) error
	GetAllByUserIDQuery(ctx context.Context, userID int64, typeFilter string, tx *sqlx.Tx) ([]*entity.Notification, error)
	SeenNotificationCommand(ctx context.Context, userID int64, notificationID int64, tx *sqlx.Tx) error
	GetOneByUserIdAndTypeAndTriggerEntityIDQuery(ctx context.Context, userId int64, typeFilter string, triggerEntityID int64, tx *sqlx.Tx) (*entity.Notification, error)
	DeleteNotificationCommand(ctx context.Context, notificationID int64, tx *sqlx.Tx) error
	DeleteTripNotificationCommand(ctx context.Context, receiverId, sender, tripId int64, tx *sqlx.Tx) error
}
