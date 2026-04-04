package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client.
type Client struct {
	ID       string
	PlayerID string
	RoomCode string
	Conn     *websocket.Conn
	Send     chan []byte
	Hub      *Hub
}

// Hub manages WebSocket connections and room routing.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]*Client        // clientID -> Client
	rooms   map[string]map[string]bool // roomCode -> set of clientIDs
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		rooms:   make(map[string]map[string]bool),
	}
}

// Register adds a client to the hub and its room.
func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[c.ID] = c
	if c.RoomCode != "" {
		if h.rooms[c.RoomCode] == nil {
			h.rooms[c.RoomCode] = make(map[string]bool)
		}
		h.rooms[c.RoomCode][c.ID] = true
	}
}

// Unregister removes a client from the hub and its room.
func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[c.ID]; ok {
		delete(h.clients, c.ID)
		close(c.Send)
	}
	if c.RoomCode != "" {
		if room, ok := h.rooms[c.RoomCode]; ok {
			delete(room, c.ID)
			if len(room) == 0 {
				delete(h.rooms, c.RoomCode)
			}
		}
	}
}

// BroadcastToRoom sends a message to all clients in a room.
func (h *Hub) BroadcastToRoom(roomCode string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, ok := h.rooms[roomCode]
	if !ok {
		return
	}
	for clientID := range room {
		if client, exists := h.clients[clientID]; exists {
			select {
			case client.Send <- message:
			default:
				log.Printf("hub: send buffer full for client %s, dropping message", clientID)
			}
		}
	}
}

// SendTo sends a message to a specific client.
func (h *Hub) SendTo(clientID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.clients[clientID]; ok {
		select {
		case client.Send <- message:
		default:
			log.Printf("hub: send buffer full for client %s", clientID)
		}
	}
}

// RoomSize returns the number of clients in a room.
func (h *Hub) RoomSize(roomCode string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, ok := h.rooms[roomCode]; ok {
		return len(room)
	}
	return 0
}

// GetClientsByRoom returns all client IDs in a room.
func (h *Hub) GetClientsByRoom(roomCode string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, ok := h.rooms[roomCode]
	if !ok {
		return nil
	}
	ids := make([]string, 0, len(room))
	for id := range room {
		ids = append(ids, id)
	}
	return ids
}

// Shutdown gracefully closes all client connections and channels.
// R5-H3: prevents WebSocket goroutine leaks on server shutdown.
func (h *Hub) Shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for id, client := range h.clients {
		close(client.Send)
		if client.Conn != nil {
			if err := client.Conn.Close(); err != nil {
				log.Printf("hub: error closing connection for client %s: %v", id, err)
			}
		}
		delete(h.clients, id)
	}

	// Clear all rooms
	for code := range h.rooms {
		delete(h.rooms, code)
	}

	log.Println("hub: all WebSocket connections shut down")
}

// JoinRoom moves a client to a room.
func (h *Hub) JoinRoom(clientID, roomCode string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	client, ok := h.clients[clientID]
	if !ok {
		return
	}

	// Leave old room
	if client.RoomCode != "" {
		if room, exists := h.rooms[client.RoomCode]; exists {
			delete(room, clientID)
			if len(room) == 0 {
				delete(h.rooms, client.RoomCode)
			}
		}
	}

	// Join new room
	client.RoomCode = roomCode
	if h.rooms[roomCode] == nil {
		h.rooms[roomCode] = make(map[string]bool)
	}
	h.rooms[roomCode][clientID] = true
}
