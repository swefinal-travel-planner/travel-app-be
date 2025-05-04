package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
)

type TripItemRepository struct {
	db *sqlx.DB
}

func NewTripItemRepository(db database.Db) repository.TripItemRepository {
	return &TripItemRepository{db: db}
}

func (repo *TripItemRepository) CreateCommand(ctx context.Context, tripItem *entity.TripItem, tx *sqlx.Tx) error {
	// Insert the new trip item
	insertQuery := `
	INSERT INTO trip_items(
		id, trip_id, place_id, trip_day, order_in_day, time_in_date
	) 
	VALUES (
		:id, :trip_id, :place_id, :trip_day, :order_in_day, :time_in_date
	)
	`
	if tx != nil {
		_, err := tx.NamedExecContext(ctx, insertQuery, tripItem)
		return err
	}

	_, err := repo.db.NamedExecContext(ctx, insertQuery, tripItem)
	return err
}

func (repo *TripItemRepository) DeleteByTripIDCommand(ctx context.Context, tripID int64, tx *sqlx.Tx) error {
	query := "DELETE FROM trip_items WHERE trip_id = ?"
	if tx != nil {
		_, err := tx.ExecContext(ctx, query, tripID)
		return err
	}
	_, err := repo.db.ExecContext(ctx, query, tripID)
	return err
}
