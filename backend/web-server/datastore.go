package webserver

import (
	"errors"
	"sync"

	"github.com/adettinger/go-quizgame/csv"
	"github.com/adettinger/go-quizgame/models"
	"github.com/google/uuid"
)

type DataStore struct {
	fileName string
	problems []models.Problem
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

func (ds *DataStore) ListProblems() []models.Problem {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Return a copy to prevent external modification TODO: Needed?
	problemsCopy := make([]models.Problem, len(ds.problems))
	copy(problemsCopy, ds.problems)
	return problemsCopy
}

func (ds *DataStore) GetProblemByIndex(index int) (models.Problem, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if index < 0 || index > (len(ds.problems)-1) {
		return models.Problem{}, errors.New("Index out of bounds")
	}

	return ds.problems[index], nil
}

func (ds *DataStore) DeleteProblemByIndex(id uuid.UUID) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	foundIndex := -1
	for i, p := range ds.problems {
		if p.Id == id {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return errors.New("Id not found")
	}

	ds.problems = append(ds.problems[:foundIndex], ds.problems[foundIndex+1:]...)
	ds.modified = true
	return nil
}

func (ds *DataStore) AddProblem(pr models.CreateProblemRequest) models.Problem {
	problem := models.Problem{
		Id:       ds.GetNewId(),
		Question: pr.Question,
		Answer:   pr.Answer,
	}
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.problems = append(ds.problems, problem)
	ds.modified = true
	return problem
}

func (ds *DataStore) SaveProblems() error {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if !ds.modified {
		return errors.New("No modifications to save")
	}
	err := csv.WriteProblems(ds.fileName, ds.problems)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DataStore) GetNewId() uuid.UUID {
	for {
		uuid := uuid.New()
		if !ds.problemIdExists(uuid) {
			return uuid
		}
	}
}

func (ds *DataStore) problemIdExists(uuid uuid.UUID) bool {
	_, err := ds.GetProblemById(uuid)
	return err == nil
}

func (ds *DataStore) GetProblemById(uuid uuid.UUID) (models.Problem, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	for _, p := range ds.problems {
		if p.Id == uuid {
			return p, nil
		}
	}
	return models.Problem{}, errors.New("Problem not found")
}
