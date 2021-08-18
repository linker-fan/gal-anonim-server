package database

import (
	"database/sql"
	"linker-fan/gal-anonim-server/server/config"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var c *config.Config

func init() {
	conf, err := config.NewConfig("./../config.yml")
	if err != nil {
		log.Fatal(err)
	}

	c = conf
}

func NewMock() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	return db, mock, nil

}

func TestConnectToPostgres(t *testing.T) {
	_, err := connectToPostgres(c)
	if err != nil {
		t.Fail()
	}
}

func TestConnectToRedis(t *testing.T) {
	_, err := connectToRedis(c)
	if err != nil {
		t.Fail()
	}
}
