package main

import (
	"fmt"
	"os"

	"github.com/adettinger/go-quizgame/webserver"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting server...")
	ds, err := webserver.NewDataStore("../../problems.csv")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	controller := webserver.NewProblemController(ds)
	router := gin.Default()

	// Server health
	router.GET("/ping", controller.Ping)
	router.GET("/hello", controller.HelloWorld)

	// Problem related
	router.GET("/problem", controller.ListProblems)
	router.GET("/problem/:id", controller.GetProblemById)
	router.DELETE("/problem/:id", controller.DeleteProblem)
	router.POST("/problem", controller.AddProblem)
	router.POST("/problem/save", controller.SaveProblems)

	router.Run("localhost:8080")
}
