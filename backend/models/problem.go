package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
)

const MaxNumChoices = 4

type ProblemType string

const (
	ProblemTypeText   ProblemType = "text"
	ProblemTypeChoice ProblemType = "choice"
)

func (pt ProblemType) String() string {
	return string(pt)
}

func (pt ProblemType) IsValid() bool {
	switch pt {
	case ProblemTypeText, ProblemTypeChoice:
		return true
	}
	return false
}

func ParseProblemType(s string) (ProblemType, error) {
	pt := ProblemType(strings.ToLower(s))
	if !pt.IsValid() {
		return "", fmt.Errorf("invalid problem type: %s", s)
	}
	return pt, nil
}

func ValidateChoices(problemType ProblemType, choices []string, answer string) error {
	if answer == "" {
		return errors.New("Answer cannot be empty string")
	}

	switch problemType {
	case ProblemTypeChoice:
		if len(choices) < 2 || len(choices) > MaxNumChoices {
			return fmt.Errorf("Choice type must have at least 2 choices and at most %d choices", MaxNumChoices)
		}
		choiceFound := false
		seen := make(map[string]struct{}, len(choices))
		for _, c := range choices {
			if c == "" {
				return errors.New("Choice cannot be empty string")
			}
			if _, exists := seen[strings.ToLower(c)]; exists {
				return fmt.Errorf("Duplicate choice found")
			}
			seen[strings.ToLower(c)] = struct{}{}
			if strings.EqualFold(c, answer) {
				choiceFound = true
			}
		}
		if !choiceFound {
			return fmt.Errorf("Answer must be one of the choices")
		}
	case ProblemTypeText:
		if len(choices) != 0 {
			return fmt.Errorf("Text problems cannot have choices")
		}
	default:
		return fmt.Errorf("Invalid problem type; %v", problemType)
	}
	return nil
}

type Problem struct {
	Id       uuid.UUID
	Type     ProblemType
	Question string
	Choices  []string
	Answer   string
}

func (p Problem) String() string {
	return fmt.Sprintf("id: %v, type: %v, question: %v, choices: %v, answer: %v", p.Id, p.Type.String(), p.Question, p.Choices, p.Answer)
}

func (p Problem) ToStringSlice() []string {
	return []string{p.Id.String(), p.Type.String(), p.Question, serializeArray(p.Choices), p.Answer}
}

func (p Problem) Equal(b Problem) bool {
	if p.Id != b.Id || p.Type != b.Type || p.Question != b.Question || !slices.Equal(p.Choices, b.Choices) || p.Answer != b.Answer {
		return false
	}
	return true
}

func serializeArray(arr []string) string {
	bytes, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

type Question struct {
	Id       uuid.UUID
	Type     ProblemType
	Question string
	Choices  []string
}

func (q Question) String() string {
	return fmt.Sprintf("id: %v, type: %v, question: %v, choices: %v", q.Id, q.Type.String(), q.Question, q.Choices)
}
