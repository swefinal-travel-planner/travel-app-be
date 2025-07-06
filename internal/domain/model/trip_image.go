package model

import "time"

type TripImageRequest struct {
	ImageURL string  `json:"imageUrl" binding:"required"`
	PlaceID  *string `json:"placeId"`
}

type TripImageResponse struct {
	ID        int64     `json:"id"`
	TripID    int64     `json:"tripId"`
	PlaceID   *string   `json:"placeId"`
	ImageURL  string    `json:"imageUrl"`
	CreatedAt time.Time `json:"createdAt"`
}

type TripImageWithUserInfoResponse struct {
	ID        int64     `json:"id"`
	TripID    int64     `json:"tripId"`
	PlaceID   *string   `json:"placeId"`
	ImageURL  string    `json:"imageUrl"`
	Author    UserInfo  `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserInfo struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	PhotoURL *string `json:"photoURL"`
}
