package database

import (
	"log"
	"time"
)

func InsertRoom(uniqueRoomID string, roomName string, passwordHash string, ownerID int) error {
	stmt, err := db.Prepare("insert into rooms (id, uniqueRoomID, roomName, passwordHash, ownerID, created, updated) values (default, $1, $2, $3, $4, $5, $6)")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(uniqueRoomID, roomName, passwordHash, ownerID, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
