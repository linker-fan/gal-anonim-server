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
	err := ConnectToPostgres(c)
	if err != nil {
		t.Fail()
	}
}

func TestConnectToRedis(t *testing.T) {
	redisClient := ConnectToRedis(c)
	redisClient.Close()
}
