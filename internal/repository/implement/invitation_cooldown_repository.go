package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
)

type InvitationCooldownRepository struct {
	db *sqlx.DB
}

func NewInvitationCooldownRepository(db database.Db) repository.InvitationCooldownRepository {
	return &InvitationCooldownRepository{db: db}
}

func (repo *InvitationCooldownRepository) CreateCommand(ctx context.Context, cooldown *entity.InvitationCooldown) error {
	insertQuery := `
		INSERT INTO invitation_cooldowns(
			user_id_1, 
			user_id_2, 
			start_cooldown_millis, 
			cooldown_duration
		) VALUES (
			:user_id_1, 
			:user_id_2, 
			:start_cooldown_millis, 
			:cooldown_duration
		)`
	_, err := repo.db.NamedExecContext(ctx, insertQuery, cooldown)
	if err != nil {
		return err
	}
	return nil
}

func (repo *InvitationCooldownRepository) GetLatestCooldownBetweenUsersQuery(ctx context.Context, userID1, userID2 int64) (*entity.InvitationCooldown, error) {
	var cooldown entity.InvitationCooldown
	query := `
		SELECT * 
		FROM invitation_cooldowns 
		WHERE 
			((user_id_1 = ? AND user_id_2 = ?) OR (user_id_1 = ? AND user_id_2 = ?))
			AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := repo.db.GetContext(ctx, &cooldown, query, userID1, userID2, userID2, userID1)
	if err != nil {
		return nil, err
	}
	return &cooldown, nil
}
