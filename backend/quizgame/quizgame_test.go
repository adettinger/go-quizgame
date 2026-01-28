package quizgame

import "testing"

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
			got := cleanInput(tt.input)
			if got != tt.want {
				t.Errorf("Got %q, expected %q", got, tt.want)
			}
		})
	}
}
