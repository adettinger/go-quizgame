package socket

import (
	"errors"
	"fmt"
	"log"
	"sync"

	livegame "github.com/adettinger/go-quizgame/liveGame"
	"github.com/adettinger/go-quizgame/models"
	"github.com/adettinger/go-quizgame/webserver"
	"github.com/google/uuid"
)

// Manager manages WebSocket connections
type Manager struct {
	PlayerClients map[uuid.UUID]*Client
	HostClient    *Client
	Register      chan *Client
	Unregister    chan *Client
	Broadcast     chan models.Message
	LiveGameStore *livegame.LiveGameStore
	QuestionStore *webserver.QuestionStore
	mutex         sync.RWMutex
}

// Creates a new WebSocket manager
func NewManager(qs *webserver.QuestionStore) *Manager {
	return &Manager{
		PlayerClients: make(map[uuid.UUID]*Client),
		HostClient:    nil,
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Broadcast:     make(chan models.Message),
		LiveGameStore: livegame.NewLiveGameStore(qs),
		QuestionStore: qs,
	}
}

// Begins the WebSocket manager's operations
func (m *Manager) Start() {
	for {
		select {
		case client := <-m.Register:
			func() {
				// TODO: Check for error response when adding clients
				m.AddClient(client)
				client.Logf("Client added to manager")
			}()
			if !client.UserData.IsHost {
				go func() {
					log.Printf("Broadcasting join message for %s", client.UserData.Name)
					m.BroadcastMessage(models.CreateMessage(
						models.MessageTypeJoin,
						client.UserData.Name,
						models.MessageTextContent{Text: "has joined the game"},
					))
				}()
			}
		case client := <-m.Unregister:
			if _, ok := m.PlayerClients[client.ID]; ok {
				leaveMsg := models.CreateMessage(
					models.MessageTypeLeave,
					client.UserData.Name,
					models.MessageTextContent{Text: "has left the game"})
				func() {
					m.mutex.Lock()
					defer m.mutex.Unlock()
					delete(m.PlayerClients, client.ID)
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
			if m.HostClient != nil && m.HostClient.ID == client.ID {
				log.Print("Host is leaving game")
				leaveMsg := models.CreateMessage(
					models.MessageTypeLeave,
					"Host",
					models.MessageTextContent{Text: "has left the game"})
				func() {
					m.mutex.Lock()
					defer m.mutex.Unlock()
					m.HostClient = nil
					close(client.Send)
					client.Logf("Client disconnected: %s")
				}()

				go func() {
					// TODO: Removing player from game store can be done in controller that writes to unregister
					log.Printf("Broadcasting leave message for %s", client.UserData.Name)
					m.BroadcastMessage(leaveMsg)
					for _, client := range m.PlayerClients {
						m.Unregister <- client
					}
					m.LiveGameStore.KillGame()
				}()
				log.Print("Killed game")
			}
		case message := <-m.Broadcast:
			log.Printf("Broadcasting message of type %s from %s", message.Type, message.PlayerName)

			// Efficiency concern: Copying clients to avoid holding lock during sends
			var clients map[uuid.UUID]*Client
			func() {
				m.mutex.Lock()
				defer m.mutex.Unlock()

				clients = make(map[uuid.UUID]*Client, len(m.PlayerClients)+1)
				if m.HostClient != nil {
					clients[m.HostClient.ID] = m.HostClient
				}
				for id, client := range m.PlayerClients {
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
	client, exists := m.PlayerClients[clientID]
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

// PlayerClientCount returns the number of connected clients
func (m *Manager) PlayerClientCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.PlayerClients)
}

func (m *Manager) CreateNewClientID() uuid.UUID {
	for {
		id := uuid.New()
		if !m.PlayerClientIDExists(id) {
			return id
		}
	}
}

func (m *Manager) PlayerClientIDExists(id uuid.UUID) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, ok := m.PlayerClients[id]
	return ok
}

func (m *Manager) AddClient(client *Client) error {
	m.mutex.RLock()
	hasHost := m.HostClient != nil
	m.mutex.RUnlock()
	if client.UserData.IsHost {
		if hasHost {
			return errors.New("Game already has a host")
		}
		m.mutex.Lock()
		defer m.mutex.Unlock()
		m.HostClient = client
	} else {
		// TODO: What if no host
		// if !hasHost {
		// 	return errors.New("Need host before player client can be added")
		// }
		if m.PlayerClientIDExists(client.ID) {
			return fmt.Errorf("Client already exists with ID: %v", client.ID.String())
		}
		m.mutex.Lock()
		defer m.mutex.Unlock()
		m.PlayerClients[client.ID] = client
	}
	return nil
}
