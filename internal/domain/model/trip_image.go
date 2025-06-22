package model

import "time"

type TripImageRequest struct {
	ImageURL string `json:"imageUrl" binding:"required"`
}

type TripImageResponse struct {
	ID        int64     `json:"id"`
	TripID    int64     `json:"tripId"`
	ImageURL  string    `json:"imageUrl"`
	CreatedAt time.Time `json:"createdAt"`
}
