package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
)

type AuthenticationRepository struct {
	db *sqlx.DB
}

func NewAuthenticationRepository(db database.Db) repository.AuthenticationRepository {
	return &AuthenticationRepository{db: db}
}

func (repo *AuthenticationRepository) CreateCommand(ctx context.Context, authentication entity.Authentication, tx *sqlx.Tx) error {
	query := `
		INSERT INTO authentications (user_id, refresh_token)
		VALUES (:user_id, :refresh_token)
	`
	if tx != nil {
		_, err := tx.NamedExecContext(ctx, query, authentication)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, query, authentication)
	return err
}

func (repo *AuthenticationRepository) UpdateCommand(ctx context.Context, authentication entity.Authentication, tx *sqlx.Tx) error {
	query := `
		UPDATE authentications
		SET refresh_token = :refresh_token
		WHERE user_id = :user_id
	`
	if tx != nil {
		_, err := tx.NamedExecContext(ctx, query, authentication)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, query, authentication)
	return err
}

// find refresh token
func (repo *AuthenticationRepository) GetOneByUserIdQuery(ctx context.Context, userId int64, tx *sqlx.Tx) (*entity.Authentication, error) {
	var authentication entity.Authentication
	query := `
		SELECT user_id, refresh_token, created_at
		FROM authentications
		WHERE user_id = ?
	`
	if tx != nil {
		err := tx.GetContext(ctx, &authentication, query, userId)
		return &authentication, err
	}
	err := repo.db.GetContext(ctx, &authentication, query, userId)
	return &authentication, err
}

func (repo *AuthenticationRepository) DeleteByRefreshToken(ctx context.Context, refreshToken string, tx *sqlx.Tx) error {
	query := `DELETE FROM authentications WHERE refresh_token = ?`
	if tx != nil {
		_, err := tx.ExecContext(ctx, query, refreshToken)
		return err
	}
	_, err := repo.db.ExecContext(ctx, query, refreshToken)
	return err
}
