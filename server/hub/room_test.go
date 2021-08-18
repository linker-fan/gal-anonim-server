package hub

import (
	"linker-fan/gal-anonim-server/server/config"
	"linker-fan/gal-anonim-server/server/database"
	"testing"
)

func TestNewRoom(t *testing.T) {
	conf, err := config.NewConfig("./../config.yml")
	if err != nil {
		t.Fail()
	}

	err = database.ConnectToPostgres(conf)
	if err != nil {
		t.Fail()
	}

	err = database.ConnectToRedis(conf)
	if err != nil {
		t.Fail()
	}

	room := NewRoom("some-unique-room-id", false)
	room.Run()
}
