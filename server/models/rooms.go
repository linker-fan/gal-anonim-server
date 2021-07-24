package models

import "time"

type Room struct {
	ID           int       `json:"id"`
	UniqueRoomID string    `json:"unique_room_id"`
	RoomName     string    `json:"room_name"`
	PasswordHash string    `json:"password_hash"`
	OwnerID      int       `json:"owner_id"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}
