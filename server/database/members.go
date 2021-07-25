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

func CheckIfUserIsAMemberOfASpecificRoom(uniqueRoomID string, userID int) error {
	var roomid int
	var id int
	err := db.QueryRow("select distinct m.userid, r.id from members as m join rooms as r on r.id = m.roomid where r.uniqueroomid=$1 and m.userid=$2", uniqueRoomID, userID).Scan(&id, roomid)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetRoomMembers(uniqueRoomID string) ([]string, error) {
	var usernames []string
	rows, err := db.Query("select u.username from rooms as r join members as m on m.roomid = r.id join users as u on m.userid = u.id where r.uniqueroomid=$1", uniqueRoomID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		usernames = append(usernames, username)
	}

	return usernames, nil
}
