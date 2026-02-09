package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateProblemRequest struct {
	Question string
	Answer   string
}

type StartQuizResponse struct {
	SessionId uuid.UUID
	Timeout   time.Time
	Questions []Question
}

type EvaluateQuizRequest struct {
	SessionID uuid.UUID
	Questions []Problem
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
