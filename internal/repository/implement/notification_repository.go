package repositoryimplement

import (
	"context"
	"database/sql"
	"strings"

	myDatabase "github.com/swefinal-travel-planner/travel-app-be/internal/database"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
)

type notificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db myDatabase.Db) repository.NotificationRepository {
	return &notificationRepository{
		db: db,
	}
}

func (r *notificationRepository) CreateCommand(ctx context.Context, notification *entity.Notification, tx *sqlx.Tx) error {
	query := `
		INSERT INTO notifications (
			user_id,
			type,
			trigger_entity_type,
			trigger_entity_avatar,
			trigger_entity_name,
			trigger_entity_id,
			reference_entity_type,
			reference_entity_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.ExecContext(
			ctx,
			query,
			notification.UserID,
			notification.Type,
			notification.TriggerEntityType,
			notification.TriggerEntityAvatar,
			notification.TriggerEntityName,
			notification.TriggerEntityID,
			notification.ReferenceEntityType,
			notification.ReferenceEntityID,
		)
	} else {
		result, err = r.db.ExecContext(
			ctx,
			query,
			notification.UserID,
			notification.Type,
			notification.TriggerEntityType,
			notification.TriggerEntityAvatar,
			notification.TriggerEntityName,
			notification.TriggerEntityID,
			notification.ReferenceEntityType,
			notification.ReferenceEntityID,
		)
	}

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	notification.ID = id
	return nil
}

func (r *notificationRepository) GetAllByUserIDQuery(ctx context.Context, userID int64, typeFilter string, tx *sqlx.Tx) ([]*entity.Notification, error) {
	query := `
		SELECT * FROM notifications WHERE user_id = ?
	`

	var notifications []*entity.Notification

	var err error
	if typeFilter != "" {
		query += " AND type IN ('" + strings.Replace(typeFilter, ",", "','", -1) + "')"
	}

	if tx != nil {
		err = tx.SelectContext(ctx, &notifications, query, userID)
	} else {
		err = r.db.SelectContext(ctx, &notifications, query, userID)
	}

	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *notificationRepository) SeenNotificationCommand(ctx context.Context, userID int64, notificationID int64, tx *sqlx.Tx) error {
	query := `
		UPDATE notifications SET is_seen = 1 WHERE id = ? AND user_id = ?
	`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, notificationID, userID)
	} else {
		_, err = r.db.ExecContext(ctx, query, notificationID, userID)
	}

	return err
}

func (r *notificationRepository) GetOneByUserIdAndTypeAndTriggerEntityIDQuery(ctx context.Context, userId int64, typeFilter string, triggerEntityID int64, tx *sqlx.Tx) (*entity.Notification, error) {
	query := `
		SELECT * FROM notifications WHERE type = ? AND trigger_entity_id = ? AND user_id = ?
	`

	var notification entity.Notification

	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &notification, query, typeFilter, triggerEntityID, userId)
	} else {
		err = r.db.GetContext(ctx, &notification, query, typeFilter, triggerEntityID, userId)
	}

	if err != nil {
		return nil, err
	}

	return &notification, nil
}

func (r *notificationRepository) DeleteNotificationCommand(ctx context.Context, notificationID int64, tx *sqlx.Tx) error {
	query := `
		DELETE FROM notifications WHERE id = ?
	`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, notificationID)
	} else {
		_, err = r.db.ExecContext(ctx, query, notificationID)
	}

	return err
}
