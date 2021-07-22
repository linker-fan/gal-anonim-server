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

func TestConnect(t *testing.T) {
	err := Connect(c)
	if err != nil {
		t.Fail()
	}
}
