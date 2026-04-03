package handler

import (
	"log"
	"net/http"
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
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: restrict to CORS_ORIGIN in production
		return true
	},
}

// WSHandler handles WebSocket connections.
type WSHandler struct {
	hub *ws.Hub
	gm  *ws.GameManager
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}

	clientID := uuid.New().String()
	client := &ws.Client{
		ID:   clientID,
		Conn: conn,
		Send: make(chan []byte, 256),
		Hub:  h.hub,
	}

	h.hub.Register(client)
	go h.writePump(client)
	go h.readPump(client)
}

func (h *WSHandler) readPump(client *ws.Client) {
	defer func() {
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
