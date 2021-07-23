package database

import (
	"log"
	"time"
)

func InsertUser(username, passwordHash string) error {
	stmt, err := db.Prepare("insert into users(id,username,passwordHash,isAdmin,created,updated) values (default, $1, $2, $3, $4, $5)")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(username, passwordHash, false, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
