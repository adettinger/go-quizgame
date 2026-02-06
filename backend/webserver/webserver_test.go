package webserver_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/adettinger/go-quizgame/models"
	"github.com/adettinger/go-quizgame/webserver"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var problemSet []models.Problem = []models.Problem{
	{Id: uuid.MustParse("c620af48-3af0-4216-a229-65c539a00202"), Question: "1+2", Answer: "3"},
	{Id: uuid.MustParse("60d1584a-9d09-4e2d-be5c-1150fafa454f"), Question: "2*2", Answer: "4"},
}

func TestGetProblemById(t *testing.T) {
	testDataStore, _ := webserver.NewDataStoreFromData(problemSet)
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Test cases
	testCases := []struct {
		name           string
		urlParam       string
		problem        models.Problem
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Success",
			urlParam:       "c620af48-3af0-4216-a229-65c539a00202",
			problem:        problemSet[0],
			expectedStatus: http.StatusOK,
			expectedBody:   problemSet[0],
		},
		{
			name:           "Invalid UUID",
			urlParam:       "invalid-uuid",
			problem:        models.Problem{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"message": "Invalid UUID"},
		},
		{
			name:           "Problem Not Found",
			urlParam:       "123e4567-e89b-12d3-a456-426614174001",
			problem:        models.Problem{},
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"message": "Id does not exist"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			problemController := webserver.NewProblemController(testDataStore)

			// Create a test HTTP response recorder and context
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			// Set up the route
			r.GET("/problem/:id", problemController.GetProblemById)

			// Create the request
			req, _ := http.NewRequest("GET", "/problem/"+tt.urlParam, nil)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: tt.urlParam}}

			// Call the handler
			problemController.GetProblemById(c)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Assert response body
			if tt.expectedBody != nil {
				var actualBody interface{}
				var unmarshalError error = nil
				if tt.expectedStatus == http.StatusOK {
					var problem models.Problem
					unmarshalError = json.Unmarshal(w.Body.Bytes(), &problem)
					actualBody = problem
				} else {
					var response gin.H
					unmarshalError = json.Unmarshal(w.Body.Bytes(), &response)
					actualBody = response
				}

				assert.NoError(t, unmarshalError)
				assert.Equal(t, tt.expectedBody, actualBody)
			}
		})
	}
}

func TestListProblems(t *testing.T) {
	testDataStore, _ := webserver.NewDataStoreFromData(problemSet)
	gin.SetMode(gin.TestMode)
	problemController := webserver.NewProblemController(testDataStore)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	// Set up the route
	r.GET("/problem/", problemController.ListProblems)

	// Create the request
	req, _ := http.NewRequest("GET", "/problem/", nil)
	c.Request = req

	// Call the handler
	problemController.ListProblems(c)

	// Assert status code
	assert.Equal(t, http.StatusOK, w.Code)

	var responseBody []models.Problem
	unmarshalError := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, unmarshalError)
	if !slices.Equal(problemSet, responseBody) {
		t.Fatalf("Expected resonse to equal problem set. Got %v", responseBody)
	}
}

func TestDeleteProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name         string
		urlParam     string
		validUUID    bool
		foundProblem bool
		expectedBody gin.H
	}{
		{
			"Delete problem",
			"c620af48-3af0-4216-a229-65c539a00202",
			true,
			true,
			gin.H{},
		},
		{
			"Attempt to delete missing problem",
			"c620af48-3af0-4216-a229-65c539a00000",
			true,
			false,
			gin.H{"message": "Index does not exist"},
		},
		{
			"Invalid UUID",
			"123",
			false,
			false,
			gin.H{"message": "Invalid UUID"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			testDataStore, _ := webserver.NewDataStoreFromData(problemSet)
			problemController := webserver.NewProblemController(testDataStore)

			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			// Set up the route
			r.DELETE("/problem/", problemController.ListProblems)

			// Create the request
			req, _ := http.NewRequest("DELETE", "/problem/"+tt.urlParam, nil)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: tt.urlParam}}

			// Call the handler
			problemController.DeleteProblem(c)

			// Assert status code
			if !tt.validUUID {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else if !tt.foundProblem {
				assert.Equal(t, http.StatusNotFound, w.Code)
			} else {
				assert.Equal(t, http.StatusNoContent, w.Code)
			}

			// Assert response body
			var response gin.H
			json.Unmarshal(w.Body.Bytes(), &response)
			if !tt.validUUID || !tt.foundProblem {
				assert.Equal(t, tt.expectedBody, response)
			}

			if tt.validUUID && tt.foundProblem {
				_, err := testDataStore.GetProblemById(uuid.MustParse(tt.urlParam))
				if err == nil {
					t.Fatal("Expected problem to be deleted")
				}
			}

		})
	}
}

func TestAddProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cases := []struct {
		name           string
		requestBody    map[string]interface{}
		validRequest   bool
		expectedStatus int
	}{
		{
			"Create problem",
			map[string]interface{}{
				"Question": "Test question",
				"Answer":   "Test answer",
			},
			true,
			http.StatusCreated,
		},
		{
			"Invalid request",
			map[string]interface{}{},
			false,
			http.StatusBadRequest,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			testDataStore, _ := webserver.NewDataStoreFromData(problemSet)
			problemController := webserver.NewProblemController(testDataStore)

			w := httptest.NewRecorder()
			router := gin.New()

			// Set up the route
			router.POST("/problem/", problemController.AddProblem)

			// Create the request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/problem/", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check response body
			if !tt.validRequest {
				var responseBody gin.H
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, gin.H{"message": "Invalid request"}, responseBody)
			} else {
				var responseBody models.Problem
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, responseBody.Question, tt.requestBody["Question"])
				assert.Equal(t, responseBody.Answer, tt.requestBody["Answer"])

				// Assert problem exists
				p, err := testDataStore.GetProblemById(responseBody.Id)
				assert.NoError(t, err)
				assert.Equal(t, tt.requestBody["Question"], p.Question)
				assert.Equal(t, tt.requestBody["Answer"], p.Answer)
			}
		})
	}
}
