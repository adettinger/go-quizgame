package socket

import (
	"log"
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
				log.Printf("Client connected: %s", client.ID)

			}()
			go func() {
				log.Printf("Broadcasting join message for %s", client.UserData["name"].(string))
				m.BroadcastMessage(models.Message{
					Type:       models.MessageTypeJoin,
					Timestamp:  time.Now(),
					PlayerName: client.UserData["name"].(string),
					Content:    "has joined the game",
				})
			}()
		case client := <-m.Unregister:
			if _, ok := m.Clients[client.ID]; ok {
				leaveMsg := models.Message{
					Type:       models.MessageTypeLeave,
					Timestamp:  time.Now(),
					PlayerName: client.UserData["name"].(string),
					Content:    "has left the game",
				}
				func() {
					m.mutex.Lock()
					defer m.mutex.Unlock()
					delete(m.Clients, client.ID)
					close(client.Send)
					log.Printf("Client disconnected: %s", client.ID)
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
						log.Printf("Message queued for client %s", id)
					default:
						log.Printf("Client %s send buffer full, closing", id)

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
