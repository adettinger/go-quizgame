package problem

import "fmt"

type Problem struct {
	Question string
	Answer   string
}

func (p Problem) String() string {
	return fmt.Sprintf("question: %v, answer: %v", p.Question, p.Answer)
}
