package repository

import (
	"context"

	"github.com/swefinal-travel-planner/travel-app-be/internal/domain/entity"
)

type InvitationFriendRepository interface {
	CreateCommand(ctx context.Context, user *entity.InvitationFriend) error
	GetByReceiverIdQuery(ctx context.Context, receiverId int64) ([]*entity.InvitationFriend, error)
	GetBySenderIdQuery(ctx context.Context, senderId int64) ([]*entity.InvitationFriend, error)
	GetBySenderAndReceiverIdQuery(ctx context.Context, senderId, receiverId int64) (*entity.InvitationFriend, error)
	GetOneByIDQuery(ctx context.Context, id int64) (*entity.InvitationFriend, error)
	DeleteByIDCommand(ctx context.Context, id int64) error
}
