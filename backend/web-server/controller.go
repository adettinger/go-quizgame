package webserver

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	index, err := getIndex(c)
	if err != nil {
		return
	}
	problem, err := wc.ds.GetProblemByIndex(index)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "error retrieving problem"})
		return
	}

	c.IndentedJSON(http.StatusOK, problem)
}

func (wc Controller) DeleteProblem(c *gin.Context) {
	index, err := getIndex(c)
	if err != nil {
		return
	}
	err = wc.ds.DeleteProblemByIndex(index)
	c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Deleted index %d", index)})
}

func getIndex(c *gin.Context) (int, error) {
	input := c.Param("index")
	index, err := strconv.Atoi(input)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad index"})
		return -1, errors.New("Bad index")
	}
	return index, nil
}
