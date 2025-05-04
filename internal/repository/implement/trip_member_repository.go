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
