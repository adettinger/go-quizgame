package webserver

import (
	"sync"
	"testing"

	"github.com/adettinger/go-quizgame/models"
	"github.com/google/uuid"
)

// TODO: Update these tests
var problems = map[uuid.UUID]models.Problem{
	uuid.MustParse("c620af48-3af0-4216-a229-65c539a00202"): models.Problem{Id: uuid.MustParse("c620af48-3af0-4216-a229-65c539a00202"), Type: "text", Question: "1+2", Choices: []string{}, Answer: "3"},
	uuid.MustParse("60d1584a-9d09-4e2d-be5c-1150fafa454f"): models.Problem{Id: uuid.MustParse("60d1584a-9d09-4e2d-be5c-1150fafa454f"), Type: "text", Question: "2*2", Choices: []string{}, Answer: "4"},
}

var uuids []uuid.UUID = []uuid.UUID{uuid.MustParse("c620af48-3af0-4216-a229-65c539a00202"), uuid.MustParse("60d1584a-9d09-4e2d-be5c-1150fafa454f")}

var ds = DataStore{
	fileName: "test",
	problems: problems,
	mu:       sync.RWMutex{},
	modified: false,
}

func TestProblemIdExists(t *testing.T) {
	t.Run("UUID exists", func(t *testing.T) {
		result := ds.problemIdExists(uuids[0])
		AssertEquals(t, result, true)
	})

	t.Run("UUID does not exist", func(t *testing.T) {
		result := ds.problemIdExists(uuid.MustParse("c620af48-3af0-4216-a229-65c539a00000"))
		AssertEquals(t, result, false)
	})
}

func TestGetProblemByID(t *testing.T) {
	t.Run("Find UUID exists", func(t *testing.T) {
		result, err := ds.GetProblemById(uuids[0])
		AssertNoError(t, err)
		AssertEquals(t, uuids[0], result.Id)
	})

	t.Run("UUID does not exist", func(t *testing.T) {
		_, err := ds.GetProblemById(uuid.MustParse("c620af48-3af0-4216-a229-65c539a00000"))
		AssertError(t, err)
	})
}

func AssertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("Expected no error. Found %v", got)
	}
}

func AssertError(t testing.TB, got error) {
	t.Helper()
	if got == nil {
		t.Fatal("Expected error. Found none")
	}
}

func AssertEquals[T comparable](t testing.TB, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("Expected got = want. Got: %v Want %v", got, want)
	}
}
