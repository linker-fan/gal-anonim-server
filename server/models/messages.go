package models

import "time"

type Message struct {
	ID      int       `json:"id"`
	RoomID  int       `json:"room_id"`
	UserID  int       `json:"user_id"`
	Text    string    `json:"text"`
	Created time.Time `json:"created"`
}
