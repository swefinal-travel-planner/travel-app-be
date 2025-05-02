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

func (repo *InvitationFriendRepository) CreateCommand(ctx context.Context, invitation *entity.InvitationFriend, tx *sqlx.Tx) error {
	// Insert the new invitation
	insertQuery := `INSERT INTO invitation_friends(sender_id, receiver_id) VALUES (:sender_id, :receiver_id)`
	if tx != nil {
		_, err := tx.NamedExecContext(ctx, insertQuery, invitation)
		return err
	}
	_, err := repo.db.NamedExecContext(ctx, insertQuery, invitation)
	return err
}

func (repo *InvitationFriendRepository) GetByReceiverIdQuery(ctx context.Context, receiverId int64, tx *sqlx.Tx) ([]*entity.InvitationFriend, error) {
	var invitationFriend []*entity.InvitationFriend
	query := "SELECT * FROM invitation_friends WHERE receiver_id = ?"
	if tx != nil {
		err := tx.SelectContext(ctx, &invitationFriend, query, receiverId)
		return invitationFriend, err
	}
	err := repo.db.SelectContext(ctx, &invitationFriend, query, receiverId)
	return invitationFriend, err
}

func (repo *InvitationFriendRepository) GetBySenderIdQuery(ctx context.Context, senderId int64, tx *sqlx.Tx) ([]*entity.InvitationFriend, error) {
	var invitationFriend []*entity.InvitationFriend
	query := "SELECT * FROM invitation_friends WHERE sender_id = ?"
	if tx != nil {
		err := tx.SelectContext(ctx, &invitationFriend, query, senderId)
		return invitationFriend, err
	}
	err := repo.db.SelectContext(ctx, &invitationFriend, query, senderId)
	return invitationFriend, err
}

func (repo *InvitationFriendRepository) GetBySenderAndReceiverIdQuery(ctx context.Context, senderId, receiverId int64, tx *sqlx.Tx) (*entity.InvitationFriend, error) {
	var invitationFriend entity.InvitationFriend
	query := `
		SELECT * 
		FROM invitation_friends 
		WHERE 
			((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))
	`
	if tx != nil {
		err := tx.GetContext(ctx, &invitationFriend, query, senderId, receiverId, receiverId, senderId)
		return &invitationFriend, err
	}
	err := repo.db.GetContext(ctx, &invitationFriend, query, senderId, receiverId, receiverId, senderId)
	return &invitationFriend, err
}

func (repo *InvitationFriendRepository) GetOneByIDQuery(ctx context.Context, id int64, tx *sqlx.Tx) (*entity.InvitationFriend, error) {
	var invitationFriend entity.InvitationFriend
	query := "SELECT * FROM invitation_friends WHERE id = ?"
	if tx != nil {
		err := tx.GetContext(ctx, &invitationFriend, query, id)
		return &invitationFriend, err
	}
	err := repo.db.GetContext(ctx, &invitationFriend, query, id)
	return &invitationFriend, err
}

func (repo *InvitationFriendRepository) DeleteByIDCommand(ctx context.Context, id int64, tx *sqlx.Tx) error {
	query := "DELETE FROM invitation_friends WHERE id = ?"
	if tx != nil {
		_, err := tx.ExecContext(ctx, query, id)
		return err
	}
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}
