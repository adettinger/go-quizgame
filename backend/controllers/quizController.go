package controllers

import (
	"net/http"
	"time"

	"github.com/adettinger/go-quizgame/models"
	"github.com/adettinger/go-quizgame/webserver"
	"github.com/gin-gonic/gin"
)

const QuizTimeout = time.Duration(5 * time.Minute)

type QuizController struct {
	ds *webserver.DataStore
	qs *webserver.QuizService
	ss *webserver.SessionStore
}

// Note: Instantiates a quizService
func NewQuizController(ds *webserver.DataStore) *QuizController {
	return &QuizController{
		ds: ds,
		qs: webserver.NewQuizService(ds),
		ss: webserver.NewSessionStore(),
	}
}

func (qc QuizController) GetQuestions(c *gin.Context) {
	c.JSON(http.StatusOK, qc.ds.GetQuestions())
}

func (qc QuizController) StartQuiz(c *gin.Context) {
	session := qc.ss.CreateSession(QuizTimeout)
	c.JSON(http.StatusOK, models.StartQuizResponse{
		SessionId: session.Id,
		Timeout:   session.Timeout,
		Questions: qc.ds.GetQuestions(),
	})
}

func (qc QuizController) SubmitQuiz(c *gin.Context) {
	var request = []models.Problem{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	response, err := qc.qs.EvaluateQuiz(request)
	if err != nil {
		// Warning: Sending service error directly to frontend
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	c.JSON(http.StatusOK, response)
}
