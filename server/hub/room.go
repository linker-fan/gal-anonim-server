package hub

type Room struct {
	Name        string
	Clients     map[*Client]bool
	ActiveUsers map[string]*Client
}

func NewRoom(name string) *Room {
	return &Room{
		Name:    name,
		Clients: make(map[*Client]bool),
	}
}

func (r *Room) Join(c *Client) error {
	if _, ok := r.Clients[c]; ok {
		c.Conn.WriteMessage()
	}

	return nil
}

func (r *Room) Leave(c *Client) error {
	return nil
}

func (r *Room) Boardcast() {

}
