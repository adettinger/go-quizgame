package quizgame

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/adettinger/go-quizgame/csvparser"
	"github.com/adettinger/go-quizgame/problem"
	"github.com/adettinger/go-quizgame/utils"
)

type quizgame struct {
	problems []problem.Problem
	score    int
}

func QuizGame(in io.Reader, fileName string, timeLimit int, random bool) {
	game := setupGame(fileName, random)
	quizCompleted := make(chan bool, 1)
	go game.startGame(in, quizCompleted)
	select {
	case <-quizCompleted:
	case <-time.After(time.Duration(timeLimit) * time.Second):
		fmt.Println("Time's up!")
	}

	fmt.Printf("Final score: %d out of %d\n", game.score, len(game.problems))
}

func setupGame(fileName string, random bool) quizgame {
	problems, err := csvparser.ParseProblems(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Read %d problems\n", len(problems))

	if random {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}
	return quizgame{problems, 0}
}

func (qg *quizgame) startGame(in io.Reader, done chan<- bool) {
	reader := bufio.NewScanner(in)
	for _, problem := range qg.problems {
		fmt.Println(problem.Question)
		answer := utils.CleanInput(readLine(reader))
		if answer == problem.Answer {
			fmt.Println("Correct!")
			qg.score++
		} else {
			fmt.Println("Wrong Answer!")
		}
	}
	done <- true
}

func readLine(in *bufio.Scanner) string {
	in.Scan()
	return in.Text()
}
