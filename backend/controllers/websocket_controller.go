package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adettinger/go-quizgame/livegame"
	"github.com/adettinger/go-quizgame/models"
	"github.com/adettinger/go-quizgame/socket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebSocketController handles WebSocket connections
type WebSocketController struct {
	manager       *socket.Manager
	upgrader      websocket.Upgrader
	liveGameStore *livegame.LiveGameStore
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
		liveGameStore: livegame.NewLiveGameStore(),
	}
}

// HandleConnection handles a new WebSocket connection
func (wsc *WebSocketController) HandleConnection(c *gin.Context) {
	log.Printf("=== NEW CONNECTION REQUEST ===")
	log.Printf("Request path: %s", c.Request.URL.Path)
	log.Printf("Request query: %s", c.Request.URL.RawQuery)
	log.Printf("Request headers: %v", c.Request.Header)
	log.Printf("Request method: %s", c.Request.Method)

	playerName := c.Param("playerName")

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
		Send:     make(chan models.Message, 256),
		UserData: make(map[string]interface{}),
	}
	client.UserData["name"] = playerName
	log.Print("New Client created")

	_, err = wsc.liveGameStore.AddPlayer(playerName)
	if err != nil {
		//TODO: handle duplicate player names
		errorToSend := models.Message{
			Type:       models.MessageTypeError,
			Timestamp:  time.Now(),
			PlayerName: "System",
			Content:    "Player name is already used",
		}
		jsonMessage, err := json.Marshal(errorToSend)
		if err != nil {
			log.Printf("[%s] Error marshaling message: %v", client.ID, err)
			close(client.Send)
			client.Conn.Close()
			return
		}
		if err := client.Conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
			log.Printf("[%s] Error writing message: %v", client.ID, err)
			close(client.Send)
			client.Conn.Close()
			return
		}
		close(client.Send)
		client.Conn.Close()
		return
		// TODO: Kill session?
	}

	wsc.manager.Register <- client
	log.Print("New Client registered")

	// Start goroutines for reading and writing
	go wsc.readPump(client)
	go wsc.writePump(client)
	log.Print("Read/Write pumps started")

	welcomeMsg := models.Message{
		Type:       models.MessageTypeChat,
		Timestamp:  time.Now(),
		PlayerName: "System",
		Content:    fmt.Sprintf("Welcome, %s!", playerName),
	}
	client.Send <- welcomeMsg
	log.Print("Welcome msg sent")
}

// readPump pumps messages from the WebSocket connection to the manager
func (wsc *WebSocketController) readPump(client *socket.Client) {
	log.Printf("[%s] Starting readPump", client.ID)
	defer func() {
		log.Printf("[%s] readPump exiting, unregistering client", client.ID)
		wsc.manager.Unregister <- client
		client.Conn.Close()
	}()

	// Set read deadline
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		log.Printf("[%s] Received pong", client.ID)
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		log.Printf("[%s] Waiting for next message...", client.ID)
		messageType, rawMessage, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[%s] Error reading message: %v", client.ID, err)
			} else {
				log.Printf("[%s] Connection closed: %v", client.ID, err)
			}
			break
		}
		log.Printf("[%s] Received message type %d: %s", client.ID, messageType, string(rawMessage))

		// Parse the incoming JSON message
		var message models.Message
		if err := json.Unmarshal(rawMessage, &message); err != nil {
			log.Printf("[%s] Error parsing message: %v", client.ID, err)

			// Send error message back to client
			errorMsg := models.Message{
				Type:      models.MessageTypeError,
				Timestamp: time.Now(),
				Content:   "Invalid message format",
			}
			select {
			case client.Send <- errorMsg:
				log.Printf("[%s] Sent error message", client.ID)
			default:
				log.Printf("[%s] Client send channel full, skipping error message", client.ID)
			}
			continue
		}

		// Add sender information to the message
		message.PlayerName = client.UserData["name"].(string)
		message.Timestamp = time.Now()

		log.Printf("[%s] Processing message of type %s", client.ID, message.Type)

		// Process the message based on its type
		switch message.Type {
		case models.MessageTypeChat:
			log.Printf("[%s] Broadcasting chat message", client.ID)
			wsc.manager.BroadcastMessage(message)
		case models.MessageTypeGameUpdate:
			log.Printf("[%s] Broadcasting game update", client.ID)
			wsc.manager.BroadcastMessage(message)
		default:
			log.Printf("[%s] Unknown message type: %s", client.ID, message.Type)
		}
	}
}

// writePump pumps messages from the manager to the WebSocket connection
func (wsc *WebSocketController) writePump(client *socket.Client) {
	log.Printf("[%s] Starting writePump", client.ID)
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		log.Printf("[%s] writePump exiting", client.ID)
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			log.Printf("[%s] Got message to send, channel ok: %v", client.ID, ok)

			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The manager closed the channel.
				log.Printf("[%s] Send channel closed, sending close message", client.ID)
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			log.Printf("[%s] Marshaling message of type %s", client.ID, message.Type)
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				log.Printf("[%s] Error marshaling message: %v", client.ID, err)
				continue
			}

			log.Printf("[%s] Writing message: %s", client.ID, string(jsonMessage))
			if err := client.Conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
				log.Printf("[%s] Error writing message: %v", client.ID, err)
				return
			}
			log.Printf("[%s] Message written successfully", client.ID)

		case <-ticker.C:
			log.Printf("[%s] Sending ping", client.ID)
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[%s] Error sending ping: %v", client.ID, err)
				return
			}
		}
	}
}

// GetManager returns the WebSocket manager
func (wsc *WebSocketController) GetManager() *socket.Manager {
	return wsc.manager
}
