package controllers

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/adettinger/go-quizgame/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleConnectionFull(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	const playerName = "Alex"

	// Create a router with the websocket handler
	router := gin.New()
	controller := NewWebSocketController()

	// Setup the route
	router.GET("/ws/:playerName", controller.HandleConnection)

	// Create a test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Convert the HTTP URL to a WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/" + playerName

	// Connect to the WebSocket server
	dialer := websocket.Dialer{}
	clientConn, _, err := dialer.Dial(wsURL, nil)
	require.NoError(t, err, "Failed to connect to WebSocket server")
	defer clientConn.Close()

	time.Sleep(1000 * time.Millisecond)

	manager := controller.GetManager()

	player, err := manager.LiveGameStore.GetPlayerByName(playerName)
	assert.NoError(t, err)

	// Assert manager called to register client
	select {
	case client := <-manager.Register:
		assert.Equal(t, playerName, client.UserData.Name)
		assert.Equal(t, player.Id, client.UserData.PlayerId)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for manager to get registration request")
	}

	// Wait for the welcome message
	var welcomeMsg models.Message
	err = clientConn.ReadJSON(&welcomeMsg)
	require.NoError(t, err, "Failed to read welcome message")

	// Verify the welcome message
	assert.Equal(t, models.MessageTypeChat, welcomeMsg.Type)
	assert.Equal(t, "System", welcomeMsg.PlayerName)
	textContent, ok := welcomeMsg.Content.(map[string]interface{})
	require.True(t, ok, "Content should be a map")
	assert.Equal(t, fmt.Sprintf("Welcome, %s!", playerName), textContent["Text"])

	// Wait for the player list message
	var playerListMsg models.Message
	err = clientConn.ReadJSON(&playerListMsg)
	require.NoError(t, err, "Failed to read player list message")

	// Verify the player list message
	assert.Equal(t, models.MessageTypePlayerList, playerListMsg.Type)
	assert.Equal(t, "System", playerListMsg.PlayerName)
}

// Test invalid player name
// func TestHandleConnectionInvalidName(t *testing.T) {
// 	// Setup
// 	gin.SetMode(gin.TestMode)

// 	// Create a router with the websocket handler
// 	router := gin.New()
// 	controller := NewWebSocketController()

// 	// Setup the route
// 	router.GET("/ws/:playerName", controller.HandleConnection)

// 	// Create a test server
// 	server := httptest.NewServer(router)
// 	defer server.Close()

// 	// Test with an empty name (invalid)
// 	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/"

// 	// Connect to the WebSocket server
// 	dialer := websocket.Dialer{}
// 	_, _, err := dialer.Dial(wsURL, nil)

// 	// Should fail with invalid name
// 	assert.Error(t, err, "Connection should fail with empty player name")

// 	// Test with a name that's too long
// 	longName := strings.Repeat("X", PlayerNameMaxLength+1)
// 	wsURL = "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/" + longName

// 	// Connect to the WebSocket server
// 	clientConn, _, err = dialer.Dial(wsURL, nil)

// 	// Should fail with name too long
// 	assert.Error(t, err, "Connection should fail with name too long")
// }

// // Test concurrent connections
// func TestHandleConnectionConcurrent(t *testing.T) {
// 	// Setup
// 	gin.SetMode(gin.TestMode)

// 	// Create a router with the websocket handler
// 	router := gin.New()
// 	controller := NewWebSocketController()

// 	// Setup the route
// 	router.GET("/ws/:playerName", controller.HandleConnection)

// 	// Create a test server
// 	server := httptest.NewServer(router)
// 	defer server.Close()

// 	// Number of concurrent connections
// 	numConnections := 5
// 	connections := make([]*websocket.Conn, numConnections)

// 	// Connect multiple clients
// 	for i := 0; i < numConnections; i++ {
// 		playerName := fmt.Sprintf("Player%d", i)
// 		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/" + playerName

// 		dialer := websocket.Dialer{}
// 		conn, _, err := dialer.Dial(wsURL, nil)
// 		require.NoError(t, err, "Failed to connect client %d", i)

// 		connections[i] = conn

// 		// Read welcome message
// 		var msg models.Message
// 		err = conn.ReadJSON(&msg)
// 		require.NoError(t, err, "Failed to read welcome message for client %d", i)

// 		// Read player list message
// 		err = conn.ReadJSON(&msg)
// 		require.NoError(t, err, "Failed to read player list message for client %d", i)
// 	}

// 	// Wait for all clients to be registered
// 	time.Sleep(100 * time.Millisecond)

// 	// Verify all clients are registered
// 	assert.Equal(t, numConnections, len(controller.manager.Clients), "All clients should be registered")

// 	// Verify player list
// 	// playerNames := mockGameStore.GetPlayerNameList()
// 	// assert.Equal(t, numConnections, len(playerNames), "All players should be in the game store")

// 	// Clean up connections
// 	for i, conn := range connections {
// 		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
// 		require.NoError(t, err, "Failed to close connection %d", i)
// 		conn.Close()
// 	}

// 	// Wait for all clients to be unregistered
// 	time.Sleep(100 * time.Millisecond)
// 	assert.Equal(t, 0, len(controller.manager.Clients), "All clients should be unregistered")
}
