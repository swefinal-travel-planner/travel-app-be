package model

type FriendResponse struct {
	Id       int64   `json:"id"`
	Username string  `json:"username"`
	ImageURL *string `json:"imageURL"`
}
