package model

type TripItemRequest struct {
	PlaceID    string `json:"placeID" binding:"required"`
	TripDay    int64  `json:"tripDay" binding:"required,min=1"`
	OrderInDay int64  `json:"orderInDay" binding:"required,min=1"`
	TimeInDate string `json:"timeInDate" binding:"required"`
}

type TripItemResponse struct {
	ID         int64  `json:"id"`
	TripID     int64  `json:"tripID"`
	PlaceID    string `json:"placeID"`
	TripDay    int64  `json:"tripDay"`
	OrderInDay int64  `json:"orderInDay"`
	TimeInDate string `json:"timeInDate"`
}
