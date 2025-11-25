package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Event types
const (
	EventProductCreated = "product.created"
	EventProductUpdated = "product.updated"
	EventProductDeleted = "product.deleted"
	EventStockUpdated   = "stock.updated"
	EventCategoryCreated = "category.created"
	EventCategoryUpdated = "category.updated"
	EventCategoryDeleted = "category.deleted"
)

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("WebSocket client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastMessage sends a message to all connected clients
func (h *Hub) BroadcastMessage(eventType string, payload interface{}) {
	msg := Message{
		Type:    eventType,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return
	}

	h.broadcast <- data
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

