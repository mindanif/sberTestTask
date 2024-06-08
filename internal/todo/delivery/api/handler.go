package api

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sberTestTask/internal/todo"
	"sberTestTask/internal/todo/service"
	"strconv"
	"time"
)

const (
	defaultPage  = 1
	defaultLimit = 10
)

type Handler struct {
	uc service.TodoUsecase
}

func NewHandler(uc service.TodoUsecase) *Handler {
	return &Handler{uc: uc}
}

// @Summary Create a new task
// @Description Create a new task with the input payload
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param task body todo.Task true "Task to create"
// @Success 201 {object} todo.Task "Task created successfully"
// @Failure 400 {object} todo.ErrorResponse "Bad Request"
// @Failure 500 {object} todo.ErrorResponse "Internal Server Error"
// @Router /tasks [post]
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task todo.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validateTask(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.uc.CreateTask(r.Context(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// @Summary Get a task by ID
// @Description Get a task by ID
// @Tags tasks
// @Produce  json
// @Param id path int true "Task ID"
// @Success 200 {object} todo.Task "Task found"
// @Failure 400 {object} todo.ErrorResponse "Bad Request"
// @Failure 404 {object} todo.ErrorResponse "Not Found"
// @Router /tasks/{id} [get]
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task, err := h.uc.GetTask(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(task)
}

// @Summary Update a task
// @Summary Update a task
// @Description Update a task with the input payload
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path int true "Task ID"
// @Param task body map[string]interface{} true "Task updates"
// @Success 200 {object} todo.Task "Task updated successfully"
// @Failure 400 {object} todo.ErrorResponse "Bad Request"
// @Failure 404 {object} todo.ErrorResponse "Not Found"
// @Failure 500 {object} todo.ErrorResponse "Internal Server Error"
// @Router /tasks/{id} [put]
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingTask, err := h.uc.GetTask(r.Context(), id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTask, err := applyUpdates(*existingTask, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.uc.UpdateTask(r.Context(), &updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}

// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} todo.ErrorResponse "Bad Request"
// @Failure 404 {object} todo.ErrorResponse "Not Found"
// @Failure 500 {object} todo.ErrorResponse "Internal Server Error"
// @Router /tasks/{id} [delete]
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := h.uc.GetTask(r.Context(), id); err != nil {
		http.Error(w, service.ErrIdNotFound.Error(), http.StatusNotFound)
		return
	}
	err = h.uc.DeleteTask(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary List tasks
// @Description Get a list of tasks with optional filters
// @Tags tasks
// @Produce  json
// @Param completed query bool false "Filter by completion status"
// @Param date query string false "Filter by due date" Format(date) example(2024-06-07) name(2024-06-07)
// @Param limit query int false "Number of tasks per page"
// @Param page query int false "Page number"
// @Success 200 {object} todo.Pages "List of tasks"
// @Failure 400 {object} todo.ErrorResponse "Bad Request"
// @Failure 500 {object} todo.ErrorResponse "Internal Server Error"
// @Router /tasks [get]
func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {

	var completed *bool
	if completedStr := r.URL.Query().Get("completed"); completedStr != "" {
		completedVal, err := strconv.ParseBool(completedStr)
		if err != nil {
			http.Error(w, "invalid completed flag", http.StatusBadRequest)
			return
		}
		completed = &completedVal
	}

	var dueDate *time.Time
	if dateStr := r.URL.Query().Get("date"); dateStr != "" {
		parsedDate, err := time.Parse(time.DateOnly, dateStr)
		if err != nil {
			http.Error(w, "invalid date format", http.StatusBadRequest)
			return
		}
		dueDate = &parsedDate
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = defaultPage
	}

	pages, err := h.uc.ListTasks(r.Context(), completed, dueDate, limit, page)
	if err != nil {
		http.Error(w, "error retrieving tasks", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pages)
}
func applyUpdates(task todo.Task, updates map[string]interface{}) (todo.Task, error) {
	for key, value := range updates {
		switch key {
		case "title":
			if v, ok := value.(string); ok {
				task.Title = v
			}
		case "description":
			if v, ok := value.(string); ok {
				task.Description = v
			}
		case "due_date":
			if v, ok := value.(string); ok {
				dueDate, err := time.Parse(time.RFC3339, v)
				if err != nil {
					return task, err
				}
				task.DueDate = &dueDate
			}
		case "completed":
			if v, ok := value.(bool); ok {
				task.Completed = v
			}
		}
	}
	return task, nil
}
func validateTask(task *todo.Task) error {
	if task.Title == "" {
		return errors.New("title cannot be empty")
	}
	if task.DueDate == nil {
		return errors.New("missed data field")
	}
	return nil
}
