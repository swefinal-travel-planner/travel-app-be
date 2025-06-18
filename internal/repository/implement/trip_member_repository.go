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

func (repo *TripMemberRepository) IsUserTripAdminQuery(ctx context.Context, tripID int64, userID int64, tx *sqlx.Tx) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM trip_members 
		WHERE trip_id = ? AND user_id = ? AND role IN ('administrator') AND deleted_at IS NULL
	`
	if tx != nil {
		err := tx.GetContext(ctx, &count, query, tripID, userID)
		return count > 0, err
	}
	err := repo.db.GetContext(ctx, &count, query, tripID, userID)
	return count > 0, err
}

func (repo *TripMemberRepository) GetTripMembersQuery(ctx context.Context, tripID int64, tx *sqlx.Tx) ([]entity.TripMemberWithUser, error) {
	members := make([]entity.TripMemberWithUser, 0)
	query := `
		SELECT tm.*, u.name as name, u.photo_url 
		FROM trip_members tm
		JOIN users u ON tm.user_id = u.id
		WHERE tm.trip_id = ? AND tm.deleted_at IS NULL
		ORDER BY tm.created_at ASC
	`
	if tx != nil {
		err := tx.SelectContext(ctx, &members, query, tripID)
		return members, err
	}
	err := repo.db.SelectContext(ctx, &members, query, tripID)
	return members, err
}

func (repo *TripMemberRepository) DeleteMemberCommand(ctx context.Context, tripID int64, userID int64, tx *sqlx.Tx) error {
	query := `
		DELETE FROM trip_members 
		WHERE trip_id = ? AND user_id = ?
	`
	if tx != nil {
		_, err := tx.ExecContext(ctx, query, tripID, userID)
		return err
	}
	_, err := repo.db.ExecContext(ctx, query, tripID, userID)
	return err
}
