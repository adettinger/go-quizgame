package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/adettinger/go-quizgame/models"
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
	log.Printf("=== NEW CONNECTION REQUEST ===")
	log.Printf("Request path: %s", c.Request.URL.Path)
	log.Printf("Request query: %s", c.Request.URL.RawQuery)
	log.Printf("Request headers: %v", c.Request.Header)
	log.Printf("Request method: %s", c.Request.Method)
	log.Printf("=== ===")

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

	// Validate player name
	playerName = strings.TrimSpace(playerName)
	if len(playerName) == 0 {
		client.ErrorAndKill("Invalid name")
		return
	}

	_, err = wsc.manager.LiveGameStore.AddPlayer(playerName)
	if err != nil {
		client.ErrorAndKill("Player name is already used")
		return
	}

	wsc.manager.Register <- client
	log.Print("New Client registered")

	// Start goroutines for reading and writing
	go wsc.readPump(client)
	go wsc.writePump(client)
	client.Logf("Read/Write pumps started")

	welcomeMsg := models.Message{
		Type:       models.MessageTypeChat,
		Timestamp:  time.Now(),
		PlayerName: "System",
		Content:    models.MessageTextContent{Text: fmt.Sprintf("Welcome, %s!", playerName)},
	}
	client.Send <- welcomeMsg

	existingPlayers := models.Message{
		Type:       models.MessageTypePlayerList,
		Timestamp:  time.Now(),
		PlayerName: "System",
		Content:    models.PlayerListMessageContent{Names: wsc.manager.LiveGameStore.GetPlayerNameList()},
	}
	client.Send <- existingPlayers
}

// readPump pumps messages from the WebSocket connection to the manager
func (wsc *WebSocketController) readPump(client *socket.Client) {
	client.Logf("Starting readPump")
	defer func() {
		client.Logf("readPump exiting, unregistering client")
		wsc.manager.Unregister <- client
		client.Conn.Close()
	}()

	// Set read deadline
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Logf("Received pong")
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		client.Logf("Waiting for next message...")
		messageType, rawMessage, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				client.Logf("Error reading message", err)
			} else {
				client.Logf("Connection closed", err)
			}
			break
		}
		client.Logf(fmt.Sprintf("Received message type %d: %s", messageType, string(rawMessage)))

		// Parse the incoming JSON message
		var message models.Message
		if err := json.Unmarshal(rawMessage, &message); err != nil {
			client.Logf("Error parsing message", err)

			// Send error message back to client
			errorMsg := models.Message{
				Type:      models.MessageTypeError,
				Timestamp: time.Now(),
				Content:   models.MessageTextContent{Text: "Invalid message format"},
			}
			select {
			case client.Send <- errorMsg:
				client.Logf("Sent error message")
			default:
				client.Logf("Client send channel full, skipping error message")
			}
			continue
		}

		// Add sender information to the message
		message.PlayerName = client.UserData["name"].(string)
		message.Timestamp = time.Now()

		client.Logf(fmt.Sprintf("Processing message of type %s", message.Type))

		// Process the message based on its type
		switch message.Type {
		case models.MessageTypeChat:
			client.Logf("Broadcasting chat message")
			wsc.manager.BroadcastMessage(message)
		case models.MessageTypeGameUpdate:
			client.Logf("Broadcasting game update")
			wsc.manager.BroadcastMessage(message)
		default:
			client.Logf(fmt.Sprintf("Unknown message type: %s", message.Type))
		}
	}
}

// writePump pumps messages from the manager to the WebSocket connection
func (wsc *WebSocketController) writePump(client *socket.Client) {
	client.Logf("Starting writePump")
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		client.Logf("writePump exiting")
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Logf(fmt.Sprintf("Got message to send, channel ok: %v", ok))

			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The manager closed the channel.
				client.Logf("Send channel closed, sending close message")
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			client.Logf(fmt.Sprintf("Marshaling message of type %s", message.Type))
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				client.Logf("Error marshaling message", err)
				continue
			}

			client.Logf(fmt.Sprintf("Writing message: %s", string(jsonMessage)))
			if err := client.Conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
				client.Logf("Error writing message", err)
				return
			}
			client.Logf("Message written successfully")

		case <-ticker.C:
			client.Logf("Sending ping")
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				client.Logf("Error sending ping: %v", err)
				return
			}
		}
	}
}

// GetManager returns the WebSocket manager
func (wsc *WebSocketController) GetManager() *socket.Manager {
	return wsc.manager
}
