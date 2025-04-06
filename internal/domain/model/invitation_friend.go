package model

type InvitationFriendRequest struct {
	ReceiverID int64 `json:"receiverID" binding:"required"`
}

type InvitationFriendResponse struct {
	Id               int64   `json:"id" binding:"required"`
	ReceiverUsername string  `json:"receiverUsername" binding:"required"`
	ReceiverImageURL *string `json:"receiverImageURL" binding:"required"`
}
