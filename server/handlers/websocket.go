package handlers

import (
	"linker-fan/gal-anonim-server/server/database"
	"linker-fan/gal-anonim-server/server/hub"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var chatHub *hub.ChatHub

func init() {
	chatHub = hub.NewChatHub()
	chatHub.Run()
}

func ChatWebsocket(c *gin.Context) {

	id, exists := c.Get("id")
	if !exists {
		c.Status(http.StatusInternalServerError)
		return
	}

	uniqueRoomID := c.Param("uniqueRoomID")
	if uniqueRoomID == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	err := database.CheckIfUserIsAMemberOfASpecificRoom(uniqueRoomID, id.(int))
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	ws, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	client := &hub.Client{Hub: chatHub, Conn: ws, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	go client.Read()
	go client.Write()

}
