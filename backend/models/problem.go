package models

import (
	"encoding/json"
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
	if problemType == ProblemTypeChoice {
		if len(choices) < 2 || len(choices) > MaxNumChoices {
			return fmt.Errorf("Choice type must have at least 2 choices and at most %d choices", MaxNumChoices)
		}
		choiceFound := false
		for _, c := range choices {
			if c == "" {
				return fmt.Errorf("Choice cannot be empty string")
			}
			if c == answer {
				choiceFound = true
			}
		}
		if !choiceFound {
			return fmt.Errorf("Answer must be one of the choices")
		}
	}
	if problemType == ProblemTypeText && len(choices) != 0 {
		return fmt.Errorf("Text problems cannot have choices")
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
	return fmt.Sprintf("id: %v, question: %v, answer: %v", p.Id, p.Question, p.Answer)
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
	Question string
}

func (q Question) String() string {
	return fmt.Sprintf("id: %v, question: %v", q.Id, q.Question)
}
