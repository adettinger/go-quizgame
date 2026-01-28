package utils

import "strings"

func CleanInput(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}
