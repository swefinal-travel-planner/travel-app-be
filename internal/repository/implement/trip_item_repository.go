package repositoryimplement

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
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

// avoid TOCTOU
func (repo *TripItemRepository) GetTripItemsByTripIDCommand(ctx context.Context, tripID int64, userId int64, tx *sqlx.Tx) ([]entity.TripItem, error) {
	query := `
		SELECT tm.user_id as 'userId', ti.*
		FROM trip_members tm
		LEFT JOIN trip_items ti ON ti.trip_id = tm.trip_id
		WHERE tm.trip_id = ? AND tm.user_id = ? AND tm.deleted_at IS NULL
	`

	var tripItems []entity.TripItem
	var tripItemsWithUserId []struct {
		UserId int64 `db:"userId"`
		entity.TripItemTOCTOU
	}

	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &tripItemsWithUserId, query, tripID, userId)
	} else {
		err = repo.db.SelectContext(ctx, &tripItemsWithUserId, query, tripID, userId)
	}
	if err != nil {
		return nil, err
	}
	if len(tripItemsWithUserId) == 0 {
		return nil, errors.New(error_utils.SystemErrorMessage.NotMemberOfTrip)
	}

	if tripItemsWithUserId[0].TripID.Int64 == 0 {
		return make([]entity.TripItem, 0), nil
	}

	for _, tripItem := range tripItemsWithUserId {
		tripItems = append(tripItems, entity.TripItem{
			ID:         tripItem.TripID.Int64,
			TripID:     tripItem.TripID.Int64,
			PlaceID:    tripItem.PlaceID.String,
			TripDay:    tripItem.TripDay.Int64,
			OrderInDay: tripItem.OrderInDay.Int64,
			TimeInDate: tripItem.TimeInDate.String,
			CreatedAt:  tripItem.CreatedAt.Time,
			UpdatedAt:  tripItem.UpdatedAt.Time,
			DeletedAt:  tripItem.DeletedAt,
		})
	}
	return tripItems, nil
}
