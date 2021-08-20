package handlers

import (
	"linker-fan/gal-anonim-server/server/config"
	"linker-fan/gal-anonim-server/server/database"
	"linker-fan/gal-anonim-server/server/hub"
)

type API struct {
	dw       *database.DatabaseWrapper
	wsServer *hub.Hub
}

func NewAPIWrapper(c *config.Config) (*API, error) {
	dw, err := database.NewDatabaseWrapper(c)
	if err != nil {
		return nil, err
	}

	hub, err := hub.NewHub(dw)
	if err != nil {
		return nil, err
	}

	api := API{
		dw:       dw,
		wsServer: hub,
	}

	go api.wsServer.Run()

	return &api, nil
}
