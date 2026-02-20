package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adettinger/go-quizgame/controllers"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWebSocketController(t *testing.T) {
	controller := controllers.NewWebSocketController()

	assert.NotNil(t, controller)
	assert.NotNil(t, *controller.GetManager())
	// This is all that is public so it is all that can be tested on creation
}

// TestHandleConnection tests the HandleConnection method
func TestHandleConnection(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	const playerName = "Alex"

	// Create a test server with a websocket handler
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		// Read the first message sent by the client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Echo the message back
		conn.WriteMessage(websocket.TextMessage, msg)
	}))
	defer server.Close()

	// Replace the server URL scheme from http to ws
	// wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Create a test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "playerName", Value: playerName}}

	// Mock the request
	req, err := http.NewRequest("GET", "/ws/"+playerName, nil)
	require.NoError(t, err)
	c.Request = req

	// Create a controller with mocked dependencies
	controller := controllers.NewWebSocketController()

	// This is a partial test since we can't fully mock the websocket connection
	// In a real test, you would need to set up a test server and client
	// Here we're just ensuring the function doesn't panic
	// assert.NotPanics(t, func() {
	// 	// This will fail because we can't upgrade the connection in a test context
	// 	// but it shouldn't panic
	// 	controller.HandleConnection(c)
	// })
	manager := controller.GetManager()
	// go manager.Start()

	controller.HandleConnection(c)

	time.Sleep(1000 * time.Millisecond)

	// TODO: Assert connection is upgraded

	// Assert player added to game store
	player, err := manager.LiveGameStore.GetPlayerByName(playerName)
	assert.NoError(t, err)

	// Assert correct client was added to manager
	select {
	case client := <-manager.Register:
		assert.Equal(t, playerName, client.UserData.Name)
		assert.Equal(t, player.Id, client.UserData.PlayerId)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for registration message")
	}
	// assert.Equal(t, 1, manager.ClientCount())

	// var client *socket.Client
	// for _, c := range manager.Clients {
	// 	client = c
	// 	break
	// }
	// assert.Equal(t, playerName, client.UserData.Name)
	// assert.Equal(t, player.Id, client.UserData.PlayerId)

}

// // TestReadPump tests the readPump method
// func TestReadPump(t *testing.T) {
// 	// Create controller
// 	controller := controllers.NewWebSocketController()

// 	// Create mock connection
// 	mockConn := new(mocks.MockWebSocketConnection)

// 	// Setup mock expectations
// 	mockConn.On("SetReadDeadline", mock.Anything).Return(nil)
// 	mockConn.On("SetPongHandler", mock.Anything).Return()

// 	// First message is a valid chat message
// 	chatMsg := models.Message{
// 		Type:    models.MessageTypeChat,
// 		Content: models.MessageTextContent{Text: "Hello"},
// 	}
// 	chatMsgBytes, _ := json.Marshal(chatMsg)

// 	// Second message is a valid game update
// 	gameMsg := models.Message{
// 		Type:    models.MessageTypeGameUpdate,
// 		Content: models.MessageTextContent{Text: "Update"},
// 	}
// 	gameMsgBytes, _ := json.Marshal(gameMsg)

// 	// Third message is invalid JSON
// 	invalidMsg := []byte("invalid json")

// 	// Setup read sequence
// 	mockConn.On("ReadMessage").Return(websocket.TextMessage, chatMsgBytes, nil).Once()
// 	mockConn.On("ReadMessage").Return(websocket.TextMessage, gameMsgBytes, nil).Once()
// 	mockConn.On("ReadMessage").Return(websocket.TextMessage, invalidMsg, nil).Once()
// 	mockConn.On("ReadMessage").Return(0, []byte{}, websocket.ErrCloseSent).Once()

// 	// Mock manager
// 	mockManager := NewMockManager()
// 	mockManager.On("CreateNewClientID").Return("client-123")
// 	mockManager.On("BroadcastMessage", mock.Anything).Return()

// 	clientId := uuid.New()
// 	// Create client
// 	client := &socket.Client{
// 		ID:       clientId,
// 		Conn:     mockConn,
// 		Manager:  mockManager,
// 		Send:     make(chan models.Message, 256),
// 		UserData: socket.UserData{Name: "TestPlayer", PlayerId: clientId},
// 	}

// 	// Run readPump in a goroutine
// 	go controller.readPump(client)

// 	// Wait for all messages to be processed
// 	time.Sleep(100 * time.Millisecond)

// 	// Verify expectations
// 	mockConn.AssertExpectations(t)
// 	mockManager.AssertExpectations(t)
// }

// // TestWritePump tests the writePump method
// func TestWritePump(t *testing.T) {
// 	// Create controller
// 	controller := NewWebSocketController()

// 	// Create mock connection
// 	mockConn := new(MockWebSocketConnection)

// 	// Setup mock expectations
// 	mockConn.On("SetWriteDeadline", mock.Anything).Return(nil)
// 	mockConn.On("WriteMessage", websocket.TextMessage, mock.Anything).Return(nil)
// 	mockConn.On("WriteMessage", websocket.PingMessage, mock.Anything).Return(nil)
// 	mockConn.On("Close").Return(nil)

// 	// Create client
// 	client := &socket.Client{
// 		ID:       "client-123",
// 		Conn:     mockConn,
// 		Manager:  socket.NewManager(),
// 		Send:     make(chan models.Message, 256),
// 		UserData: socket.UserData{Name: "TestPlayer", PlayerId: "player-123"},
// 	}

// 	// Create test message
// 	testMsg := models.CreateMessage(
// 		models.MessageTypeChat,
// 		"TestPlayer",
// 		models.MessageTextContent{Text: "Hello"},
// 	)

// 	// Start writePump in a goroutine
// 	go func() {
// 		// Send a message to the client
// 		client.Send <- testMsg

// 		// Wait for the message to be processed
// 		time.Sleep(50 * time.Millisecond)

// 		// Close the send channel to terminate the writePump
// 		close(client.Send)
// 	}()

// 	// Run writePump
// 	controller.writePump(client)

// 	// Verify expectations
// 	mockConn.AssertExpectations(t)
// }
