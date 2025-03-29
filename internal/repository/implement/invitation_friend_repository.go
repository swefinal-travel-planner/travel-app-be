package repositoryimplement

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
)

type InvitationFriendRepository struct {
	db *sqlx.DB
}

func NewInvitationFriendRepository(db database.Db) repository.InvitationFriendRepository {
	return &InvitationFriendRepository{db: db}
}

func (repo *InvitationFriendRepository) CreateCommand(ctx context.Context, invitation *entity.InvitationFriend) error {
	// Insert the new invitation
	insertQuery := `INSERT INTO invitation_friends(sender_id, receiver_id, status) VALUES (:sender_id, :receiver_id, :status)`
	_, err := repo.db.NamedExecContext(ctx, insertQuery, invitation)
	if err != nil {
		return err
	}
	return nil
}

func (repo *InvitationFriendRepository) GetByReceiverIdCommand(ctx context.Context, receiverId int64) ([]*entity.InvitationFriend, error) {
	var invitationFriend []*entity.InvitationFriend
	query := "SELECT * FROM invitation_friends WHERE status = 'pending' AND receiver_id = ?"
	err := repo.db.SelectContext(ctx, &invitationFriend, query, receiverId)
	if err != nil {
		return nil, err
	}
	return invitationFriend, nil
}

func (repo *InvitationFriendRepository) GetBySenderAndReceiverIdQuery(ctx context.Context, senderId, receiverId int64) (*entity.InvitationFriend, error) {
	var invitationFriend entity.InvitationFriend
	query := `
		SELECT * 
		FROM invitation_friends 
		WHERE 
			((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))
	`
	err := repo.db.GetContext(ctx, &invitationFriend, query, senderId, receiverId, receiverId, senderId)
	if err != nil {
		return nil, err
	}
	return &invitationFriend, nil
}
