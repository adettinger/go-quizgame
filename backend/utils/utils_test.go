package utils_test

import (
	"testing"

	"github.com/adettinger/go-quizgame/testutils"
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

func TestIsAlphanumeric(t *testing.T) {
	t.Run("letters and numbers", func(t *testing.T) {
		got := utils.IsAlphanumeric("abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		testutils.AssertTrue(t, got)
	})

	symbols := "!@#$%^&*()<>?:\"{}_+,.;'[]-=/|\\"
	for _, sym := range symbols {
		t.Run(string(sym), func(t *testing.T) {
			testutils.AssertFalse(t, utils.IsAlphanumeric(string(sym)))
		})
	}
}

func TestIsPlayerNameValid(t *testing.T) {
	cases := []struct {
		testName   string
		playerName string
		maxLength  int
		want       bool
	}{
		{"valid player name", "valid name 1", 100, true},
		{"begins with space", " player", 100, false},
		{"end with space", "player ", 100, false},
		{"is only spaces", "   ", 100, false},
		{"has special characters", "!", 100, false},
		{"longer than max", "hello", 1, false},
	}

	for _, tt := range cases {
		t.Run(tt.testName, func(t *testing.T) {
			got := utils.IsPlayerNameValid(tt.playerName, tt.maxLength)
			testutils.AssertEqual(t, got, tt.want)
		})
	}
}
