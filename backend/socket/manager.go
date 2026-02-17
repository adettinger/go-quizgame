package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	livegame "github.com/adettinger/go-quizgame/liveGame"
	"github.com/adettinger/go-quizgame/models"
	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connection
type Client struct {
	ID       string
	Conn     *websocket.Conn
	Manager  *Manager
	Send     chan models.Message
	UserData map[string]interface{} // For storing user-specific data
}

func (client Client) Logf(message string, errs ...error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %v", client.ID, message))
	for _, err := range errs {
		sb.WriteString(": " + err.Error())
	}
	log.Print(sb.String())
}

func (client *Client) ErrorAndKill(errorMsg string) {
	errorToSend := models.Message{
		Type:       models.MessageTypeError,
		Timestamp:  time.Now(),
		PlayerName: "System",
		Content:    models.MessageTextContent{Text: errorMsg},
	}
	jsonMessage, err := json.Marshal(errorToSend)
	if err != nil {
		client.Logf("Error marshaling message: ", err)
		close(client.Send)
		client.Conn.Close()
		return
	}
	if err := client.Conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
		client.Logf("Error writing message", err)
		close(client.Send)
		client.Conn.Close()
		return
	}
	close(client.Send)
	client.Conn.Close()
}

// Manager manages WebSocket connections
type Manager struct {
	Clients       map[string]*Client
	Register      chan *Client
	Unregister    chan *Client
	Broadcast     chan models.Message
	LiveGameStore *livegame.LiveGameStore
	mutex         sync.Mutex
}

// Creates a new WebSocket manager
func NewManager() *Manager {
	return &Manager{
		Clients:       make(map[string]*Client),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Broadcast:     make(chan models.Message),
		LiveGameStore: livegame.NewLiveGameStore(),
	}
}

// Begins the WebSocket manager's operations
func (m *Manager) Start() {
	for {
		select {
		case client := <-m.Register:
			func() {
				m.mutex.Lock()
				defer m.mutex.Unlock()
				m.Clients[client.ID] = client
				client.Logf("Client connected")

			}()
			go func() {
				log.Printf("Broadcasting join message for %s", client.UserData["name"].(string))
				m.BroadcastMessage(models.Message{
					Type:       models.MessageTypeJoin,
					Timestamp:  time.Now(),
					PlayerName: client.UserData["name"].(string),
					Content:    models.MessageTextContent{Text: "has joined the game"},
				})
			}()
		case client := <-m.Unregister:
			if _, ok := m.Clients[client.ID]; ok {
				leaveMsg := models.Message{
					Type:       models.MessageTypeLeave,
					Timestamp:  time.Now(),
					PlayerName: client.UserData["name"].(string),
					Content:    models.MessageTextContent{Text: "has left the game"},
				}
				func() {
					m.mutex.Lock()
					defer m.mutex.Unlock()
					delete(m.Clients, client.ID)
					close(client.Send)
					client.Logf("Client disconnected: %s")
				}()

				go func() {
					m.LiveGameStore.RemovePlayerByName(client.UserData["name"].(string))
					log.Printf("Broadcasting leave message for %s", client.UserData["name"].(string))
					m.BroadcastMessage(leaveMsg)
				}()
			}
		case message := <-m.Broadcast:
			log.Printf("Broadcasting message of type %s from %s", message.Type, message.PlayerName)

			// Efficiency concern: Copying clients to avoid holding lock during sends
			var clients map[string]*Client
			func() {
				m.mutex.Lock()
				defer m.mutex.Unlock()

				clients = make(map[string]*Client, len(m.Clients))
				for id, client := range m.Clients {
					clients[id] = client
				}
			}()

			log.Printf("Sending to %d clients", len(clients))

			for id, client := range clients {
				func(id string, client *Client) {
					select {
					case client.Send <- message:
						client.Logf("Message queued")
					default:
						client.Logf("Send buffer full, closing")

						// Use a goroutine to avoid deadlock when unregistering
						go func() {
							m.Unregister <- client
						}()
					}
				}(id, client)
			}
		}
	}
}

// SendToClient sends a message to a specific client
func (m *Manager) SendToClient(clientID string, message models.Message) bool {
	m.mutex.Lock()
	client, exists := m.Clients[clientID]
	m.mutex.Unlock()

	if !exists {
		return false
	}

	client.Send <- message
	return true
}

// BroadcastMessage sends a message to all connected clients
func (m *Manager) BroadcastMessage(message models.Message) {
	log.Printf("Queueing broadcast message of type %s", message.Type)
	m.Broadcast <- message
}

// ClientCount returns the number of connected clients
func (m *Manager) ClientCount() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return len(m.Clients)
}
