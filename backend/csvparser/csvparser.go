package csvparser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/adettinger/go-quizgame/problem"
	"github.com/adettinger/go-quizgame/utils"
)

func ParseProblems(fileName string) ([]problem.Problem, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open problems file. %v", err.Error())
	}
	defer file.Close()
	reader := csv.NewReader(file)

	expectedFieldCount := reflect.TypeOf(problem.Problem{}).NumField()
	lineCount := 0
	problems := make([]problem.Problem, 0)
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
		problems = append(problems, problem.Problem{
			Question: record[0],
			Answer:   utils.CleanInput(record[1]),
		})
	}
	if lineCount == 0 {
		return nil, errors.New("Expected to found at least 1 problem. Found 0")
	}

	return problems, nil
}
