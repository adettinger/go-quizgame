package csv

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

func WriteProblems(fileName string, problems []problem.Problem) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return errors.New("Failed to open file")
	}
	buffer := make([][]string, len(problems))
	for i, p := range problems {
		buffer[i] = p.ToStringSlice()
	}
	writer := csv.NewWriter(file)
	err = writer.WriteAll(buffer)
	if err != nil {
		return errors.New("Failed to write to file")
	}

	return nil
}
