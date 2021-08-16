package hub

import (
	"context"
	"encoding/json"
	"linker-fan/gal-anonim-server/server/database"
	"linker-fan/gal-anonim-server/server/models"
	"log"

	"github.com/go-redis/redis/v8"
)

const PubSubGeneralChannel = "general"

type Hub struct {
	Clients     map[*Client]bool
	Register    chan *Client
	Unregister  chan *Client
	broadcast   chan []byte
	rooms       map[*Room]bool
	users       []*models.User
	redisClient *redis.Client
}

func NewHub() (*Hub, error) {
	users, err := database.GetAllUsers()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		rooms:      make(map[*Room]bool),
		users:      users,
	}, nil
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.UnregisterClient(client)
		case message := <-h.broadcast:
			h.broadcastToClients(message)
		}
	}
}

func (h *Hub) broadcastToClients(message []byte) {
	for client := range h.Clients {
		client.send <- message
	}
}

func (h *Hub) registerClient(client *Client) {
	h.notifyClientJoined(client)
	h.listOnlineClients(client)
	h.Clients[client] = true
}

func (h *Hub) UnregisterClient(client *Client) {
	if _, ok := h.Clients[client]; ok {
		delete(h.Clients, client)
		h.notifyClientLeft(client)

	}
}

func (h *Hub) CreateRoom(id string, private bool) *Room {
	room := NewRoom(id, private)
	go room.Run()
	h.rooms[room] = true
	return room
}

func (h *Hub) FindRoomByID(id string) *Room {
	for r := range h.rooms {
		if r.GetID() == id {
			return r
		}
	}

	return nil
}

func (h *Hub) findClientByID(id int) *Client {
	for c := range h.Clients {
		if c.GetID() == id {
			return c
		}
	}

	return nil
}

func (h *Hub) notifyClientJoined(c *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: c,
	}

	h.broadcastToClients(message.encode())
}

func (h *Hub) notifyClientLeft(c *Client) {
	message := &Message{
		Action: UserLeftAction,
		Sender: c,
	}

	h.broadcastToClients(message.encode())
}

func (h *Hub) listOnlineClients(c *Client) {
	/*
		for existingClient := range h.Clients {
			message := &Message{
				Action: UserJoinedAction,
				Sender: existingClient,
			}
			c.send <- message.encode()
		}
	*/

	for _, users := range h.users {
	}
}

func (h *Hub) findRoomByName(name string) *Room {
	var room *Room
	return room
}

func (h *Hub) publishClientleft(c *Client) error {

	msg := &Message{
		Action: UserLeftAction,
		Sender: c,
	}

	err := database.RedisClient.Publish(context.Background(), PubSubGeneralChannel, msg.encode()).Err()
	if err != nil {
		log.Println(err)
		return nil
	}

	return nil
}

func (h *Hub) publishClientJoined(c *Client) error {

	msg := &Message{
		Action: UserJoinedAction,
		Sender: c,
	}

	if err := database.RedisClient.Publish(context.Background(), PubSubGeneralChannel, msg.encode()).Err(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (h *Hub) listenPubSubChannel() {
	pubsub := database.RedisClient.Subscribe(context.Background(), PubSubGeneralChannel)
	channel := pubsub.Channel()

	for msg := range channel {

		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Printf("Error while unmarshaling JSON message: %v\n", err)
			return
		}

		switch message.Action {
		case UserJoinedAction:
			h.handleUserJoined(message)
		case UserLeftAction:
			h.handleUserLeft(message)
		}
	}
}

func (h *Hub) handleUserJoined(message Message) {

}

func (h *Hub) handleUserLeft(message Message) {

}
