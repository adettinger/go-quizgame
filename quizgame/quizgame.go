package quizgame

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
)

type problem struct {
	question string
	answer   string
}

func (p problem) String() string {
	return fmt.Sprintf("question: %v, answer: %v", p.question, p.answer)
}

// TODO: Take out as a param (dependency inject)
func StartGame(in io.Reader) {
	reader := bufio.NewScanner(in)
	problems, err := parseProblems("problems.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Read %d problems", len(problems))

	for _, problem := range problems {
		fmt.Println(problem.String())
	}

	score := 0
	for _, problem := range problems {
		fmt.Println(problem.question)
		answer := readLine(reader)
		if answer == problem.answer {
			fmt.Println("Correct!")
			score++
		} else {
			fmt.Println("Wrong Answer!")
		}
	}
	fmt.Printf("Final score: %d out of %d\n", score, len(problems))
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
			answer:   record[1],
		})
	}
	if lineCount == 0 {
		return nil, errors.New("Expected to found at least 1 problem. Found 0")
	}

	return problems, nil
}
