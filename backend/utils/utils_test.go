package utils_test

import (
	"testing"

	"github.com/adettinger/go-quizgame/utils"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"to lower case", "ABC", "abc"},
		{"trim whitespace", " abc ", "abc"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.CleanInput(tt.input)
			if got != tt.want {
				t.Errorf("Got %q, expected %q", got, tt.want)
			}
		})
	}
}
