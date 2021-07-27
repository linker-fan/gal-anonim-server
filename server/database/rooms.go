package database

import (
	"errors"
	"log"
	"time"
)

func InsertRoom(uniqueRoomID string, roomName string, passwordHash string, ownerID int) (int, error) {
	stmt, err := db.Prepare("insert into rooms (id, uniqueRoomID, roomName, passwordHash, ownerID, created, updated) values (default, $1, $2, $3, $4, $5, $6) returning id")
	if err != nil {
		log.Println(err)
		return 0, err
	}

	var roomID int
	err = stmt.QueryRow(uniqueRoomID, roomName, passwordHash, ownerID, time.Now(), time.Now()).Scan(&roomID)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return roomID, nil
}

func CheckIfUniqueRoomIDExists(uniqueRoomID string) error {
	var id string
	err := db.QueryRow("select uniqueRoomID from rooms where uniqueRoomID = $1", uniqueRoomID).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func ChceckIfUserIsOwnerOfTheRoom(uniqueRoomID string, userID int) error {
	var ownerID int
	err := db.QueryRow("select ownerID from rooms where uniqueRoomID=$1", uniqueRoomID).Scan(&ownerID)
	if err != nil {
		log.Println(err)
		return err
	}

	if ownerID != userID {
		return errors.New("Not the owner")
	}

	return nil

}

func DeleteRoom(uniqueRoomID string) error {
	stmt, err := db.Prepare("delete from rooms where uniqueRoomID=$1")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(uniqueRoomID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func UpdateRoomName(name, uniqueRoomID string) error {
	stmt, err := db.Prepare("update rooms set roomName=$1 where uniqueRoomID=$2")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(name, uniqueRoomID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func UpdateRoomPassword(passowrdHash, uniqueRoomID string) error {
	stmt, err := db.Prepare("update rooms set passwordHash=$1 where uniqueRoomID=$2")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(passowrdHash, uniqueRoomID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func UpdateRoom(name, passwordHash, uniqueRoomID string) error {
	stmt, err := db.Prepare("update rooms set roomName=$1, passwordHash=$2, updated=$3 where uniqueRoomID=$4")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(name, passwordHash, time.Now(), uniqueRoomID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetRoomIDByUniqueRoomID(uniqueRoomID string) (int, error) {
	var id int
	err := db.QueryRow("select id from rooms where uniqueRoomID=$1", uniqueRoomID).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}
