package model

type InvitationFriendRequest struct {
	ReceiverEmail string `json:"receiverEmail" binding:"required"`
}
type InvitationFriendResponse struct {
	Id             int64   `json:"id" binding:"required"`
	SenderUsername string  `json:"senderUsername" binding:"required"`
	SenderImageURL *string `json:"senderImageURL" binding:"required"`
}
