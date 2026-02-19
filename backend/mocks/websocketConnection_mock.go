package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockConnection struct {
	mock.Mock
}

// Close mocks the Close method
func (m *MockConnection) Close() error {
	args := m.Called()
	return args.Error(0)
}

// SetPongHandler mocks the SetPongHandler method
func (m *MockConnection) SetPongHandler(h func(appData string) error) {
	m.Called(h)
	// No return value needed as the original function doesn't return anything
}

// ReadMessage mocks the ReadMessage method
func (m *MockConnection) ReadMessage() (messageType int, p []byte, err error) {
	args := m.Called()
	return args.Int(0), args.Get(1).([]byte), args.Error(2)
}

// WriteMessage mocks the WriteMessage method
func (m *MockConnection) WriteMessage(messageType int, data []byte) error {
	args := m.Called(messageType, data)
	return args.Error(0)
}

// SetReadDeadline mocks the SetReadDeadline method
func (m *MockConnection) SetReadDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

// SetWriteDeadline mocks the SetWriteDeadline method
func (m *MockConnection) SetWriteDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}
