package model

type friendStatus struct {
	Friend     string
	Sent       string
	Received   string
	Stranger   string
	Restricted string
}

var FriendStatus = friendStatus{
	Friend:     "friend",
	Sent:       "sent",
	Received:   "received",
	Stranger:   "stranger",
	Restricted: "restricted",
}

type FriendResponse struct {
	Id            int64   `json:"id"`
	Email         string  `json:"email"`
	Username      string  `json:"username"`
	ImageURL      *string `json:"imageURL"`
	Status        string  `json:"status"`
	TimeRemaining *int64  `json:"timeRemaining"`
}
