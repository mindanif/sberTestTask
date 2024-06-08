package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"sberTestTask/internal/todo"
	"sberTestTask/internal/todo/service"
	"sberTestTask/internal/todo/tests/mocks/serviceMock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupRouterWithMock() (*chi.Mux, *serviceMock.MockTodoUsecase) {
	mockUsecase := new(serviceMock.MockTodoUsecase)
	handler := NewHandler(mockUsecase)
	router := chi.NewRouter()
	router.Get("/tasks/{id}", handler.GetTask)
	router.Post("/tasks", handler.CreateTask)

	return router, mockUsecase
}
func TestCreateTaskServerError(t *testing.T) {
	router, mockUsecase := setupRouterWithMock()

	// Mock responses and test cases
	date, _ := time.Parse(time.RFC3339, "2024-06-07T15:00:00Z")
	tests := []struct {
		name           string
		task           *todo.Task
		body           []byte
		mockReturn     error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Internal Server Error",
			task: &todo.Task{
				Title:       "Test Task",
				Description: "This is a test",
				DueDate:     &date,
			},
			mockReturn:     service.ErrOnServer,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "error on server\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			if tt.body != nil {
				body = tt.body
			} else {
				body, err = json.Marshal(tt.task)
				assert.NoError(t, err)
				mockUsecase.On("CreateTask", mock.Anything, tt.task).Return(tt.mockReturn)
			}

			req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}
