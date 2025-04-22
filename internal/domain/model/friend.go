package model

type FriendStatus string

const (
	Friend     FriendStatus = "friend"
	Sent       FriendStatus = "sent"
	Received   FriendStatus = "received"
	Stranger   FriendStatus = "stranger"
	Restricted FriendStatus = "restricted"
)

type FriendResponse struct {
	Id            int64        `json:"id"`
	Email         string       `json:"email"`
	Username      string       `json:"username"`
	ImageURL      *string      `json:"imageURL"`
	Status        FriendStatus `json:"status"`
	TimeRemaining *int64       `json:"timeRemaining"`
}
