package webserver

import (
	"fmt"

	"github.com/adettinger/go-quizgame/models"
	"github.com/google/uuid"
)

type quizService struct {
	ds *DataStore
}

type ProblemEvaluation struct {
	Id      uuid.UUID
	Answer  string
	Correct bool
}

func (qs quizService) evaluateQuiz(submission []models.Problem) ([]ProblemEvaluation, error) {
	response := make([]ProblemEvaluation, len(submission))
	for i, s := range submission {
		matchingProblem, err := qs.ds.GetProblemById(s.Id)
		if err != nil {
			return []ProblemEvaluation{}, fmt.Errorf("Cannot find problem with Id %v", s.Id)
		}
		response[i] = ProblemEvaluation{
			Id:      s.Id,
			Answer:  matchingProblem.Answer,
			Correct: s.Answer == matchingProblem.Answer,
		}
	}
	return response, nil
}
