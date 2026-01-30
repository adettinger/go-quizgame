package mocks

// This is currently unused

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/adettinger/go-quizgame/models"
)

type MockDataStore struct {
	mock.Mock
	problems []models.Problem
}

func (m *MockDataStore) ListProblems() []models.Problem {
	return m.problems
}

func (m *MockDataStore) GetProblemById(id uuid.UUID) (models.Problem, error) {
	args := m.Called(id)
	return args.Get(0).(models.Problem), args.Error(1)
}
