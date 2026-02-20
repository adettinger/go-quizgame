package socket_test

import (
	"sync"
	"testing"
	"time"

	"github.com/adettinger/go-quizgame/mocks"
	"github.com/adettinger/go-quizgame/models"
	"github.com/adettinger/go-quizgame/socket"
	"github.com/adettinger/go-quizgame/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const testClientBufferSize = 10

func createTestClient(t testing.TB, name string) *socket.Client {
	t.Helper()

	mockConn := new(mocks.MockConnection)

	return &socket.Client{
		ID:      uuid.New(),
		Conn:    mockConn,
		Manager: socket.NewManager(),
		Send:    make(chan models.Message, testClientBufferSize),
		UserData: socket.UserData{
			PlayerId: uuid.New(),
			Name:     name,
		},
	}
}

func TestNewManager(t *testing.T) {
	manager := socket.NewManager()

	assert.NotNil(t, manager)
	assert.NotNil(t, manager.Clients)
	assert.NotNil(t, manager.Register)
	assert.NotNil(t, manager.Unregister)
	assert.NotNil(t, manager.Broadcast)
	assert.NotNil(t, manager.LiveGameStore)
	assert.Equal(t, 0, manager.ClientCount())
}

func TestManager_ClientCount(t *testing.T) {
	manager := socket.NewManager()

	// Initially zero clients
	assert.Equal(t, 0, manager.ClientCount())

	// Add a client
	client := createTestClient(t, "testUser")
	manager.Clients[client.ID] = client

	assert.Equal(t, 1, manager.ClientCount())
}

func TestManager_CreateNewClientID(t *testing.T) {
	manager := socket.NewManager()

	// Create a new ID
	id1 := manager.CreateNewClientID()
	assert.NotEqual(t, uuid.Nil, id1)

	// Create another ID - should be different
	id2 := manager.CreateNewClientID()
	assert.NotEqual(t, uuid.Nil, id2)
	assert.NotEqual(t, id1, id2)

	// Create another ID - should not be id1
	id3 := manager.CreateNewClientID()
	assert.NotEqual(t, uuid.Nil, id3)
	assert.NotEqual(t, id1, id3)
	assert.NotEqual(t, id2, id3)
}

func TestManager_AddClient(t *testing.T) {
	manager := socket.NewManager()

	// Can add client
	client1 := createTestClient(t, "client1")
	err := manager.AddClient(client1)
	testutils.AssertNoError(t, err)
	assert.Equal(t, 1, manager.ClientCount())

	// Cannot add duplicate id client
	client2 := createTestClient(t, "client2")
	client2.ID = client1.ID
	err = manager.AddClient(client2)
	testutils.AssertHasError(t, err)

}

func TestManager_ClientIDExists(t *testing.T) {
	manager := socket.NewManager()

	// Create a client
	client := createTestClient(t, "testUser")
	manager.Clients[client.ID] = client

	// Check if the client ID exists
	assert.True(t, manager.ClientIDExists(client.ID))

	// Check if a non-existent ID exists
	assert.False(t, manager.ClientIDExists(uuid.New()))
}

func TestManager_SendToClient(t *testing.T) {
	manager := socket.NewManager()

	// Create a client
	client := createTestClient(t, "testUser")
	manager.Clients[client.ID] = client

	// Create a message
	message := models.CreateMessage(
		models.MessageTypeChat,
		"System",
		models.MessageTextContent{Text: "Test message"},
	)

	// Send the message to the client
	success := manager.SendToClient(client.ID, message)
	assert.True(t, success)

	// Verify the message was sent to the client
	select {
	case receivedMsg := <-client.Send:
		assert.Equal(t, message, receivedMsg)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for message")
	}

	// Try to send to a non-existent client
	success = manager.SendToClient(uuid.New(), message)
	assert.False(t, success)
}

func TestManager_BroadcastMessage(t *testing.T) {
	manager := socket.NewManager()

	// Create a message
	message := models.CreateMessage(
		models.MessageTypeChat,
		"System",
		models.MessageTextContent{Text: "Test broadcast"},
	)

	// Start a goroutine to receive from the broadcast channel
	var receivedMsg models.Message
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		receivedMsg = <-manager.Broadcast
	}()

	// Broadcast the message
	manager.BroadcastMessage(message)

	// Wait for the goroutine to complete
	wg.Wait()

	// Verify the message was broadcast
	assert.Equal(t, message, receivedMsg)
}

