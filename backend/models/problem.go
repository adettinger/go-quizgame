package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Problem struct {
	Id       uuid.UUID
	Question string
	Answer   string
}

func (p Problem) String() string {
	return fmt.Sprintf("id: %v, question: %v, answer: %v", p.Id, p.Question, p.Answer)
}

func (p Problem) ToStringSlice() []string {
	return []string{p.Id.String(), p.Question, p.Answer}
}

type Question struct {
	Id       uuid.UUID
	Question string
}

func (q Question) String() string {
	return fmt.Sprintf("id: %v, question: %v", q.Id, q.Question)
}
