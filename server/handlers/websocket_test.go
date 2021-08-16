package handlers

import (
	"linker-fan/gal-anonim-server/server/config"
	"linker-fan/gal-anonim-server/server/database"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestServeWS(t *testing.T) {

	config, err := config.NewConfig("./../config.yml")
	if err != nil {
		t.Fatalf("Loading config failed: %v\n", err)
	}

	err = database.ConnectToPostgres(config)
	if err != nil {
		t.Fatalf("Connecting to postgres failed: %v\n", err)
	}

	err = database.ConnectToRedis(config)
	if err != nil {
		t.Fatalf("Connecting to Redis failed: %v\n", err)
	}

	//init test server
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveWS(w, r, "test-username", 1)
	}))

	defer s.Close()

	url := "ws" + strings.TrimPrefix(s.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	defer ws.Close()
}
