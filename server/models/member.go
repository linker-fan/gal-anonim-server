package models

import "time"

type Member struct {
	ID     int       `json:"id"`
	RoomID int       `json:"room_id"`
	UserID int       `json:"user_id"`
	Joined time.Time `json:"joined"`
}
