package hub

import (
	"encoding/json"
	"log"
	"strconv"
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
	id       int
	username string
	conn     *websocket.Conn
	hub      *Hub
	send     chan []byte
	rooms    map[*Room]bool
}

func NewClient(conn *websocket.Conn, hub *Hub, username string, id int) *Client {
	return &Client{
		id:       id,
		username: username,
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
		if room := c.hub.FindRoomByID(roomID.GetID()); room != nil {
			room.broadcast <- &message
		}
	case JoinRoomAction:
		c.handleJoinRoomMessage(message)
	case LeaveRoomAction:
		c.handleLeaveRoomMessage(message)
	case JoinPrivateRoomAction:
		c.handleJoinRoomPrivateMessage(message)
	}
}

func (c *Client) handleJoinRoomMessage(m Message) {
	roomName := m.Message
	c.joinRoom(roomName, nil)
}

func (c *Client) handleJoinRoomPrivateMessage(m Message) {
	target := c.hub.findClientByID(m.Sender.GetID())
	if target == nil {
		return
	}

	roomName := m.Message + strconv.Itoa(c.id)
	c.joinRoom(roomName, target)
	target.joinRoom(roomName, c)
}

func (c *Client) handleLeaveRoomMessage(m Message) {
	room := c.hub.FindRoomByID(m.Target.GetID())
	if room == nil {
		return
	}
	if _, ok := c.rooms[room]; ok {
		delete(c.rooms, room)
	}

	room.unregister <- c
}

func (c *Client) GetID() int {
	return c.id
}

func (c *Client) joinRoom(roomName string, sender *Client) {
	room := c.hub.findRoomByName(roomName)
	if room == nil {
		room = c.hub.CreateRoom(roomName, sender != nil)
	}

	if sender == nil && room.Private {
		return
	}

	if !c.isInRoom(room) {
		c.rooms[room] = true
		room.register <- c
		c.notifyRoomJoined(room, sender)
	}
}

func (c *Client) isInRoom(r *Room) bool {
	if _, ok := c.rooms[r]; ok {
		return true
	}
	return false
}

func (c *Client) notifyRoomJoined(r *Room, s *Client) {
	m := Message{
		Action: RoomJoinedAction,
		Target: r,
		Sender: s,
	}

	c.send <- m.encode()
}
