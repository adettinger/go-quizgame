package main

import (
	"fmt"
	"os"

	"github.com/adettinger/go-quizgame/quizgame"
)

func main() {
	fmt.Println("Welcome to quizgame!")
	quizgame.StartGame(os.Stdin)
}
