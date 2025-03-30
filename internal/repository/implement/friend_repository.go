package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
)

type FriendRepository struct {
	db *sqlx.DB
}

func NewFriendRepository(db database.Db) repository.FriendRepository {
	return &FriendRepository{db: db}
}

func (repo *FriendRepository) CreateCommand(ctx context.Context, friend *entity.Friend) error {
	// Insert new friend
	insertQuery := `INSERT INTO friends(user_id_1, user_id_2) VALUES (:user_id_1, :user_id_2)`
	_, err := repo.db.NamedExecContext(ctx, insertQuery, friend)
	if err != nil {
		return err
	}
	return nil
}
