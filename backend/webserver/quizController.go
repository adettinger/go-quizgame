package webserver

import (
	"net/http"

	"github.com/adettinger/go-quizgame/models"
	"github.com/gin-gonic/gin"
)

type QuizController struct {
	ds *DataStore
	qs *quizService
}

// Note: Instantiates a quizService
func NewQuizController(ds *DataStore) *QuizController {
	return &QuizController{
		ds: ds,
		qs: &quizService{ds: ds},
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

	response, err := qc.qs.evaluateQuiz(request)
	if err != nil {
		// Warning: Sending service error directly to frontend
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	c.JSON(http.StatusOK, response)
}
