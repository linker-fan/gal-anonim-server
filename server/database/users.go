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

func GetIDAndPasswordByUsername(username string) (int, string, bool, error) {
	var id int
	var passwordHash string
	var isAdmin bool
	err := db.QueryRow("select id, passwordHash, isAdmin from users where username=$1", username).Scan(&id, &passwordHash, &isAdmin)
	if err != nil {
		return 0, "", false, err
	}

	return id, passwordHash, isAdmin, nil
}

func GetUserIDByUsername(username string) (int, error) {
	var id int
	err := db.QueryRow("select id from users where username=$1", username).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, nil
	}

	return id, nil
}

func SetPin(pin string, id int) error {
	stmt, err := db.Prepare("update users set pin=$1 where id=$2")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(pin, id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
