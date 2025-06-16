package model

type TripItemRequest struct {
	PlaceID    string `json:"placeID" binding:"required"`
	TripDay    int64  `json:"tripDay" binding:"required,min=1"`
	OrderInDay int64  `json:"orderInDay" binding:"required,min=1"`
	TimeInDate string `json:"timeInDate" binding:"required"`
}

type PlaceInfo struct {
	Address  string   `json:"address"`
	ID       string   `json:"id"`
	Images   []string `json:"images"`
	Location struct {
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
	} `json:"location"`
	Name       string   `json:"name"`
	Properties []string `json:"properties"`
	Type       string   `json:"type"`
}

type TripItemResponse struct {
	ID         int64      `json:"id"`
	TripID     int64      `json:"tripID"`
	PlaceID    string     `json:"placeID"`
	TripDay    int64      `json:"tripDay"`
	OrderInDay int64      `json:"orderInDay"`
	TimeInDate string     `json:"timeInDate"`
	PlaceInfo  *PlaceInfo `json:"placeInfo"`
}

type TripItemFromAIResponse struct {
	TripID     int64  `json:"trip_id"`
	TripDay    int64  `json:"trip_day"`
	OrderInDay int64  `json:"order_in_day"`
	TimeInDay  string `json:"time_in_day"`
	PlaceID    string `json:"place_id"`
}
