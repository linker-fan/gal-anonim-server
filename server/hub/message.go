package hub

import (
	"linker-fan/gal-anonim-server/server/database"
	"time"
)

type Message struct {
	Room     string
	Username string
	Text     string
	Time     time.Time
}

func NewMessage(room, username, text string) *Message {
	return &Message{
		Room:     room,
		Username: username,
		Text:     text,
		Time:     time.Now(),
	}
}

func (m *Message) MapIntoTheDatabase() error {
	roomID, err := database.GetRoomIDByName(m.Room)
	if err != nil {
		return err
	}

	userID, err := database.GetUserIDByUsername(m.Username)
	if err != nil {
		return err
	}

	err = database.InsertMessage(roomID, userID, m.Text, m.Time)
	if err != nil {
		return err
	}

	return nil
}
