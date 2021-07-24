package database

import (
	"log"
	"time"
)

func InsertMember(roomID, userID int) error {
	stmt, err := db.Prepare("insert into members (id, roomID, userID, joined) values (default, $1, $2, $3)")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(roomID, userID, time.Now())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func CheckIfUserIsAMemberOfASpecificRoom(uniqueRoomID string, userID string) error {
	err := db.QueryRow("select distinct m.userid, r.id from members as m join rooms as r on r.id = m.roomid where r.uniqueroomid=$1 and m.userid=$2;
	")
}
