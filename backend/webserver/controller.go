package webserver

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adettinger/go-quizgame/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	ds *DataStore
	qs *quizService
}

func NewProblemController(ds *DataStore) *Controller {
	return &Controller{
		ds: ds,
		qs: &quizService{ds: ds},
	}
}

func (wc Controller) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (wc Controller) HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}

func (wc Controller) ListProblems(c *gin.Context) {
	c.JSON(http.StatusOK, wc.ds.ListProblems())
}

func (wc Controller) GetProblemById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID"})
		return
	}
	problem, err := wc.ds.GetProblemById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Id does not exist"})
		return
	}

	c.JSON(http.StatusOK, problem)
}

func (wc Controller) DeleteProblem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid UUID"})
		return
	}
	err = wc.ds.DeleteProblemByIndex(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Index does not exist"})
	}
	c.JSON(http.StatusNoContent, struct{}{})
}

func (wc Controller) AddProblem(c *gin.Context) {
	var problemRequest models.CreateProblemRequest
	if err := c.BindJSON(&problemRequest); err != nil || problemRequest.Question == "" || problemRequest.Answer == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	problem := wc.ds.AddProblem(problemRequest)
	c.JSON(http.StatusCreated, problem)
}

func (wc Controller) SaveProblems(c *gin.Context) {
	err := wc.ds.SaveProblems()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save problems"})
		return
	}
	// TODO: Is this correct http code?
	c.JSON(http.StatusAccepted, gin.H{"message": "Saved problems"})
}

func (wc Controller) GetQuestions(c *gin.Context) {
	c.JSON(http.StatusOK, wc.ds.GetQuestions())
}

func (wc Controller) SubmitQuiz(c *gin.Context) {
	var request = []models.Problem{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	// Warning: Sending service error directly to frontend
	fmt.Printf("request: %v\n", request)
	response, err := wc.qs.evaluateQuiz(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	fmt.Printf("Eval response: %v\n", response)
	c.JSON(http.StatusOK, response)
}

func parseIndex(c *gin.Context) (int, error) {
	input := c.Param("index")
	index, err := strconv.Atoi(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad index"})
		return -1, errors.New("Bad index")
	}
	return index, nil
}
