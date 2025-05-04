package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
)

type TripMemberRepository struct {
	db *sqlx.DB
}

func NewTripMemberRepository(db database.Db) repository.TripMemberRepository {
	return &TripMemberRepository{db: db}
}

func (repo *TripMemberRepository) CreateCommand(ctx context.Context, tripMember *entity.TripMember, tx *sqlx.Tx) error {
	// Insert the new trip member
	insertQuery := `
	INSERT INTO trip_members(
		trip_id, user_id, role
	) 
	VALUES (
		:trip_id, :user_id, :role
	)
	`
	if tx != nil {
		_, err := tx.NamedExecContext(ctx, insertQuery, tripMember)
		return err
	}

	_, err := repo.db.NamedExecContext(ctx, insertQuery, tripMember)
	return err
}

func (repo *TripMemberRepository) IsUserInTripQuery(ctx context.Context, tripID int64, userID int64, tx *sqlx.Tx) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM trip_members 
		WHERE trip_id = ? AND user_id = ? AND deleted_at IS NULL
	`
	if tx != nil {
		err := tx.GetContext(ctx, &count, query, tripID, userID)
		return count > 0, err
	}
	err := repo.db.GetContext(ctx, &count, query, tripID, userID)
	return count > 0, err
}

func (repo *TripMemberRepository) IsUserTripAdminOrStaffQuery(ctx context.Context, tripID int64, userID int64, tx *sqlx.Tx) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM trip_members 
		WHERE trip_id = ? AND user_id = ? AND role IN ('administrator', 'staff') AND deleted_at IS NULL
	`
	if tx != nil {
		err := tx.GetContext(ctx, &count, query, tripID, userID)
		return count > 0, err
	}
	err := repo.db.GetContext(ctx, &count, query, tripID, userID)
	return count > 0, err
}
