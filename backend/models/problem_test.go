package models_test

import (
	"strconv"
	"testing"

	"github.com/adettinger/go-quizgame/models"
	"github.com/adettinger/go-quizgame/testutils"
)

func TestParseProblemType(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		pType   models.ProblemType
		isValid bool
	}{
		{
			"Text type",
			"Text",
			models.ProblemTypeText,
			true,
		},
		{
			"Choice type",
			"Choice",
			models.ProblemTypeChoice,
			true,
		},
		{
			"Invalid type",
			"invalid",
			models.ProblemTypeText,
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := models.ParseProblemType(tt.input)
			if tt.isValid {
				testutils.AssertNoError(t, err)
				testutils.AssertEqual(t, got, tt.pType)
			} else {
				testutils.AssertHasError(t, err)
			}
		})
	}
}

func TestValidateChoices(t *testing.T) {
	cases := []struct {
		name    string
		pType   models.ProblemType
		choices []string
		answer  string
		isValid bool
	}{
		{
			"text type",
			models.ProblemTypeText,
			[]string{},
			"nala",
			true,
		},
		{
			"answer empty string",
			models.ProblemTypeText,
			[]string{},
			"",
			false,
		},
		{
			"Text with choices",
			models.ProblemTypeText,
			[]string{"Nala", "Hasse", "Murphey"},
			"nala",
			false,
		},
		{
			"choice type",
			models.ProblemTypeChoice,
			[]string{"Nala", "Hasse", "Murphey"},
			"nala",
			true,
		},
		{
			"Only 1 choice",
			models.ProblemTypeChoice,
			[]string{"Nala"},
			"nala",
			false,
		},
		{
			"too many choices",
			models.ProblemTypeChoice,
			func() []string {
				toReturn := []string{}
				for i := range models.MaxNumChoices + 1 {
					toReturn = append(toReturn, strconv.Itoa(i))
				}
				return toReturn
			}(),
			"nala",
			false,
		},
		{
			"Empty choice string",
			models.ProblemTypeChoice,
			[]string{"Nala", "", "Murphey"},
			"nala",
			false,
		},
		{
			"duplicate choice",
			models.ProblemTypeChoice,
			[]string{"Nala", "naLa", "Murphey"},
			"nala",
			false,
		},
		{
			"Answer not one of choices",
			models.ProblemTypeChoice,
			[]string{"Nala", "Hasse", "Murphey"},
			"Frankie",
			false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := models.ValidateChoices(tt.pType, tt.choices, tt.answer)
			switch tt.isValid {
			case true:
				testutils.AssertNoError(t, err)
			case false:
				testutils.AssertHasError(t, err)
			}
		})
	}
}
