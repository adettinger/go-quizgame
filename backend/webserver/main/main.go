package main

import (
	"fmt"
	"os"
	"time"

	"github.com/adettinger/go-quizgame/controllers"
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

	problemController := controllers.NewProblemController(ds)
	quizController := controllers.NewQuizController(ds)
	wsController := controllers.NewWebSocketController()

	go wsController.GetManager().Start()

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
	router.GET("/ping", problemController.Ping)
	router.GET("/hello", problemController.HelloWorld)

	// Problem endpoints
	router.GET("/problem", problemController.ListProblems)
	router.GET("/problem/:id", problemController.GetProblemById)
	router.DELETE("/problem/:id", problemController.DeleteProblem)
	router.POST("/problem", problemController.AddProblem)
	router.POST("/problem/save", problemController.SaveProblems)

	// Quiz endpoints
	router.GET("/quiz/questions", quizController.GetQuestions)
	router.GET("/quiz/start", quizController.StartQuiz)
	router.POST("/quiz/submit", quizController.SubmitQuiz)

	router.GET("/liveGame/player/:playerName", wsController.HandleConnection)

	router.Run("localhost:8080")
}
