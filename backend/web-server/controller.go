package webserver

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adettinger/go-quizgame/models"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	ds *DataStore
}

func NewProblemController(ds *DataStore) *Controller {
	return &Controller{ds: ds}
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
	c.IndentedJSON(http.StatusOK, wc.ds.ListProblems())
}

func (wc Controller) GetProblemByIndex(c *gin.Context) {
	index, err := parseIndex(c)
	if err != nil {
		return
	}
	problem, err := wc.ds.GetProblemByIndex(index)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Index does not exist"})
		return
	}

	c.IndentedJSON(http.StatusOK, problem)
}

func (wc Controller) DeleteProblem(c *gin.Context) {
	index, err := parseIndex(c)
	if err != nil {
		return
	}
	err = wc.ds.DeleteProblemByIndex(index)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Index does not exist"})
	}
	c.IndentedJSON(http.StatusNoContent, gin.H{"message": fmt.Sprintf("Deleted index %d", index)})
}

func (wc Controller) AddProblem(c *gin.Context) {
	var problemRequest models.CreateProblemRequest
	if err := c.BindJSON(&problemRequest); err != nil {
		return
	}
	fmt.Println("Bound problem request")

	problem := wc.ds.AddProblem(problemRequest)
	c.IndentedJSON(http.StatusCreated, problem)
}

func (wc Controller) SaveProblems(c *gin.Context) {
	err := wc.ds.SaveProblems()
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to save problems"})
		return
	}
	// TODO: Is this correct http code?
	c.IndentedJSON(http.StatusAccepted, gin.H{"message": "Saved problems"})
}

func parseIndex(c *gin.Context) (int, error) {
	input := c.Param("index")
	index, err := strconv.Atoi(input)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad index"})
		return -1, errors.New("Bad index")
	}
	return index, nil
}
