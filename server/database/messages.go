package database

import (
	"log"
	"time"
)

func (d *DatabaseWrapper) InsertMessage(roomID, userID int, text string, created time.Time) error {
	stmt, err := d.db.Prepare("insert into messages(id, roomID, userID, messageText, created) values(default, $1, $2, $3, $4)")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(roomID, userID, text, created)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
