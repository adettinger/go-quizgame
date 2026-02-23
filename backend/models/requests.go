package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateProblemRequest struct {
	Type     string
	Question string
	Choices  []string
	Answer   string
}

type StartQuizResponse struct {
	SessionId uuid.UUID
	Timeout   time.Time
	Questions []Question
}

type QuestionSubmission struct {
	QuestionId uuid.UUID
	Answer     string
}

// TODO: Change request object
type EvaluateQuizRequest struct {
	SessionID uuid.UUID
	QuestionSubmissions []QuestionSubmission
}

type EvaluateQuizResponse struct {
	Score   int
	Answers []QuestionResponse
}

type QuestionResponse struct {
	Id      uuid.UUID
	Answer  string
	Correct bool
}
