package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
)

type TripRepository struct {
	db *sqlx.DB
}

func NewTripRepository(db database.Db) repository.TripRepository {
	return &TripRepository{db: db}
}

func (repo *TripRepository) CreateCommand(ctx context.Context, trip *entity.Trip, tx *sqlx.Tx) (int64, error) {
	// Insert the new trip
	insertQuery := `
	INSERT INTO trips(
		title, city, start_date, days, budget, num_members, 
		vi_location_attributes, vi_food_attributes, vi_special_requirements, vi_medical_conditions,
		en_location_attributes, en_food_attributes, en_special_requirements, en_medical_conditions,
		status
	) 
	VALUES (
		:title, :city, :start_date, :days, :budget, :num_members, 
		:vi_location_attributes, :vi_food_attributes, :vi_special_requirements, :vi_medical_conditions,
		:en_location_attributes, :en_food_attributes, :en_special_requirements, :en_medical_conditions,
		:status
	)
	`
	if tx != nil {
		result, err := tx.NamedExecContext(ctx, insertQuery, trip)
		if err != nil {
			return 0, err
		}

		id, err := result.LastInsertId()
		return id, err
	}

	result, err := repo.db.NamedExecContext(ctx, insertQuery, trip)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

func (repo *TripRepository) GetOneByIDQuery(ctx context.Context, id int64, tx *sqlx.Tx) (*entity.Trip, error) {
	var trip entity.Trip
	query := "SELECT * FROM trips WHERE id = ?"
	if tx != nil {
		err := tx.GetContext(ctx, &trip, query, id)
		return &trip, err
	}
	err := repo.db.GetContext(ctx, &trip, query, id)
	return &trip, err
}

func (repo *TripRepository) LockTripRowByIDCommand(ctx context.Context, id int64, tx *sqlx.Tx) (*entity.Trip, error) {
	var trip entity.Trip
	query := "SELECT * FROM trips WHERE id = ? FOR UPDATE"
	if tx != nil {
		err := tx.GetContext(ctx, &trip, query, id)
		return &trip, err
	}
	err := tx.GetContext(ctx, &trip, query, id)
	return &trip, err
}
