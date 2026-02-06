package models

import "github.com/google/uuid"

type CreateProblemRequest struct {
	Question string
	Answer   string
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
