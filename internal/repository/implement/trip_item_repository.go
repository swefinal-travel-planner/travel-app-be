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

func (repo *TripItemRepository) CreateCommand(ctx context.Context, tripItem *entity.TripItem) error {
	// Insert the new trip item
	insertQuery := `
	INSERT INTO trip_items(
		id, trip_id, place_id, trip_day, order_in_day, tag
	) 
	VALUES (
		:id, :trip_id, :place_id, :trip_day, :order_in_day, :tag
	)
	`
	_, err := repo.db.NamedExecContext(ctx, insertQuery, tripItem)
	if err != nil {
		return err
	}

	return nil
}
