package handlers

import (
	"linker-fan/gal-anonim-server/server/hub"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func ChatWebsocket(c *gin.Context) {

}

func serveWS(wsServer *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := hub.NewClient(conn, wsServer)

	go client.WritePump()
	go client.ReadPump()

	wsServer.Register <- client
}
