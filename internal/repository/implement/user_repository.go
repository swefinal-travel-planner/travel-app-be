package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db database.Db) repository.UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateCommand(ctx context.Context, user *entity.User, tx *sqlx.Tx) error {
	// Insert the new user
	insertQuery := `INSERT INTO users(email, name, phone_number, id_token, photo_url, password) VALUES (:email, :name, :phone_number, :id_token, :photo_url, :password)`
	if tx != nil {
		_, err := tx.NamedExecContext(ctx, insertQuery, user)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, insertQuery, user)
	return err
}

func (repo *UserRepository) GetOneByEmailQuery(ctx context.Context, email string, tx *sqlx.Tx) (*entity.User, error) {
	var customer entity.User
	query := "SELECT * FROM users WHERE email = ? AND users.deleted_at IS NULL"
	if tx != nil {
		err := tx.QueryRowxContext(ctx, query, email).StructScan(&customer)
		return &customer, err
	}
	err := repo.db.QueryRowxContext(ctx, query, email).StructScan(&customer)
	return &customer, err
}

func (repo *UserRepository) GetIdByEmailQuery(ctx context.Context, email string, tx *sqlx.Tx) (int64, error) {
	var user entity.User
	query := "SELECT * FROM users WHERE email = ? AND users.deleted_at IS NULL"
	if tx != nil {
		err := tx.QueryRowxContext(ctx, query, email).StructScan(&user)
		return user.Id, err
	}
	err := repo.db.QueryRowxContext(ctx, query, email).StructScan(&user)
	return user.Id, err
}

func (repo *UserRepository) UpdatePasswordByIdQuery(ctx context.Context, id int64, password string, tx *sqlx.Tx) error {
	query := "UPDATE users SET password = ? WHERE id = ?"
	if tx != nil {
		_, err := tx.ExecContext(ctx, query, password, id)
		return err
	}
	_, err := repo.db.ExecContext(ctx, query, password, id)
	return err
}

func (repo *UserRepository) GetOneByIDQuery(ctx context.Context, id int64, tx *sqlx.Tx) (*entity.User, error) {
	var customer entity.User
	query := "SELECT * FROM users WHERE id = ? AND users.deleted_at IS NULL"
	if tx != nil {
		err := tx.QueryRowxContext(ctx, query, id).StructScan(&customer)
		return &customer, err
	}
	err := repo.db.QueryRowxContext(ctx, query, id).StructScan(&customer)
	return &customer, err
}

func (repo *UserRepository) UpdateNotificationTokenCommand(ctx context.Context, id int64, token string, tx *sqlx.Tx) error {
	query := "UPDATE users SET notification_token = ? WHERE id = ?"
	if tx != nil {
		_, err := tx.ExecContext(ctx, query, token, id)
		return err
	}
	_, err := repo.db.ExecContext(ctx, query, token, id)
	return err
}
