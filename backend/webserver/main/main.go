package main

import (
	"fmt"
	"os"
	"time"

	"github.com/adettinger/go-quizgame/webserver"
	"github.com/gin-contrib/cors"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Server health
	router.GET("/ping", controller.Ping)
	router.GET("/hello", controller.HelloWorld)

	// Problem related
	router.GET("/problem", controller.ListProblems)
	router.GET("/problem/:id", controller.GetProblemById)
	router.DELETE("/problem/:id", controller.DeleteProblem)
	router.POST("/problem", controller.AddProblem)
	router.POST("/problem/save", controller.SaveProblems)

	router.GET("/quiz/questions", controller.GetQuestions)
	router.POST("/quiz/submit", controller.SubmitQuiz)

	router.Run("localhost:8080")
}
