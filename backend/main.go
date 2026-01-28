package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adettinger/go-quizgame/quizgame"
)

func main() {
	timeLimit := flag.Int("time", 5, "time limit in seconds")
	shuffleOder := flag.Bool("random", false, "should the question be random order")
	flag.Parse()

	fmt.Println("Welcome to quizgame!")
	quizgame.QuizGame(os.Stdin, *timeLimit, *shuffleOder)
}
