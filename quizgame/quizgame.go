package quizgame

import (
	"encoding/csv"
	"fmt"
	"os"
)

// var csv_headers = `problem, solution`
func StartGame() {
	problems, err := parseProblems("problems.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Read %d problems", len(problems))
}

func parseProblems(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open problems file. %v", err.Error())
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to parse problems. %v", err.Error())
	}
	// TODO: Check records matches expected format
	return records, nil
}
