package model

type TripItemRequest struct {
	TripID     int64  `json:"tripID" binding:"required"`
	PlaceID    string `json:"placeID" binding:"required"`
	TripDay    int64  `json:"tripDay" binding:"required,min=1"`
	OrderInDay int64  `json:"orderInDay" binding:"required,min=1"`
	Tag        string `json:"tag" binding:"required"`
}

type TripItemResponse struct {
	TripID     int64  `json:"tripID"`
	PlaceID    string `json:"placeID"`
	TripDay    int64  `json:"tripDay"`
	OrderInDay int64  `json:"orderInDay"`
	Tag        string `json:"tag"`
}
