package webserver

import (
	"errors"
	"sync"

	"github.com/adettinger/go-quizgame/csv"
	"github.com/adettinger/go-quizgame/problem"
)

type DataStore struct {
	fileName string
	problems []problem.Problem
	mu       sync.RWMutex
	modified bool
}

func NewDataStore(fileName string) (*DataStore, error) {
	problems, err := csv.ParseProblems(fileName)
	if err != nil {
		return nil, err
	}
	return &DataStore{
		fileName: fileName,
		problems: problems,
		mu:       sync.RWMutex{},
		modified: false,
	}, nil
}

func (ds *DataStore) ListProblems() []problem.Problem {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Return a copy to prevent external modification TODO: Needed?
	problemsCopy := make([]problem.Problem, len(ds.problems))
	copy(problemsCopy, ds.problems)
	return problemsCopy
}

func (ds *DataStore) GetProblemByIndex(index int) (problem.Problem, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if index < 0 || index > (len(ds.problems)-1) {
		return problem.Problem{}, errors.New("Index out of bounds")
	}

	return ds.problems[index], nil
}

func (ds *DataStore) DeleteProblemByIndex(index int) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if index < 0 || index > (len(ds.problems)-1) {
		return errors.New("Index out of bounds")
	}

	ds.problems = append(ds.problems[:index], ds.problems[index+1:]...)
	ds.modified = true
	return nil
}

func (ds *DataStore) AddProblem(problem problem.Problem) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.problems = append(ds.problems, problem)
	ds.modified = true
}

func (ds *DataStore) SaveProblems() error {
	if !ds.modified {
		return errors.New("No modifications to save")
	}
	err := csv.WriteProblems(ds.fileName, ds.problems)
	if err != nil {
		return err
	}
	return nil
}
