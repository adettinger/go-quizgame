package main

import (
	"fmt"
	"os"

	webserver "github.com/adettinger/go-quizgame/web-server"
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
	router.GET("/problem/:index", controller.GetProblemByIndex)
	router.DELETE("/problem/:index", controller.DeleteProblem)
	router.POST("/problem", controller.AddProblem)
	router.POST("/")

	router.Run("localhost:8080")
}
