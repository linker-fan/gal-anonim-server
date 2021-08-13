package hub

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second
	// Max time till next pong from peer
	pongWait = 60 * time.Second
	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	ID       int
	Username string
	conn     *websocket.Conn
	hub      *Hub
	send     chan []byte
	rooms    map[*Room]bool
}

func NewClient(conn *websocket.Conn, hub *Hub, username string, id int) *Client {
	return &Client{
		ID:       id,
		Username: username,
		conn:     conn,
		hub:      hub,
		send:     make(chan []byte),
		rooms:    make(map[*Room]bool),
	}
}

func (c *Client) disconnect() {
	c.hub.Unregister <- c
	for room := range c.rooms {
		room.unregister <- c
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.disconnect()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, jsonMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		c.handleNewMessage(jsonMessage)
	}
}

func (c *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %v\n", err)
	}

	message.Sender = c

	switch message.Action {
	case SendMessageAction:
		roomID := message.Target
		if room := c.hub.FindRoomByID(roomID); room != nil {
			room.broadcast <- &message
		}
	case JoinRoomAction:
		c.handleJoinRoomMessage(message)
	case LeaveRoomAction:
		c.handleLeaveRoomMessage(message)
	}
}

func (c *Client) handleJoinRoomMessage(m Message) {
	roomID := m.Message
	room := c.hub.FindRoomByID(roomID)
	if room == nil {
		//room = c.hub.CreateRoom(roomID)
	} else {
		c.rooms[room] = true
		room.register <- c
	}
}

func (c *Client) handleLeaveRoomMessage(m Message) {
	room := c.hub.FindRoomByID(m.Target)
	if _, ok := c.rooms[room]; ok {
		delete(c.rooms, room)
	}

	room.unregister <- c
}
