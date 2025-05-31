package repositoryimplement

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
)

type InvitationTripRepository struct {
	db *sqlx.DB
}

func NewInvitationTripRepository(db database.Db) repository.InvitationTripRepository {
	return &InvitationTripRepository{db: db}
}

func (repo *InvitationTripRepository) CreateCommand(ctx context.Context, invitation *entity.InvitationTrip, tx *sqlx.Tx) error {
	insertQuery := `INSERT INTO invitation_trips(trip_id, sender_id, receiver_id, status) VALUES (:trip_id, :sender_id, :receiver_id, :status)`
	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.NamedExecContext(ctx, insertQuery, invitation)
	} else {
		result, err = repo.db.NamedExecContext(ctx, insertQuery, invitation)
	}

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	invitation.ID = id
	return nil
}

func (repo *InvitationTripRepository) GetByTripIDQuery(ctx context.Context, tripId int64, tx *sqlx.Tx) ([]entity.InvitationTrip, error) {
	var invitations []entity.InvitationTrip
	query := "SELECT * FROM invitation_trips WHERE trip_id = ?"
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &invitations, query, tripId)
	} else {
		err = repo.db.SelectContext(ctx, &invitations, query, tripId)
	}
	if err != nil && err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
		return make([]entity.InvitationTrip, 0), nil
	}
	return invitations, err
}

func (repo *InvitationTripRepository) GetByReceiverIDQuery(ctx context.Context, receiverId int64, tx *sqlx.Tx) ([]entity.InvitationTrip, error) {
	var invitations []entity.InvitationTrip
	query := "SELECT * FROM invitation_trips WHERE receiver_id = ?"
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &invitations, query, receiverId)
	} else {
		err = repo.db.SelectContext(ctx, &invitations, query, receiverId)
	}
	if err != nil && err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
		return make([]entity.InvitationTrip, 0), nil
	}
	return invitations, err
}

func (repo *InvitationTripRepository) GetBySenderIDQuery(ctx context.Context, senderId int64, tx *sqlx.Tx) ([]entity.InvitationTrip, error) {
	var invitations []entity.InvitationTrip
	query := "SELECT * FROM invitation_trips WHERE sender_id = ?"
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &invitations, query, senderId)
	} else {
		err = repo.db.SelectContext(ctx, &invitations, query, senderId)
	}
	if err != nil && err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
		return make([]entity.InvitationTrip, 0), nil
	}
	return invitations, err
}

func (repo *InvitationTripRepository) GetOneByIDQuery(ctx context.Context, id int64, tx *sqlx.Tx) (*entity.InvitationTrip, error) {
	var invitation entity.InvitationTrip
	query := "SELECT * FROM invitation_trips WHERE id = ?"
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &invitation, query, id)
	} else {
		err = repo.db.GetContext(ctx, &invitation, query, id)
	}
	if err != nil && err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
		return nil, nil
	}
	return &invitation, err
}

func (repo *InvitationTripRepository) DeleteByIDCommand(ctx context.Context, id int64, tx *sqlx.Tx) error {
	query := "DELETE FROM invitation_trips WHERE id = ?"
	if tx != nil {
		_, err := tx.ExecContext(ctx, query, id)
		return err
	}
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *InvitationTripRepository) GetOneByReceiverIdAndTripIDQuery(ctx context.Context, userId int64, tripId int64, tx *sqlx.Tx) (*entity.InvitationTrip, error) {
	var invitation entity.InvitationTrip
	query := "SELECT * FROM invitation_trips WHERE receiver_id = ? AND trip_id = ?"
	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &invitation, query, userId, tripId)
	} else {
		err = repo.db.GetContext(ctx, &invitation, query, userId, tripId)
	}
	if err != nil && err.Error() == error_utils.SystemErrorMessage.SqlxNoRow {
		return nil, nil
	}
	return &invitation, err
}
