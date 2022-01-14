package chat

type ChatHub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *ChatHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			_, ok := h.clients[client]
			if ok {
				delete(h.clients, client)
				close(client.channel)
			}
		}
	}
}

func (h *ChatHub) Register(c *Client) {
	h.register <- c
}

func (h *ChatHub) DeRegister(c *Client) {
	h.unregister <- c
}

func (h *ChatHub) GetClients() []*Client {
	var clients []*Client
	for client, _ := range h.clients {
		clients = append(clients, client)
	}

	return clients
}
