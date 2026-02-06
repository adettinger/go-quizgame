package controllers

import (
	"net/http"

	"github.com/adettinger/go-quizgame/models"
	"github.com/adettinger/go-quizgame/webserver"
	"github.com/gin-gonic/gin"
)

type QuizController struct {
	ds *webserver.DataStore
	qs *webserver.QuizService
}

// Note: Instantiates a quizService
func NewQuizController(ds *webserver.DataStore) *QuizController {
	return &QuizController{
		ds: ds,
		qs: webserver.NewQuizService(ds),
	}
}

func (qc QuizController) GetQuestions(c *gin.Context) {
	c.JSON(http.StatusOK, qc.ds.GetQuestions())
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
