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

func (repo *UserRepository) CreateCommand(ctx context.Context, user *entity.User) error {
	// Insert the new user
	insertQuery := `INSERT INTO users(email, name, phone_number, password) VALUES (:email, :name, :phone_number, :password)`
	_, err := repo.db.NamedExecContext(ctx, insertQuery, user)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetOneByEmailQuery(ctx context.Context, email string) (*entity.User, error) {
	var customer entity.User
	query := "SELECT * FROM users WHERE email = ? AND users.deleted_at IS NULL"
	err := repo.db.QueryRowxContext(ctx, query, email).StructScan(&customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (repo *UserRepository) GetIdByEmailQuery(ctx context.Context, email string) (int64, error) {
	var user entity.User
	query := "SELECT * FROM users WHERE email = ? AND users.deleted_at IS NULL"
	err := repo.db.QueryRowxContext(ctx, query, email).StructScan(&user)
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (repo *UserRepository) UpdatePasswordByIdQuery(ctx context.Context, id int64, password string) error {
	query := "UPDATE users SET password = ? WHERE id = ? AND users.deleted_at IS NULL"
	_, err := repo.db.ExecContext(ctx, query, password, id)
	if err != nil {
		return err
	}

	return nil
}
