package hub

type Room struct {
	id         string
	name       string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

func NewRoom(id string, name string) *Room {
	r := &Room{
		id:         id,
		name:       name,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}

	return r
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.registerClientInRoom(client)
		case client := <-r.unregister:
			r.unregisterClientInRoom(client)
		case message := <-r.broadcast:
			r.broadcastToClientsInRoom(message.encode())
		}
	}
}

func (r *Room) registerClientInRoom(client *Client) {
	r.notifyClientJoined(client)
	r.clients[client] = true
}

func (r *Room) unregisterClientInRoom(client *Client) {
	if _, ok := r.clients[client]; ok {
		delete(r.clients, client)
	}
}

func (r *Room) broadcastToClientsInRoom(message []byte) {
	for client := range r.clients {
		client.send <- message
	}
}

func (r *Room) notifyClientJoined(c *Client) {

}

func (r *Room) GetID() string {
	return r.id
}
