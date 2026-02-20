package socket

import (
	"fmt"
	"log"
	"sync"

	livegame "github.com/adettinger/go-quizgame/liveGame"
	"github.com/adettinger/go-quizgame/models"
	"github.com/google/uuid"
)

// Manager manages WebSocket connections
type Manager struct {
	Clients       map[uuid.UUID]*Client
	Register      chan *Client
	Unregister    chan *Client
	Broadcast     chan models.Message
	LiveGameStore *livegame.LiveGameStore
	mutex         sync.RWMutex
}

// Creates a new WebSocket manager
func NewManager() *Manager {
	return &Manager{
		Clients:       make(map[uuid.UUID]*Client),
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
				m.AddClient(client)
				client.Logf("Client added to manager")

			}()
			go func() {
				log.Printf("Broadcasting join message for %s", client.UserData.Name)
				m.BroadcastMessage(models.CreateMessage(
					models.MessageTypeJoin,
					client.UserData.Name,
					models.MessageTextContent{Text: "has joined the game"},
				))
			}()
		case client := <-m.Unregister:
			if _, ok := m.Clients[client.ID]; ok {
				leaveMsg := models.CreateMessage(
					models.MessageTypeLeave,
					client.UserData.Name,
					models.MessageTextContent{Text: "has left the game"})
				func() {
					m.mutex.Lock()
					defer m.mutex.Unlock()
					delete(m.Clients, client.ID)
					close(client.Send)
					client.Logf("Client disconnected: %s")
				}()

				go func() {
					// TODO: Removing player from game store can be done in controller that writes to unregister
					m.LiveGameStore.RemovePlayerByName(client.UserData.Name)
					log.Printf("Broadcasting leave message for %s", client.UserData.Name)
					m.BroadcastMessage(leaveMsg)
				}()
			}
		case message := <-m.Broadcast:
			log.Printf("Broadcasting message of type %s from %s", message.Type, message.PlayerName)

			// Efficiency concern: Copying clients to avoid holding lock during sends
			var clients map[uuid.UUID]*Client
			func() {
				m.mutex.Lock()
				defer m.mutex.Unlock()

				clients = make(map[uuid.UUID]*Client, len(m.Clients))
				for id, client := range m.Clients {
					clients[id] = client
				}
			}()

			log.Printf("Sending to %d clients", len(clients))

			for id, client := range clients {
				func(id uuid.UUID, client *Client) {
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
func (m *Manager) SendToClient(clientID uuid.UUID, message models.Message) bool {
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
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.Clients)
}

func (m *Manager) CreateNewClientID() uuid.UUID {
	for {
		id := uuid.New()
		if !m.ClientIDExists(id) {
			return id
		}
	}
}

func (m *Manager) ClientIDExists(id uuid.UUID) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, ok := m.Clients[id]
	return ok
}

func (m *Manager) AddClient(client *Client) error {
	if m.ClientIDExists(client.ID) {
		return fmt.Errorf("Client already exists with ID: %v", client.ID.String())
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Clients[client.ID] = client
	return nil
}
