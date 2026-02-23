package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/adettinger/go-quizgame/models"
	"github.com/adettinger/go-quizgame/webserver"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProblemController struct {
	ds *webserver.DataStore
}

func NewProblemController(ds *webserver.DataStore) *ProblemController {
	return &ProblemController{
		ds: ds,
	}
}

func (wc ProblemController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (wc ProblemController) HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}

func (wc ProblemController) ListProblems(c *gin.Context) {
	c.JSON(http.StatusOK, wc.ds.ListProblems())
}

func (wc ProblemController) GetProblemById(c *gin.Context) {
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

func (wc ProblemController) DeleteProblem(c *gin.Context) {
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

func (wc ProblemController) AddProblem(c *gin.Context) {
	var problemRequest models.CreateProblemRequest
	if err := c.BindJSON(&problemRequest); err != nil || problemRequest.Question == "" || problemRequest.Answer == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	problem, err := wc.ds.AddProblem(problemRequest)
	if err != nil {
		log.Printf("Error adding problem: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	c.JSON(http.StatusCreated, problem)
}

// func (wc ProblemController) EditProblem(c *gin.Context) {
// 	var problemRequest models.EditProblemRequest
// 	if err := c.BindJSON(&problemRequest); err != nil || problemRequest.Question == "" || problemRequest.Answer == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
// 		return
// 	}
// 	fmt.Printf("ProblemRequest: %v\n", problemRequest)

// 	problem, err := wc.ds.GetProblemById(problemRequest.Id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "Id does not exist"})
// 		return
// 	}
// }

func (wc ProblemController) SaveProblems(c *gin.Context) {
	err := wc.ds.SaveProblems()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save problems"})
		return
	}
	// TODO: Is this correct http code?
	c.JSON(http.StatusAccepted, gin.H{"message": "Saved problems"})
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
