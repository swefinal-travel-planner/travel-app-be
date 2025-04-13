package model

type InvitationFriendRequest struct {
	ReceiverEmail string `json:"receiverEmail" binding:"required"`
}
type InvitationFriendReceivedResponse struct {
	Id             int64   `json:"id" binding:"required"`
	SenderUsername string  `json:"senderUsername" binding:"required"`
	SenderImageURL *string `json:"senderImageURL" binding:"required"`
}

type InvitationFriendRequestedResponse struct {
	Id               int64   `json:"id" binding:"required"`
	ReceiverUsername string  `json:"receiverUsername" binding:"required"`
	ReceiverImageURL *string `json:"receiverImageURL" binding:"required"`
}
