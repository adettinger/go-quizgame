package controllers

import (
	"log"
	"net/http"

	"github.com/adettinger/go-quizgame/socket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebSocketController handles WebSocket connections
type WebSocketController struct {
	manager  *socket.Manager
	upgrader websocket.Upgrader
}

// NewWebSocketController creates a new WebSocket controller
func NewWebSocketController() *WebSocketController {
	return &WebSocketController{
		manager: socket.NewManager(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for now - you might want to restrict this in production
				return true
			},
		},
	}
}

// HandleConnection handles a new WebSocket connection
func (wsc *WebSocketController) HandleConnection(c *gin.Context) {
	conn, err := wsc.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	clientID := uuid.New().String()
	client := &socket.Client{
		ID:       clientID,
		Conn:     conn,
		Manager:  wsc.manager,
		Send:     make(chan []byte, 256),
		UserData: make(map[string]interface{}),
	}

	wsc.manager.Register <- client

	// Start goroutines for reading and writing
	go wsc.readPump(client)
	go wsc.writePump(client)
}

// readPump pumps messages from the WebSocket connection to the manager
func (wsc *WebSocketController) readPump(client *socket.Client) {
	defer func() {
		wsc.manager.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Process the message here
		// For now, we'll just broadcast it back to all clients
		wsc.manager.BroadcastMessage(message)
	}
}

// writePump pumps messages from the manager to the WebSocket connection
func (wsc *WebSocketController) writePump(client *socket.Client) {
	defer client.Conn.Close()

	for message := range client.Send {
		w, err := client.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		// Add queued messages to the current WebSocket message
		n := len(client.Send)
		for range n {
			w.Write(<-client.Send)
		}

		if err := w.Close(); err != nil {
			return
		}
	}
}

// GetManager returns the WebSocket manager
func (wsc *WebSocketController) GetManager() *socket.Manager {
	return wsc.manager
}
