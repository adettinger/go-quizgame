package webserver_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adettinger/go-quizgame/models"
	webserver "github.com/adettinger/go-quizgame/web-server"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetProblemById(t *testing.T) {
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
			name:     "Success",
			urlParam: "123e4567-e89b-12d3-a456-426614174000",
			problem: models.Problem{
				Id:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Question: "What is 2+2?",
				Answer:   "4",
			},
			expectedStatus: http.StatusOK,
			expectedBody: models.Problem{ //TODO: Redundant
				Id:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Question: "What is 2+2?",
				Answer:   "4",
			},
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
			testDS, _ := webserver.NewDataStoreFromData([]models.Problem{tt.problem})
			controller := webserver.NewProblemController(testDS)

			// Create a test HTTP response recorder and context
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			// Set up the route
			r.GET("/problems/:id", controller.GetProblemById)

			// Create the request
			req, _ := http.NewRequest("GET", "/problems/"+tt.urlParam, nil)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: tt.urlParam}}

			// Call the handler
			controller.GetProblemById(c)

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
