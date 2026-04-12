package websocket

import (
	"encoding/json"
	"sync"

	"github.com/vpa/quanlynhahang-backend/internal/dto"
)

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true

		case client := <-h.Unregister:
			delete(h.Clients, client)
			close(client.Send)
		}
	}
}

func (h *Hub) SendToUser(userID uint, msg dto.WSMessage) {
	data, _ := json.Marshal(msg)

	h.mu.RLock()
	defer h.mu.RUnlock()

	for c := range h.Clients {
		if c.UserID == userID {
			c.Send <- data
		}
	}
}

func (h *Hub) BroadcastToRoom(roomID uint, msg dto.WSMessage) {
	data, _ := json.Marshal(msg)

	for c := range h.Clients {

		c.Send <- data
	}
}

func (h *Hub) SendToRole(role string, msg dto.WSMessage) {
	data, _ := json.Marshal(msg)

	for c := range h.Clients {
		if c.Role == role { // 👈 cần thêm Role vào Client
			c.Send <- data
		}
	}
}
