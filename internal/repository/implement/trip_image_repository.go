package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
)

type TripImageRepository struct {
	db *sqlx.DB
}

func NewTripImageRepository(db database.Db) repository.TripImageRepository {
	return &TripImageRepository{db: db}
}

func (repo *TripImageRepository) CreateCommand(ctx context.Context, tripImage *entity.TripImage, tx *sqlx.Tx) error {
	insertQuery := `
	INSERT INTO trip_images(
		trip_id, image_url
	) 
	VALUES (
		:trip_id, :image_url
	)
	`
	if tx != nil {
		_, err := tx.NamedExecContext(ctx, insertQuery, tripImage)
		return err
	}

	_, err := repo.db.NamedExecContext(ctx, insertQuery, tripImage)
	return err
}

func (repo *TripImageRepository) GetAllQuery(ctx context.Context, tripID int64, tx *sqlx.Tx) ([]entity.TripImage, error) {
	var tripImages []entity.TripImage
	query := "SELECT * FROM trip_images WHERE trip_id = ? AND deleted_at IS NULL ORDER BY created_at ASC"

	if tx != nil {
		err := tx.SelectContext(ctx, &tripImages, query, tripID)
		return tripImages, err
	}

	err := repo.db.SelectContext(ctx, &tripImages, query, tripID)
	return tripImages, err
}

func (repo *TripImageRepository) DeleteOneByIDCommand(ctx context.Context, id int64, tx *sqlx.Tx) error {
	deleteQuery := "DELETE FROM trip_images WHERE id = ?"
	if tx != nil {
		_, err := tx.ExecContext(ctx, deleteQuery, id)
		return err
	}
	_, err := repo.db.ExecContext(ctx, deleteQuery, id)
	return err
}
