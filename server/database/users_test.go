package database

import (
	"fmt"
	"linker-fan/gal-anonim-server/server/config"
	"log"
	"testing"
)

func init() {
	c, err := config.NewConfig("./../config.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = ConnectToPostgres(c)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAllUsers(t *testing.T) {
	users, err := GetAllUsers()
	if err != nil {
		t.Fail()
	}

	fmt.Println(users)
}