func TestCreateTask(t *testing.T) {
	router, mockUsecase := setupRouterWithMock()

	date, _ := time.Parse(time.RFC3339, "2024-06-07T15:00:00Z")
	tests := []struct {
		name           string
		task           *todo.Task
		body           []byte
		mockReturn     error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid Task",
			task: &todo.Task{
				Title:       "Test Task",
				Description: "This is a test",
				DueDate:     &date,
			},
			mockReturn:     nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   fmt.Sprintf(`{"title":"Test Task","description":"This is a test","due_date":"%s`, date.Format(time.RFC3339)) + `","completed":false}` + "\n",
		},
		{
			name:           "Invalid JSON",
			body:           []byte("Invalid JSON"),
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid character 'I' looking for beginning of value\n",
		},
		{
			name: "Validation Error",
			task: &todo.Task{
				Title:       "",
				Description: "No title",
				DueDate:     &date,
			},
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "title cannot be empty\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			if tt.body != nil {
				body = tt.body
			} else {
				body, err = json.Marshal(tt.task)
				assert.NoError(t, err)
				mockUsecase.On("CreateTask", mock.Anything, tt.task).Return(tt.mockReturn)
			}

			req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

func TestGetTask(t *testing.T) {
	router, mockUsecase := setupRouterWithMock()

	date, _ := time.Parse(time.RFC3339, "2024-06-07T15:00:00Z")
	mockTask := &todo.Task{
		ID:          1,
		Title:       "Sample Task",
		Description: "This is a sample task.",
		DueDate:     &date,
	}
	mockUsecase.On("GetTask", mock.Anything, 1).Return(mockTask, nil)
	mockUsecase.On("GetTask", mock.Anything, 2).Return((*todo.Task)(nil), service.ErrIdNotFound)

	tests := []struct {
		name           string
		taskID         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Task",
			taskID:         "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"title":"Sample Task","description":"This is a sample task.","due_date":"2024-06-07T15:00:00Z"` + `,"completed":false}` + "\n",
		},
		{
			name:           "Task Not Found",
			taskID:         "2",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "id not found",
		},
		{
			name:           "Invalid ID Format",
			taskID:         "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/tasks/"+tt.taskID, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}
func setupRouterWithMockForUpdate() (*chi.Mux, *serviceMock.MockTodoUsecase) {
	mockUsecase := new(serviceMock.MockTodoUsecase)
	handler := NewHandler(mockUsecase)
	router := chi.NewRouter()
	router.Put("/tasks/{id}", handler.UpdateTask)
	return router, mockUsecase
}

func TestUpdateTask(t *testing.T) {
	router, mockUsecase := setupRouterWithMockForUpdate()

	date, _ := time.Parse(time.RFC3339, "2024-06-07T15:00:00Z")
	mockTask := &todo.Task{
		ID:          1,
		Title:       "Sample Task",
		Description: "This is a sample task.",
		DueDate:     &date,
		Completed:   false,
	}
	tests := []struct {
		name             string
		taskID           string
		existingTask     *todo.Task
		updates          map[string]interface{}
		mockGetReturn    error
		mockUpdateReturn error
		expectedStatus   int
		expectedBody     string
	}{
		{
			name:             "Valid Update",
			taskID:           "1",
			existingTask:     mockTask,
			updates:          map[string]interface{}{"title": "Updated Task"},
			mockGetReturn:    nil,
			mockUpdateReturn: nil,
			expectedStatus:   http.StatusOK,
			expectedBody:     `{"id":1,"title":"Updated Task","description":"This is a sample task.","due_date":"` + date.Format(time.RFC3339) + `","completed":false}` + "\n",
		},
		{
			name:             "Task Not Found",
			taskID:           "2",
			existingTask:     &todo.Task{},
			updates:          map[string]interface{}{"title": "Updated Task"},
			mockGetReturn:    service.ErrIdNotFound,
			mockUpdateReturn: nil,
			expectedStatus:   http.StatusNotFound,
			expectedBody:     "task not found\n",
		},
		{
			name:             "Invalid JSON",
			taskID:           "1",
			existingTask:     mockTask,
			updates:          nil,
			mockGetReturn:    nil,
			mockUpdateReturn: nil,
			expectedStatus:   http.StatusBadRequest,
			expectedBody:     "invalid character 'I' looking for beginning of value\n",
		},
		{
			name:             "Internal Server Error",
			taskID:           "1",
			existingTask:     mockTask,
			updates:          map[string]interface{}{"title": "Updated Task"},
			mockGetReturn:    nil,
			mockUpdateReturn: errors.New("internal error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedBody:     "internal error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.Calls = nil
			mockUsecase.ExpectedCalls = nil

			if tt.existingTask != nil {
				mockUsecase.On("GetTask", mock.Anything, mock.AnythingOfType("int")).Return(tt.existingTask, tt.mockGetReturn)
			} else {
				mockUsecase.On("GetTask", mock.Anything, mock.AnythingOfType("int")).Return(nil, tt.mockGetReturn)
			}
			if tt.updates != nil {
				updatedTask, _ := applyUpdates(*mockTask, tt.updates)
				mockUsecase.On("UpdateTask", mock.Anything, &updatedTask).Return(tt.mockUpdateReturn)
			}

			var body []byte
			if tt.updates != nil {
				body, _ = json.Marshal(tt.updates)
			} else {
				body = []byte("Invalid JSON")
			}

			req := httptest.NewRequest("PUT", "/tasks/"+tt.taskID, bytes.NewBuffer(body))
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}
func setupRouterWithMockForDelete() (*chi.Mux, *serviceMock.MockTodoUsecase) {
	mockUsecase := new(serviceMock.MockTodoUsecase)
	handler := NewHandler(mockUsecase)
	router := chi.NewRouter()
	router.Delete("/tasks/{id}", handler.DeleteTask)
	return router, mockUsecase
}

func TestDeleteTask(t *testing.T) {
	router, mockUsecase := setupRouterWithMockForDelete()

	tests := []struct {
		name             string
		taskID           string
		mockGetReturn    error
		mockDeleteReturn error
		expectedStatus   int
		expectedBody     string
	}{
		{
			name:             "Successful Delete",
			taskID:           "1",
			mockGetReturn:    nil,
			mockDeleteReturn: nil,
			expectedStatus:   http.StatusOK,
			expectedBody:     "",
		},
		{
			name:             "Invalid ID Format",
			taskID:           "abc",
			mockGetReturn:    nil,
			mockDeleteReturn: nil,
			expectedStatus:   http.StatusBadRequest,
			expectedBody:     "strconv.Atoi: parsing \"abc\": invalid syntax\n",
		},
		{
			name:             "Task Not Found",
			taskID:           "2",
			mockGetReturn:    service.ErrIdNotFound,
			mockDeleteReturn: nil,
			expectedStatus:   http.StatusNotFound,
			expectedBody:     "id not found\n",
		},
		{
			name:             "Internal Server Error",
			taskID:           "3",
			mockGetReturn:    nil,
			mockDeleteReturn: service.ErrOnServer,
			expectedStatus:   http.StatusInternalServerError,
			expectedBody:     "error on server\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.Calls = nil
			mockUsecase.ExpectedCalls = nil

			if tt.mockGetReturn == service.ErrIdNotFound {
				mockUsecase.On("GetTask", mock.Anything, mock.AnythingOfType("int")).Return((*todo.Task)(nil), tt.mockGetReturn)
			} else {
				mockUsecase.On("GetTask", mock.Anything, mock.AnythingOfType("int")).Return(&todo.Task{}, tt.mockGetReturn)
			}
			mockUsecase.On("DeleteTask", mock.Anything, mock.AnythingOfType("int")).Return(tt.mockDeleteReturn)

			req := httptest.NewRequest("DELETE", "/tasks/"+tt.taskID, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
func setupRouterWithMockForList() (*chi.Mux, *serviceMock.MockTodoUsecase) {
	mockUsecase := new(serviceMock.MockTodoUsecase)
	handler := NewHandler(mockUsecase)
	router := chi.NewRouter()
	router.Get("/tasks", handler.ListTasks)
	return router, mockUsecase
}

func TestListTasks(t *testing.T) {
	router, mockUsecase := setupRouterWithMockForList()

	date, _ := time.Parse(time.RFC3339, "2024-06-07T15:00:00Z")
	mockPages := &todo.Pages{
		CountPage: 1,
		CurPage:   1,
		Tasks: []*todo.Task{
			{
				ID:          1,
				Title:       "Sample Task 1",
				Description: "This is a sample task 1.",
				DueDate:     &date,
				Completed:   false,
			},
		},
	}

	tests := []struct {
		name            string
		queryParams     string
		mockReturn      *todo.Pages
		mockReturnError error
		expectedStatus  int
		expectedBody    string
		isJson          bool
	}{
		{
			name:            "Successful Request With Params",
			queryParams:     "?completed=true&date=" + date.Format(time.DateOnly) + "&limit=5&page=1",
			mockReturn:      mockPages,
			mockReturnError: nil,
			expectedStatus:  http.StatusOK,
			expectedBody:    `{"count_page":1,"cur_page":1,"tasks":[{"id":1,"title":"Sample Task 1","description":"This is a sample task 1.","due_date":"` + date.Format(time.RFC3339) + `","completed":false}]}`,
			isJson:          true,
		},
		{
			name:            "Successful Request Without Params",
			queryParams:     "",
			mockReturn:      mockPages,
			mockReturnError: nil,
			expectedStatus:  http.StatusOK,
			expectedBody:    `{"count_page":1,"cur_page":1,"tasks":[{"id":1,"title":"Sample Task 1","description":"This is a sample task 1.","due_date":"` + date.Format(time.RFC3339) + `","completed":false}]}`,
			isJson:          true,
		},
		{
			name:            "Invalid Completed Flag",
			queryParams:     "?completed=invalid",
			mockReturn:      nil,
			mockReturnError: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    "invalid completed flag\n",
			isJson:          false,
		},
		{
			name:            "Invalid Date Format",
			queryParams:     "?date=invalid-date",
			mockReturn:      nil,
			mockReturnError: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    "invalid date format\n",
			isJson:          false,
		},
		{
			name:            "Internal Server Error",
			queryParams:     "",
			mockReturn:      nil,
			mockReturnError: service.ErrOnServer,
			expectedStatus:  http.StatusInternalServerError,
			expectedBody:    "error retrieving tasks\n",
			isJson:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.Calls = nil
			mockUsecase.ExpectedCalls = nil

			mockUsecase.On("ListTasks", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.mockReturn, tt.mockReturnError)

			req := httptest.NewRequest("GET", "/tasks"+tt.queryParams, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != "" {
				if tt.isJson {
					assert.JSONEq(t, tt.expectedBody, rr.Body.String())
				} else {
					assert.Contains(t, rr.Body.String(), tt.expectedBody)
				}
			}
		})
	}
}
