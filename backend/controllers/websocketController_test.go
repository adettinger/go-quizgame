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

func TestHandleConnectionInvalidName(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create a router with the websocket handler
	router := gin.New()
	controller := NewWebSocketController()

	// Setup the route
	router.GET("/ws/:playerName", controller.HandleConnection)

	// Create a test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Test with an empty name (invalid)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/"

	// Connect to the WebSocket server
	dialer := websocket.Dialer{}
	_, _, err := dialer.Dial(wsURL, nil)

	// Should fail with invalid name
	assert.Error(t, err, "Connection should fail with empty player name")

	// Test with a name that's too long
	longName := strings.Repeat("X", PlayerNameMaxLength+1)
	wsURL = "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/" + longName

	// Connect to the WebSocket server
	_, _, err = dialer.Dial(wsURL, nil)

	// Should fail with name too long
	assert.Error(t, err, "Connection should fail with name too long")
}
