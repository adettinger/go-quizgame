package socket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connection
type Client struct {
	ID       string
	Conn     *websocket.Conn
	Manager  *Manager
	Send     chan []byte
	UserData map[string]interface{} // For storing user-specific data
}

// Manager manages WebSocket connections
type Manager struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mutex      sync.Mutex
}

// Creates a new WebSocket manager
func NewManager() *Manager {
	return &Manager{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

// Begins the WebSocket manager's operations
func (m *Manager) Start() {
	for {
		select {
		case client := <-m.Register:
			m.mutex.Lock()
			m.Clients[client.ID] = client
			m.mutex.Unlock()
			log.Printf("Client connected: %s", client.ID)
		case client := <-m.Unregister:
			if _, ok := m.Clients[client.ID]; ok {
				m.mutex.Lock()
				delete(m.Clients, client.ID)
				close(client.Send)
				m.mutex.Unlock()
				log.Printf("Client disconnected: %s", client.ID)
			}
		case message := <-m.Broadcast:
			for _, client := range m.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					m.mutex.Lock()
					delete(m.Clients, client.ID)
					m.mutex.Unlock()
				}
			}
		}
	}
}

// SendToClient sends a message to a specific client
func (m *Manager) SendToClient(clientID string, message []byte) bool {
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
func (m *Manager) BroadcastMessage(message []byte) {
	m.Broadcast <- message
}

// ClientCount returns the number of connected clients
func (m *Manager) ClientCount() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return len(m.Clients)
}
