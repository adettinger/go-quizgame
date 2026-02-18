package utils

import (
	"regexp"
	"strings"
)

func CleanInput(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func IsAlphanumeric(s string) bool {
	// Create a regex pattern that matches strings containing only letters and numbers
	pattern := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return pattern.MatchString(s)
}

func IsPlayerNameValid(name string, maxLength int) bool {
	spacesRemovedName := strings.ReplaceAll(name, " ", "")

	return (name == strings.TrimSpace(name) && //Does not begin or end with space
		len(spacesRemovedName) != 0 && //is not only spaces
		IsAlphanumeric(spacesRemovedName) && //has only alphanumeric and spaces
		len(name) <= maxLength) //Fits in max length
}
