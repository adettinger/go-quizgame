package webserver

import (
	"strings"
	"time"

	"github.com/adettinger/go-quizgame/models"
	"github.com/adettinger/go-quizgame/types"
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

func (qs *QuizService) EvaluateQuiz(sessionId uuid.UUID, submission []models.QuestionSubmission) (models.EvaluateQuizResponse, error) {
	isActive, err := qs.ss.IsSessionActive(sessionId, time.Now())
	if err != nil {
		return models.EvaluateQuizResponse{}, &types.ErrSessionNotFound{SessionID: sessionId}
	}
	defer qs.ss.DeleteSession(sessionId) //Delete session after processing this function
	if !isActive {
		return models.EvaluateQuizResponse{}, &types.ErrSessionExpired{SessionID: sessionId}
	}

	questionResponses := make([]models.QuestionResponse, len(submission))
	score := 0
	for i, s := range submission {
		matchingProblem, err := qs.ds.GetProblemById(s.QuestionId)
		if err != nil {
			return models.EvaluateQuizResponse{}, &types.ErrProblemNotFound{ProblemId: s.QuestionId}
		}
		correct := strings.EqualFold(s.Answer, matchingProblem.Answer)
		questionResponses[i] = models.QuestionResponse{
			Id:      s.QuestionId,
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
