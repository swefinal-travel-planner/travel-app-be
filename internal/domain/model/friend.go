package model

type FriendResponse struct {
	Username string  `json:"username" binding:"required"`
	ImageURL *string `json:"imageURL" binding:"required"`
}
