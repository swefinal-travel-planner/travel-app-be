package model

type TestNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	PushToken string `json:"pushToken"`
}