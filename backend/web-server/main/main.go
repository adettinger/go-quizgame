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
	router.GET("/ping", controller.Ping)
	router.GET("/hello", controller.HelloWorld)
	router.GET("/listProblems", controller.ListProblems)
	router.GET("/problem/:index", controller.GetProblemByIndex)
	router.GET("/problem/delete/:index", controller.DeleteProblem)

	router.Run("localhost:8080")
}
