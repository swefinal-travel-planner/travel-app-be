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

func (repo *FriendRepository) CreateCommand(ctx context.Context, friend *entity.Friend, tx *sqlx.Tx) error {
	// Insert new friend
	insertQuery := `INSERT INTO friends(user_id_1, user_id_2) VALUES (:user_id_1, :user_id_2)`
	if tx != nil {
		_, err := tx.NamedExecContext(ctx, insertQuery, friend)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, insertQuery, friend)
	return err
}

func (repo *FriendRepository) GetByUserIdQuery(ctx context.Context, userId int64, tx *sqlx.Tx) ([]*entity.User, error) {
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
	if tx != nil {
		err := tx.SelectContext(ctx, &users, query, userId, userId)
		return users, err
	}
	err := repo.db.SelectContext(ctx, &users, query, userId, userId)
	return users, err
}

func (repo *FriendRepository) DeleteByUserId1AndUserId2Command(ctx context.Context, userId1 int64, userId2 int64, tx *sqlx.Tx) error {
	// Delete friend by userId
	deleteQuery := `
		DELETE FROM friends
		WHERE (user_id_1 = ? AND user_id_2 = ?)
		   OR (user_id_1 = ? AND user_id_2 = ?)
	`
	if tx != nil {
		_, err := tx.ExecContext(ctx, deleteQuery, userId1, userId2, userId2, userId1)
		return err
	}
	_, err := repo.db.ExecContext(ctx, deleteQuery, userId1, userId2, userId2, userId1)
	return err
}

func (repo *FriendRepository) ExistsByUserId1AndUserId2Query(ctx context.Context, userId1 int64, userId2 int64, tx *sqlx.Tx) bool {
	var count int
	query := `
		SELECT COUNT(*) FROM friends 
		WHERE (user_id_1 = ? AND user_id_2 = ?) 
		   OR (user_id_1 = ? AND user_id_2 = ?)
	`
	if tx != nil {
		err := tx.GetContext(ctx, &count, query, userId1, userId2, userId2, userId1)
		return err == nil && count > 0
	}
	err := repo.db.GetContext(ctx, &count, query, userId1, userId2, userId2, userId1)
	return err == nil && count > 0
}
