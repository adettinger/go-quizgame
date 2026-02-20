package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adettinger/go-quizgame/models"
	"github.com/adettinger/go-quizgame/socket"
	"github.com/adettinger/go-quizgame/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const PlayerNameMaxLength = 20

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

	// TODO: reorder validation?
	conn, err := wsc.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		// TODO: better response
		return
	}

	client := &socket.Client{
		ID:       wsc.manager.CreateNewClientID(),
		Conn:     conn,
		Manager:  wsc.manager,
		Send:     make(chan models.Message, 256),
		UserData: socket.UserData{},
	}
	log.Print("New Client created")

	if !utils.IsPlayerNameValid(playerName, PlayerNameMaxLength) {
		client.ErrorAndKill("Invalid name")
		return
	}

	playerId, err := wsc.manager.LiveGameStore.AddPlayer(playerName)
	if err != nil {
		client.ErrorAndKill("Failed to add player to game store")
		return
	}
	client.UserData.Name = playerName
	client.UserData.PlayerId = playerId

	wsc.manager.Register <- client
	client.Logf("New Client registered")

	// Start goroutines for reading and writing
	go wsc.readPump(client)
	go wsc.writePump(client)
	client.Logf("Read/Write pumps started")

	// Send welcome message to new client
	client.Send <- models.CreateMessage(
		models.MessageTypeChat,
		"System",
		models.MessageTextContent{Text: fmt.Sprintf("Welcome, %s!", playerName)},
	)

	// Send list of players to new client
	client.Send <- models.CreateMessage(
		models.MessageTypePlayerList,
		"System",
		models.PlayerListMessageContent{Names: wsc.manager.LiveGameStore.GetPlayerNameList()},
	)
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
		client.Logf("Received message type: ", messageType, ": ", string(rawMessage))

		var message models.Message
		if err := json.Unmarshal(rawMessage, &message); err != nil {
			client.Logf("Error parsing message", err)

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

		message.PlayerName = client.UserData.Name
		message.Timestamp = time.Now()

		client.Logf("Processing message of type: ", message.Type)

		switch message.Type {
		case models.MessageTypeChat:
			client.Logf("Broadcasting chat message")
			wsc.manager.BroadcastMessage(message)
		case models.MessageTypeGameUpdate:
			client.Logf("Broadcasting game update")
			wsc.manager.BroadcastMessage(message)
		default:
			client.Logf("Unknown message type: ", message.Type)
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
			client.Logf("Got message to send, channel ok: ", ok)

			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The manager closed the channel.
				client.Logf("Send channel closed, sending close message")
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			client.Logf("Marshaling message of type: ", message.Type)
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				client.Logf("Error marshaling message", err)
				continue
			}

			client.Logf("Writing message: ", string(jsonMessage))
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

func (wsc *WebSocketController) GetManager() *socket.Manager {
	return wsc.manager
}
