package database

import (
	"database/sql"
	"errors"
	"linker-fan/gal-anonim-server/server/models"
	"log"
	"time"
)

//InsertUser takes username and hashed password from the register handler and creates a new row in users table. by default new user is not an admin
//@author hyperxpizza
func (d *DatabaseWrapper) InsertUser(username, passwordHash string) error {
	stmt, err := d.db.Prepare("insert into users(id,username,passwordHash,isAdmin,created,updated) values (default, $1, $2, $3, $4, $5)")
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
func (d *DatabaseWrapper) CheckIfUsernameExists(username string) error {
	var u string
	err := d.db.QueryRow("select username from users where username=$1", username).Scan(&u)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			return err
		}
	}

	return errors.New("Username already taken")
}

//GetIDAndPasswordByUsername queries id, passwordHash and isAdmin with username given as an argument
func (d *DatabaseWrapper) GetIDAndPasswordByUsername(username string) (int, string, bool, error) {
	var id int
	var passwordHash string
	var isAdmin bool
	err := d.db.QueryRow("select id, passwordHash, isAdmin from users where username=$1", username).Scan(&id, &passwordHash, &isAdmin)
	if err != nil {
		return 0, "", false, err
	}

	return id, passwordHash, isAdmin, nil
}

func (d *DatabaseWrapper) GetUserIDByUsername(username string) (int, error) {
	var id int
	err := d.db.QueryRow("select id from users where username=$1", username).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}

func (d *DatabaseWrapper) SetPin(pin string, id int) error {
	stmt, err := d.db.Prepare("update users set pin=$1 where id=$2")
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

func (d *DatabaseWrapper) GetAllUsers() ([]*models.User, error) {
	var users []*models.User

	rows, err := d.db.Query("select id, username, isAdmin, created, updated from users")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.IsAdmin, &user.Created, &user.Updated)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (d *DatabaseWrapper) DeleteUser(id int) error {
	stmt, err := d.db.Prepare("delete from users where id=$1")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (d *DatabaseWrapper) CheckIfEmailExists(email string) error {
	var e string
	err := d.db.QueryRow("select username from users where username=$1", email).Scan(&e)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			return err
		}
	}

	return errors.New("Email already taken")
}
