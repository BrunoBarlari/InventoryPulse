package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins in development
		// In production, this should be restricted
		return true
	},
}

// Handler handles WebSocket connections
type Handler struct {
	hub *Hub
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

// HandleWebSocket handles WebSocket upgrade requests
func (h *Handler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := NewClient(h.hub, conn)
	h.hub.Register(client)

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()
}

