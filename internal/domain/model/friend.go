package model

type FriendResponse struct {
	Id       int64   `json:"id" binding:"required"`
	Username string  `json:"username" binding:"required"`
	ImageURL *string `json:"imageURL" binding:"required"`
}
