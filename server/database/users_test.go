package database

import (
	"linker-fan/gal-anonim-server/server/config"
	"log"
)

func init() {
	c, err := config.NewConfig("./../config.yml")
	if err != nil {
		log.Fatal(err)
	}
	_, err = connectToPostgres(c)
	if err != nil {
		log.Fatal(err)
	}
}
