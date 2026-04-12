package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/vpa/quanlynhahang-backend/internal/dto"
	"github.com/vpa/quanlynhahang-backend/internal/usecase"
	"github.com/vpa/quanlynhahang-backend/utils"

	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // dev OK, production nên check origin
	},
}

type Handler struct {
	ChatUC *usecase.ChatUseCase
	NotiUC *usecase.NotificationUseCase
}

func HandleWS(hub *Hub, handler *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, role, err := utils.ParseToken(r)

		if err != nil {
			http.Error(w, "unauthorized", 401)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		client := NewClient(conn, hub, userID, role)
		hub.Register <- client

		go client.WritePump()

		go func() {
			defer func() {
				hub.Unregister <- client
				conn.Close()
			}()

			for {
				_, msgBytes, err := conn.ReadMessage()
				if err != nil {
					break
				}

				var msg dto.WSMessage
				if err := json.Unmarshal(msgBytes, &msg); err != nil {
					continue
				}

				handler.dispatchEvent(client, msg)
			}
		}()
	}
}

func (h *Handler) dispatchEvent(c *Client, msg dto.WSMessage) {
	switch msg.Type {

	case "send_message":
		h.ChatUC.SendMessage(c.UserID, msg)

	case "notify":
		h.NotiUC.Notify(msg)
	}
}
