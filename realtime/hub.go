package realtime

type Hub struct {
	Clients    map[string]map[*Client]bool // userID -> clients
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

var HubInstance = NewHub() // ✅ GLOBAL HUB

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if h.Clients[client.UserID] == nil {
				h.Clients[client.UserID] = make(map[*Client]bool)
			}
			h.Clients[client.UserID][client] = true

		case client := <-h.Unregister:
			if clients, ok := h.Clients[client.UserID]; ok {
				delete(clients, client)
				close(client.Send)
				if len(clients) == 0 {
					delete(h.Clients, client.UserID)
				}
			}
		}
	}
}
