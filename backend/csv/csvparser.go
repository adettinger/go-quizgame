package csv

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/adettinger/go-quizgame/models"
	"github.com/adettinger/go-quizgame/utils"
	"github.com/google/uuid"
)

// String Problem: ID, string, question, , answer
// Choice problem: Id, choice, question, choices[], answer

func ParseProblems(fileName string) ([]models.Problem, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open problems file. %v", err.Error())
	}
	defer file.Close()
	reader := csv.NewReader(file)

	expectedFieldCount := reflect.TypeOf(models.Problem{}).NumField()
	lineCount := 0
	problems := make([]models.Problem, 0)
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
		// Validate fields
		id, err := uuid.Parse(record[0])
		if err != nil {
			return nil, fmt.Errorf("Failed to parse UUID for line %d", lineCount)
		}
		questionType, err := models.ParseProblemType(record[1])
		if err != nil {
			return nil, fmt.Errorf("Failed to parse problem type for line %d", lineCount)
		}
		answer := utils.CleanInput(record[4])
		if answer == "" {
			return nil, fmt.Errorf("Answer cannot be empty string for line %d", lineCount)
		}
		choices, err := deserializeArray(record[3])
		if err != nil {
			return nil, fmt.Errorf("Failed to parse choices for line %d", lineCount)
		}
		if err := models.ValidateChoices(questionType, choices, answer); err != nil {
			return nil, fmt.Errorf("Line %d: %v", lineCount, err.Error())
		}

		problems = append(problems, models.Problem{
			Id:       id,
			Type:     questionType,
			Question: record[2],
			Choices:  choices,
			Answer:   answer,
		})
	}
	if lineCount == 0 {
		return nil, errors.New("Expected to found at least 1 problem. Found 0")
	}

	return problems, nil
}

func WriteProblems(fileName string, problems []models.Problem) error {
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

func deserializeArray(jsonStr string) ([]string, error) {
	var arr []string
	err := json.Unmarshal([]byte(jsonStr), &arr)
	return arr, err
}
