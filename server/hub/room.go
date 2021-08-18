package hub

import (
	"context"
	"fmt"
	"linker-fan/gal-anonim-server/server/database"
	"log"
)

const welcomeMessage = "%s joined the chat"

type Room struct {
	id         string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	private    bool
}

func NewRoom(id string, private bool) *Room {
	r := &Room{
		id:         id,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
		private:    private,
	}

	return r
}

func (r *Room) Run() {
	log.Printf("Running room with ID: %s\n", r.id)

	go r.subscribeToRoomMessages()

	for {
		select {
		case client := <-r.register:
			if err := r.registerClientInRoom(client); err != nil {
				log.Println(err)
			}
		case client := <-r.unregister:
			if err := r.unregisterClientInRoom(client); err != nil {
				log.Println(err)
			}
		case message := <-r.broadcast:
			r.publishRoomMessage(message.encode())
		}
	}
}

func (r *Room) registerClientInRoom(client *Client) error {
	if !r.GetPrivate() {
		r.notifyClientJoined(client)
	}
	r.clients[client] = true
	err := database.InsertMewmberWithUniqueRoomID(r.id, client.id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Room) unregisterClientInRoom(client *Client) error {
	if _, ok := r.clients[client]; ok {
		delete(r.clients, client)
		err := database.DeleteMemberWithUnqueRoomID(r.id, client.id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Room) broadcastToClientsInRoom(message []byte) {
	for client := range r.clients {
		client.send <- message
	}
}

func (r *Room) notifyClientJoined(c *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  r,
		Message: fmt.Sprintf(welcomeMessage, c.username),
	}

	r.publishRoomMessage(message.encode())
}

func (r *Room) publishRoomMessage(message []byte) error {
	err := database.RedisClient.Publish(context.Background(), r.GetID(), message).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *Room) subscribeToRoomMessages() {
	pubsub := database.RedisClient.Subscribe(context.Background(), r.GetID())
	channel := pubsub.Channel()

	for msg := range channel {
		r.broadcastToClientsInRoom([]byte(msg.Payload))
	}
}

func (r *Room) GetID() string {
	return r.id
}

func (r *Room) GetPrivate() bool {
	return r.private
}
