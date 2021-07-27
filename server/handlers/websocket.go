package handlers

import (
	"linker-fan/gal-anonim-server/server/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ChatHub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newChatHub() *ChatHub {
	return &ChatHub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (c *ChatHub) run() {
	for {
		select {
		case client := <-c.register:
			c.clients[client] = true
		case client := <-c.unregister:
			if _, ok := c.clients[client]; ok {
				delete(c.clients, client)
				close(client.send)
			}

		case message := <-c.broadcast:
			for client := range c.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(c.clients, client)
				}
			}

		}
	}
}

type Client struct {
	hub  *ChatHub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) read() {}

func (c *Client) write() {}

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

	hub := newChatHub()
	client := &Client{hub: hub, conn: ws, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.read()
	go client.write()

}
