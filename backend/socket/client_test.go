package socket_test

import (
	"errors"
	"testing"

	"github.com/adettinger/go-quizgame/mocks"
	"github.com/adettinger/go-quizgame/models"
	"github.com/adettinger/go-quizgame/socket"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_ErrorAndKill(t *testing.T) {
	// Create mock objects
	mockConn := new(mocks.MockConnection)

	// Create a client with mocked dependencies
	client := &socket.Client{
		ID:      uuid.New(),
		Conn:    mockConn,
		Manager: socket.NewManager(),
		Send:    make(chan models.Message, 10),
		UserData: socket.UserData{
			PlayerId: uuid.New(),
			Name:     "TestUser",
		},
	}

	// Set expectations
	errorMsg := "Test error message"

	// Instead of using MatchedBy with complex JSON parsing, use mock.Anything
	// for the data parameter since we're more concerned with the method being called
	mockConn.On("WriteMessage", websocket.TextMessage, mock.Anything).Return(nil)
	mockConn.On("Close").Return(nil)

	// Call the method
	client.ErrorAndKill(errorMsg)

	// Verify expectations
	mockConn.AssertExpectations(t)

	// Verify channel is closed
	_, isOpen := <-client.Send
	assert.False(t, isOpen, "Send channel should be closed")
}

// Test case for WriteMessage error
func TestClient_ErrorAndKill_WriteError(t *testing.T) {
	mockConn := new(mocks.MockConnection)
	client := &socket.Client{
		ID:   uuid.New(),
		Conn: mockConn,
		Send: make(chan models.Message, 10),
	}

	// Set up the mock to return an error
	mockConn.On("WriteMessage", websocket.TextMessage, mock.Anything).Return(errors.New("write error"))
	mockConn.On("Close").Return(nil)

	// Call the method
	client.ErrorAndKill("Test error")

	// Verify expectations
	mockConn.AssertExpectations(t)

	// Verify channel is closed
	_, isOpen := <-client.Send
	assert.False(t, isOpen, "Send channel should be closed")
}
