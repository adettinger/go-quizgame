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

func NewDataStoreFromData(problems []models.Problem) (*DataStore, error) {
	return &DataStore{
		fileName: "ignore",
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

// TODO: Move validation into business logic component
func (ds *DataStore) AddProblem(pr models.CreateProblemRequest) (models.Problem, error) {
	problemType, err := models.ParseProblemType(pr.Type)
	if err != nil {
		return models.Problem{}, err
	}
	if err = models.ValidateChoices(problemType, pr.Choices, pr.Answer); err != nil {
		return models.Problem{}, err
	}
	if pr.Question == "" {
		return models.Problem{}, errors.New("Question cannot be empty string")
	}
	if pr.Answer == "" {
		return models.Problem{}, errors.New("Answer canot be empty string")
	}

	problem := models.Problem{
		Id:       ds.getNewId(),
		Type:     problemType,
		Question: pr.Question,
		Choices:  pr.Choices,
		Answer:   pr.Answer,
	}
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.problems = append(ds.problems, problem)
	ds.modified = true
	return problem, nil
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

func (ds *DataStore) GetQuestions() []models.Question {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Return a copy to prevent external modification TODO: Needed?
	questions := make([]models.Question, len(ds.problems))
	for i, p := range ds.problems {
		questions[i] = models.Question{
			Id:       p.Id,
			Type:     p.Type,
			Question: p.Question,
			Choices:  p.Choices,
		}
	}
	return questions
}

func (ds *DataStore) getNewId() uuid.UUID {
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
