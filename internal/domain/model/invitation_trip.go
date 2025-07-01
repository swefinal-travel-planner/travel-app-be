package model

import "time"

type InvitationTripRequest struct {
	TripID     int64 `json:"tripId" binding:"required"`
	ReceiverID int64 `json:"receiverId" binding:"required"`
}

type InvitationTripReceivedResponse struct {
	ID         int64     `json:"id"`
	TripID     int64     `json:"tripId"`
	SenderID   int64     `json:"senderId"`
	ReceiverID int64     `json:"receiverId"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type InvitationTripSentResponse struct {
	ID         int64     `json:"id"`
	TripID     int64     `json:"tripId"`
	SenderID   int64     `json:"senderId"`
	ReceiverID int64     `json:"receiverId"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type InvitationTripPendingResponse struct {
	ID         int64     `json:"id"`
	TripID     int64     `json:"tripId"`
	SenderID   int64     `json:"senderId"`
	ReceiverID int64     `json:"receiverId"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
