package main

import (
	"linker-fan/gal-anonim-server/server/config"
	"linker-fan/gal-anonim-server/server/database"
	"linker-fan/gal-anonim-server/server/router"

	"log"
)

func main() {
	c, err := config.NewConfig("./config.yml")
	if err != nil {
		log.Fatal(err)
	}

	err = database.ConnectToPostgres(c)
	if err != nil {
		log.Fatal(err)
	}

	err = database.ConnectToRedis(c)
	if err != nil {
		log.Fatal(err)
	}

	//api
	router.Run(c.Server.Port, c.Server.Mode)
}
