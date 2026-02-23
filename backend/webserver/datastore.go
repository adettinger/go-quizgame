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
	problems map[uuid.UUID]models.Problem
	mu       sync.RWMutex
	modified bool
}

func NewDataStore(fileName string) (*DataStore, error) {
	problems, err := csv.ParseProblems(fileName)
	if err != nil {
		return nil, err
	}
	problemsMap := make(map[uuid.UUID]models.Problem, len(problems))
	for _, p := range problems {
		problemsMap[p.Id] = p
	}

	return &DataStore{
		fileName: fileName,
		problems: problemsMap,
		mu:       sync.RWMutex{},
		modified: false,
	}, nil
}

func NewDataStoreFromData(problems []models.Problem) (*DataStore, error) {
	problemsMap := make(map[uuid.UUID]models.Problem, len(problems))
	for _, p := range problems {
		problemsMap[p.Id] = p
	}

	return &DataStore{
		fileName: "ignore",
		problems: problemsMap,
		mu:       sync.RWMutex{},
		modified: false,
	}, nil
}

func (ds *DataStore) ListProblems() []models.Problem {
	return ds.problemsToArray()
}

func (ds *DataStore) GetProblemById(uuid uuid.UUID) (models.Problem, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	problem, ok := ds.problems[uuid]
	if !ok {
		return models.Problem{}, errors.New("Problem not found")
	}
	return problem, nil
}

func (ds *DataStore) DeleteProblemByIndex(id uuid.UUID) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	_, ok := ds.problems[id]
	if !ok {
		return errors.New("Id not found")
	}
	delete(ds.problems, id)
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

	ds.problems[problem.Id] = problem
	ds.modified = true
	return problem, nil
}

// func (ds *DataStore) Edit(pr models.EditProblemRequest) (models.Problem, error) {
// 	problemType, err := models.ParseProblemType(pr.Type)
// 	if err != nil {
// 		return models.Problem{}, err
// 	}
// 	if err = models.ValidateChoices(problemType, pr.Choices, pr.Answer); err != nil {
// 		return models.Problem{}, err
// 	}
// 	if pr.Question == "" {
// 		return models.Problem{}, errors.New("Question cannot be empty string")
// 	}
// 	if pr.Answer == "" {
// 		return models.Problem{}, errors.New("Answer canot be empty string")
// 	}

// 	problem := models.Problem{
// 		Id:       pr.Id,
// 		Type:     problemType,
// 		Question: pr.Question,
// 		Choices:  pr.Choices,
// 		Answer:   pr.Answer,
// 	}
// 	ds.mu.Lock()
// 	defer ds.mu.Unlock()

// 	ds.problems = append(ds.problems, problem)
// 	ds.modified = true
// 	return problem, nil
// }

func (ds *DataStore) SaveProblems() error {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if !ds.modified {
		return errors.New("No modifications to save")
	}
	err := csv.WriteProblems(ds.fileName, ds.problemsToArray())
	if err != nil {
		return err
	}
	return nil
}

func (ds *DataStore) GetQuestions() []models.Question {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Return a copy to prevent external modification TODO: Needed?
	questions := make([]models.Question, 0, len(ds.problems))
	index := 0
	for _, p := range ds.problems {
		questions[index] = models.Question{
			Id:       p.Id,
			Type:     p.Type,
			Question: p.Question,
			Choices:  p.Choices,
		}
		index++
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

func (ds *DataStore) problemsToArray() []models.Problem {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	problems := make([]models.Problem, 0, len(ds.problems))

	for _, problem := range ds.problems {
		problems = append(problems, problem)
	}
	return problems
}
