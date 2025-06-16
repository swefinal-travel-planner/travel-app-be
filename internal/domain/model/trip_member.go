package model

type tripMemberRole struct {
	Administrator string
	Staff         string
	NormalUser    string
}

var TripMemberRole = tripMemberRole{
	Administrator: "administrator",
	Staff:         "staff",
	NormalUser:    "normal_user",
}

type TripMemberRequest struct {
	TripID int64  `json:"tripID" binding:"required"`
	UserID string `json:"userID" binding:"required"`
	Role   int64  `json:"role" binding:"required"`
}

type TripMemberResponse struct {
	ID       int64   `json:"id"`
	TripID   int64   `json:"trip_id"`
	UserID   int64   `json:"user_id"`
	Role     string  `json:"role"`
	Name     string  `json:"name"`
	PhotoURL *string `json:"photo_url"`
}
