package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"gochat/internal/domain"
)

type Hub struct {
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *RoomMessage
	mu         sync.RWMutex
}

type RoomMessage struct {
	RoomID  string
	Message *domain.Message
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *RoomMessage, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.roomID] == nil {
				h.rooms[client.roomID] = make(map[*Client]bool)
			}
			h.rooms[client.roomID][client] = true
			h.mu.Unlock()
			log.Printf("Client registered to room: %s (total: %d)", client.roomID, len(h.rooms[client.roomID]))

		case client := <-h.unregister:
			h.mu.Lock()
			if room, ok := h.rooms[client.roomID]; ok {
				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.send)
					if len(room) == 0 {
						delete(h.rooms, client.roomID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client unregistered from room: %s", client.roomID)

		case message := <-h.broadcast:
			h.mu.RLock()
			room, exists := h.rooms[message.RoomID]
			if !exists {
				h.mu.RUnlock()
				log.Printf("Room %s does not exist for broadcast", message.RoomID)
				continue
			}

			data, err := json.Marshal(message.Message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				h.mu.RUnlock()
				continue
			}

			clientsToRemove := make([]*Client, 0)
			for client := range room {
				select {
				case client.send <- data:
				default:
					clientsToRemove = append(clientsToRemove, client)
				}
			}
			h.mu.RUnlock()

			if len(clientsToRemove) > 0 {
				h.mu.Lock()
				for _, client := range clientsToRemove {
					if room, ok := h.rooms[message.RoomID]; ok {
						if _, ok := room[client]; ok {
							close(client.send)
							delete(room, client)
							if len(room) == 0 {
								delete(h.rooms, message.RoomID)
							}
						}
					}
				}
				h.mu.Unlock()
			}

			log.Printf("Broadcasted message to room %s (%d clients)", message.RoomID, len(room)-len(clientsToRemove))
		}
	}
}

func (h *Hub) BroadcastMessage(roomID string, message *domain.Message) {
	h.broadcast <- &RoomMessage{
		RoomID:  roomID,
		Message: message,
	}
}
