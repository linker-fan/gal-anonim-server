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

	err = database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	router.Run(c.Server.Port, c.Server.Mode)
}
