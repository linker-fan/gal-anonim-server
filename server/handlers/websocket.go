package handlers

import (
	"linker-fan/gal-anonim-server/server/hub"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

var wsServer *hub.Hub

func RunWsServer() {
	wsServer, err := hub.NewHub()
	if err != nil {
		log.Fatalf("hub.NewHub error: %v\n", err)
		return
	}
	go wsServer.Run()
}

func ChatWebsocket(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.Status(http.StatusInternalServerError)
	}

	id, exists := c.Get("id")
	if !exists {
		c.Status(http.StatusInternalServerError)
	}

	serveWS(c.Writer, c.Request, username.(string), id.(int))
}

func serveWS(w http.ResponseWriter, r *http.Request, username string, id int) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := hub.NewClient(conn, wsServer, username, id)

	go client.WritePump()
	go client.ReadPump()

	wsServer.Register <- client
}
