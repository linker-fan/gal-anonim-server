package database

import (
	"linker-fan/gal-anonim-server/server/config"
	"log"
	"testing"
)

var c *config.Config

func init() {
	conf, err := config.NewConfig("./../config.yml")
	if err != nil {
		log.Fatal(err)
	}

	c = conf
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
