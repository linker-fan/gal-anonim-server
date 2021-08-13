package hub

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	broadcast  chan []byte
	rooms      map[*Room]bool
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		rooms:      make(map[*Room]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)
		case client := <-h.Unregister:
			h.unregisterClient(client)
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
	h.Clients[client] = true
}

func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.Clients[client]; ok {
		delete(h.Clients, client)
	}
}

func (h *Hub) CreateRoom(id string) {
	room := NewRoom(id)
	go room.Run()
	h.rooms[room] = true
}

func (h *Hub) FindRoomByID(id string) *Room {
	for r := range h.rooms {
		if r.GetID() == id {
			return r
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

func (h)
