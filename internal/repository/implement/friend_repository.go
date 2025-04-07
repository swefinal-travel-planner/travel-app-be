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

func (repo *FriendRepository) GetByUserIdQuery(ctx context.Context, userId int64) ([]*entity.User, error) {
	var users []*entity.User
	query := `
		SELECT u.id, u.name, u.photo_url 
		FROM friends f
		JOIN users u ON f.user_id_2 = u.id
		WHERE f.user_id_1 = ?

		UNION

		SELECT u.id, u.name, u.photo_url 
		FROM friends f
		JOIN users u ON f.user_id_1 = u.id
		WHERE f.user_id_2 = ?;
	`
	err := repo.db.SelectContext(ctx, &users, query, userId, userId)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *FriendRepository) DeleteByUserId1AndUserId2Command(ctx context.Context, userId1 int64, userId2 int64) error {
	// Delete friend by userId
	deleteQuery := `
		DELETE FROM friends
		WHERE (user_id_1 = ? AND user_id_2 = ?)
		   OR (user_id_1 = ? AND user_id_2 = ?)
	`
	_, err := repo.db.ExecContext(ctx, deleteQuery, userId1, userId2, userId2, userId1)
	if err != nil {
		return err
	}
	return nil
}

func (repo *FriendRepository) GetByUserId1AndUserId2Query(ctx context.Context, userId1 int64, userId2 int64) error {
	var friend entity.Friend
	query := `
		SELECT * FROM friends 
		WHERE (user_id_1 = ? AND user_id_2 = ?) 
		   OR (user_id_1 = ? AND user_id_2 = ?)
	`
	err := repo.db.GetContext(ctx, &friend, query, userId1, userId2, userId2, userId1)
	if err != nil {
		return err
	}
	return nil
}
