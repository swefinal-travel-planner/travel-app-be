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

func (repo *TripRepository) CreateCommand(ctx context.Context, trip *entity.Trip) (int64, error) {
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
	result, err := repo.db.NamedExecContext(ctx, insertQuery, trip)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
