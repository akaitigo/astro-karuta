package handler

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/akaitigo/astro-karuta/backend/internal/ws"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMsgSize = 4096
	// maxWSConnections limits the total number of concurrent WebSocket connections.
	maxWSConnections = 200
)

// allowedOrigins returns the set of allowed WebSocket origins from the
// CORS_ORIGIN environment variable. Multiple origins can be separated by
// commas. When the variable is empty the default development origin is used.
func allowedOrigins() map[string]bool {
	raw := os.Getenv("CORS_ORIGIN")
	if raw == "" {
		raw = "http://localhost:3000"
	}
	origins := make(map[string]bool)
	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			origins[o] = true
		}
	}
	return origins
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Allow connections without an Origin header (e.g. non-browser clients).
			return true
		}
		return allowedOrigins()[origin]
	},
}

// WSHandler handles WebSocket connections.
type WSHandler struct {
	hub       *ws.Hub
	gm        *ws.GameManager
	connCount atomic.Int64
}

// NewWSHandler creates a new WSHandler.
func NewWSHandler(hub *ws.Hub, gm *ws.GameManager) *WSHandler {
	return &WSHandler{hub: hub, gm: gm}
}

// RegisterRoutes registers WebSocket routes.
func (h *WSHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/ws", h.ServeWS)
}

// ServeWS upgrades an HTTP connection to WebSocket.
func (h *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	if h.connCount.Load() >= maxWSConnections {
		http.Error(w, "too many connections", http.StatusServiceUnavailable)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}

	h.connCount.Add(1)

	clientID := uuid.New().String()
	client := &ws.Client{
		ID:   clientID,
		Conn: conn,
		Send: make(chan []byte, 256),
		Hub:  h.hub,
	}

	h.hub.Register(client)
	go h.writePump(client)
	go h.readPump(client, &h.connCount)
}

func (h *WSHandler) readPump(client *ws.Client, connCount *atomic.Int64) {
	defer func() {
		connCount.Add(-1)
		h.gm.HandleDisconnect(client.ID)
		h.hub.Unregister(client)
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(maxMsgSize)
	if err := client.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("ws: failed to set read deadline: %v", err)
		return
	}
	client.Conn.SetPongHandler(func(string) error {
		return client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("ws read error: %v", err)
			}
			break
		}
		h.gm.ProcessMessage(client.ID, message)
	}
}

func (h *WSHandler) writePump(client *ws.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if err := client.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("ws: failed to set write deadline: %v", err)
				return
			}
			if !ok {
				if err := client.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("ws: failed to send close message: %v", err)
				}
				return
			}
			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("ws write error: %v", err)
				return
			}

		case <-ticker.C:
			if err := client.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("ws: failed to set write deadline: %v", err)
				return
			}
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
