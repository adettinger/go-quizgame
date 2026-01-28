package quizgame

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"
)

type quizgame struct {
	problems []problem
	score    int
}

type problem struct {
	question string
	answer   string
}

func (p problem) String() string {
	return fmt.Sprintf("question: %v, answer: %v", p.question, p.answer)
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
	problems, err := parseProblems(fileName)
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
		fmt.Println(problem.question)
		answer := cleanInput(readLine(reader))
		if answer == problem.answer {
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

func parseProblems(fileName string) ([]problem, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open problems file. %v", err.Error())
	}
	reader := csv.NewReader(file)

	expectedFieldCount := reflect.TypeOf(problem{}).NumField()
	lineCount := 0
	problems := make([]problem, 0)
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("Failed to parse problems. %v", err.Error())
		}
		lineCount++
		if len(record) != expectedFieldCount {
			return nil, fmt.Errorf("Expected %d columns per row. Found %d on line %d", expectedFieldCount, len(record), lineCount)
		}
		problems = append(problems, problem{
			question: record[0],
			answer:   cleanInput(record[1]),
		})
	}
	if lineCount == 0 {
		return nil, errors.New("Expected to found at least 1 problem. Found 0")
	}

	return problems, nil
}

func cleanInput(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}
