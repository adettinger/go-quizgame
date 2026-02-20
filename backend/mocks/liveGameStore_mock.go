package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockLiveGameStore is a mock implementation of the LiveGameStore
type MockLiveGameStore struct {
	mock.Mock
}

// RemovePlayerByName mocks the RemovePlayerByName method
func (m *MockLiveGameStore) RemovePlayerByName(name string) {
	m.Called(name)
}

// Add other methods from the LiveGameStore interface as needed
