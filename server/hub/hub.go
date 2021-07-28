package hub

type ChatHub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (c *ChatHub) Run() {
	for {
		select {
		case client := <-c.Register:
			c.Clients[client] = true
		case client := <-c.Unregister:
			if _, ok := c.Clients[client]; ok {
				delete(c.Clients, client)
				close(client.Send)
			}

		case message := <-c.Broadcast:
			for client := range c.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(c.Clients, client)
				}
			}

		}
	}
}
