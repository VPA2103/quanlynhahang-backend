package realtime

import "encoding/json"

func PushToUser(userID string, data any) {
	if HubInstance == nil {
		return
	}

	msg, _ := json.Marshal(data)

	if clients, ok := HubInstance.Clients[userID]; ok {
		for client := range clients {
			client.Send <- msg
		}
	}
}
