package main

import (
	"fmt"

	webserver "github.com/adettinger/go-quizgame/web-server"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting server...")
	router := gin.Default()
	router.GET("/ping", webserver.Ping)
	router.GET("/hello", webserver.HelloWorld)

	router.Run("localhost:8080")
}
