package webserver

import (
	"errors"
	"fmt"
	"time"

	"github.com/adettinger/go-quizgame/models"
	"github.com/google/uuid"
)

type QuizService struct {
	ds *DataStore
	ss *SessionStore
}

func NewQuizService(ds *DataStore, ss *SessionStore) *QuizService {
	return &QuizService{
		ds: ds,
		ss: ss,
	}
}

func (qs *QuizService) EvaluateQuiz(sessionId uuid.UUID, submission []models.Problem) (models.EvaluateQuizResponse, error) {
	isActive, err := qs.ss.IsSessionActive(sessionId, time.Now())
	if err != nil {
		return models.EvaluateQuizResponse{}, errors.New("Cannot find session")
	}
	defer qs.ss.DeleteSession(sessionId) //Delete session after processing this function
	if !isActive {
		return models.EvaluateQuizResponse{}, errors.New("Session is expired")
	}

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
