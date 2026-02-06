package webserver

import (
	"fmt"

	"github.com/adettinger/go-quizgame/models"
)

type QuizService struct {
	ds *DataStore
}

func NewQuizService(ds *DataStore) *QuizService {
	return &QuizService{
		ds: ds,
	}
}

func (qs *QuizService) EvaluateQuiz(submission []models.Problem) (models.EvaluateQuizResponse, error) {
	questionResponses := make([]models.QuestionResponse, len(submission))
	score := 0
	for i, s := range submission {
		matchingProblem, err := qs.ds.GetProblemById(s.Id)
		if err != nil {
			return models.EvaluateQuizResponse{}, fmt.Errorf("Cannot find problem with Id %v", s.Id)
		}
		correct := s.Answer == matchingProblem.Answer
		questionResponses[i] = models.QuestionResponse{
			Id:      s.Id,
			Answer:  matchingProblem.Answer,
			Correct: correct,
		}
		if correct {
			score++
		}
	}
	return models.EvaluateQuizResponse{
		Score:   score,
		Answers: questionResponses,
	}, nil
}
