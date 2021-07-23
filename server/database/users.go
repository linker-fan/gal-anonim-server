package database

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

//InsertUser takes username and hashed password from the register handler and creates a new row in users table. by default new user is not an admin
//@author hyperxpizza
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

//CheckIfUsernameExists queries the users table to check if row with username given as an argument to the function already exists in the database
//@author hyperxpizza
func CheckIfUsernameExists(username string) error {
	var u string
	err := db.QueryRow("select username from users where username=$1", username).Scan(&u)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			return err
		}
	}

	return errors.New("Username already taken")
}