func TestManager_Start_Register(t *testing.T) {
	manager := socket.NewManager()

	// Create a clients
	client1 := createTestClient(t, "testUser")
	client2 := createTestClient(t, "testUser2")

	// Start the manager in a goroutine
	go manager.Start()

	// Register the client
	manager.Register <- client1

	// Give some time for processing
	time.Sleep(1000 * time.Millisecond)

	// Verify the client was registered
	assert.Equal(t, 1, manager.ClientCount())
	assert.True(t, manager.ClientIDExists(client1.ID))

	// Verify a join message was broadcast to client
	select {
	case msg := <-client1.Send:
		assert.Equal(t, models.MessageTypeJoin, msg.Type)
		assert.Equal(t, client1.UserData.Name, msg.PlayerName)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for broadcast message to client1 for client1 joining")
	}

	// Register the 2nd client
	manager.Register <- client2

	// Give some time for processing
	time.Sleep(100 * time.Millisecond)

	// Verify the client was registered
	assert.Equal(t, 2, manager.ClientCount())
	assert.True(t, manager.ClientIDExists(client2.ID))

	// Verify a join message was broadcast to client
	select {
	case msg := <-client1.Send:
		assert.Equal(t, models.MessageTypeJoin, msg.Type)
		assert.Equal(t, client2.UserData.Name, msg.PlayerName)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for broadcast message to client1 for client2 joining")
	}

	select {
	case msg := <-client2.Send:
		assert.Equal(t, models.MessageTypeJoin, msg.Type)
		assert.Equal(t, client2.UserData.Name, msg.PlayerName)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for broadcast message to client2 for client2 joining")
	}
}

func TestManager_Start_Unregister(t *testing.T) {
	const playerName = "testPlayer"
	manager := socket.NewManager()

	// Create a client1
	client1 := createTestClient(t, playerName)
	client2 := createTestClient(t, "player2")

	// Add the client directly to the manager
	manager.Clients[client1.ID] = client1
	manager.Clients[client2.ID] = client2
	// Add player to the liveGameStore
	manager.LiveGameStore.AddPlayer(playerName)
	manager.LiveGameStore.AddPlayer("player2")

	// Start the manager in a goroutine
	go func() {
		manager.Start()
	}()

	// Unregister the client
	manager.Unregister <- client1

	// Give some time for processing
	time.Sleep(100 * time.Millisecond)

	// Verify the client was unregistered
	assert.Equal(t, 1, manager.ClientCount())
	assert.False(t, manager.ClientIDExists(client1.ID))

	// Verify player was removed from gamestore
	assert.False(t, manager.LiveGameStore.PlayerExistsByName(playerName))

	// Verify a leave message was broadcast
	select {
	case msg := <-client2.Send:
		assert.Equal(t, models.MessageTypeLeave, msg.Type)
		assert.Equal(t, client1.UserData.Name, msg.PlayerName)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for broadcast message")
	}

	// Verify the client's send channel was closed
	_, isOpen := <-client1.Send
	assert.False(t, isOpen, "Client send channel should be closed")
}

func TestManager_Start_Broadcast(t *testing.T) {
	manager := socket.NewManager()

	// Create clients
	client1 := createTestClient(t, "testUser")
	client2 := createTestClient(t, "testUser")

	// Add clients to the manager
	manager.Clients[client1.ID] = client1
	manager.Clients[client2.ID] = client2

	// Start the manager in a goroutine
	go func() {
		manager.Start()
	}()

	// Create a message to broadcast
	message := models.CreateMessage(
		models.MessageTypeChat,
		"System",
		models.MessageTextContent{Text: "Test broadcast"},
	)

	// Broadcast the message
	manager.Broadcast <- message

	// Give some time for processing
	time.Sleep(100 * time.Millisecond)

	// Verify both clients received the message
	select {
	case msg1 := <-client1.Send:
		assert.Equal(t, message, msg1)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Client 1 timed out waiting for message")
	}

	select {
	case msg2 := <-client2.Send:
		assert.Equal(t, message, msg2)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Client 2 timed out waiting for message")
	}
}

func TestManager_Start_BroadcastFullBuffer(t *testing.T) {
	manager := socket.NewManager()

	// Create a client with a full send buffer
	client := createTestClient(t, "testUser")
	// Fill the buffer
	for i := 0; i < testClientBufferSize; i++ {
		client.Send <- models.CreateMessage(
			models.MessageTypeChat,
			"System",
			models.MessageTextContent{Text: "Buffer message"},
		)
	}

	// Add client to the manager
	manager.Clients[client.ID] = client

	// Start the manager in a goroutine
	go func() {
		manager.Start()
	}()

	// Create a message to broadcast
	message := models.CreateMessage(
		models.MessageTypeChat,
		"System",
		models.MessageTextContent{Text: "Test broadcast"},
	)

	// Broadcast the message
	manager.Broadcast <- message

	// Give some time for processing
	time.Sleep(100 * time.Millisecond)

	// Verify the client was unregistered due to full buffer
	assert.Equal(t, 0, manager.ClientCount())
}
