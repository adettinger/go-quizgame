package webserver

import (
	"errors"
	"sync"

	"github.com/adettinger/go-quizgame/csvparser"
	"github.com/adettinger/go-quizgame/problem"
)

type DataStore struct {
	problems []problem.Problem
	mu       sync.RWMutex
}

func NewDataStore(fileName string) (*DataStore, error) {
	problems, err := csvparser.ParseProblems(fileName)
	if err != nil {
		return nil, err
	}
	return &DataStore{
		problems: problems,
		mu:       sync.RWMutex{},
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
	return nil
}

func (ds *DataStore) AddProblem(problem problem.Problem) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.problems = append(ds.problems, problem)
}
