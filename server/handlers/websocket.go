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

func (a *API) ChatWebsocket(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.Status(http.StatusInternalServerError)
	}

	id, exists := c.Get("id")
	if !exists {
		c.Status(http.StatusInternalServerError)
	}

	a.serveWS(c.Writer, c.Request, username.(string), id.(int))
}

func (a *API) serveWS(w http.ResponseWriter, r *http.Request, username string, id int) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := hub.NewClient(conn, a.wsServer, username, id)

	go client.WritePump()
	go client.ReadPump()

	a.wsServer.Register <- client
}
