package model

type InvitationFriendRequest struct {
	ReceiverEmail string `json:"receiverEmail" binding:"required"`
}
type InvitationFriendReceivedResponse struct {
	Id             int64   `json:"id"`
	SenderUsername string  `json:"senderUsername"`
	SenderImageURL *string `json:"senderImageURL"`
}

type InvitationFriendRequestedResponse struct {
	Id               int64   `json:"id"`
	ReceiverUsername string  `json:"receiverUsername"`
	ReceiverImageURL *string `json:"receiverImageURL"`
}
