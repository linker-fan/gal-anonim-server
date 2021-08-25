package handlers

import (
	"linker-fan/gal-anonim-server/server/config"
	"linker-fan/gal-anonim-server/server/database"
	"linker-fan/gal-anonim-server/server/filestorage"
	"linker-fan/gal-anonim-server/server/hub"
)

type API struct {
	dw       *database.DatabaseWrapper
	wsServer *hub.Hub
	fs       *filestorage.FileStorage
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

	fs, err := filestorage.NewFileStorage(c.FileStorage.Endpoint, c.FileStorage.AccessKeyID, c.FileStorage.SecretAccessKey, true)
	if err != nil {
		return nil, err
	}

	api := API{
		dw:       dw,
		wsServer: hub,
		fs:       fs,
	}

	go api.wsServer.Run()

	return &api, nil
}
