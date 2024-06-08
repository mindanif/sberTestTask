package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sberTestTask/internal/todo"
	"sberTestTask/internal/todo/tests/mocks/repositoryMock"
)

func TestCreateTask(t *testing.T) {
	mockRepo := new(repositoryMock.MockTodoRepository)
	svc := NewTodoUsecase(mockRepo)

	task := &todo.Task{
		Title:       "Test Task",
		Description: "This is a test task",
	}

	t.Run("Successful CreateTask", func(t *testing.T) {
		mockRepo.On("CreateTask", mock.Anything, task).Return(nil)

		err := svc.CreateTask(context.Background(), task)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CreateTask Error", func(t *testing.T) {
		mockRepo.Calls = nil
		mockRepo.ExpectedCalls = nil
		mockRepo.On("CreateTask", mock.Anything, task).Return(ErrOnServer)

		err := svc.CreateTask(context.Background(), task)
		assert.Error(t, err)
		assert.Equal(t, ErrOnServer, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetTask(t *testing.T) {
	mockRepo := new(repositoryMock.MockTodoRepository)
	svc := NewTodoUsecase(mockRepo)

	task := &todo.Task{
		ID:          1,
		Title:       "Test Task",
		Description: "This is a test task",
	}

	mockRepo.On("GetTask", mock.Anything, 1).Return(task, nil)

	result, err := svc.GetTask(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, task, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTask(t *testing.T) {
	mockRepo := new(repositoryMock.MockTodoRepository)
	svc := NewTodoUsecase(mockRepo)

	task := &todo.Task{
		ID:          1,
		Title:       "Updated Task",
		Description: "This is an updated test task",
	}

	mockRepo.On("UpdateTask", mock.Anything, task).Return(nil)

	err := svc.UpdateTask(context.Background(), task)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTask(t *testing.T) {
	mockRepo := new(repositoryMock.MockTodoRepository)
	svc := NewTodoUsecase(mockRepo)

	mockRepo.On("DeleteTask", mock.Anything, 1).Return(nil)

	err := svc.DeleteTask(context.Background(), 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestListTasks(t *testing.T) {
	mockRepo := new(repositoryMock.MockTodoRepository)
	svc := NewTodoUsecase(mockRepo)

	date := time.Now()
	completed := true
	tasks := []*todo.Task{
		{ID: 1, Title: "Test Task 1", Description: "This is a test task 1", DueDate: &date, Completed: completed},
		{ID: 2, Title: "Test Task 2", Description: "This is a test task 2", DueDate: &date, Completed: completed},
	}

	t.Run("Successful ListTasks", func(t *testing.T) {
		mockRepo.Calls = nil
		mockRepo.ExpectedCalls = nil
		mockRepo.On("CountTasks", mock.Anything, &completed, &date).Return(2, nil)
		mockRepo.On("ListTasks", mock.Anything, &completed, &date, 10, 0).Return(tasks, nil)

		result, err := svc.ListTasks(context.Background(), &completed, &date, 10, 1)
		expectedPages := &todo.Pages{
			CountPage: 1,
			CurPage:   1,
			Tasks:     tasks,
		}

		assert.NoError(t, err)
		assert.Equal(t, expectedPages, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Counting Tasks", func(t *testing.T) {
		mockRepo.Calls = nil
		mockRepo.ExpectedCalls = nil
		mockRepo.On("CountTasks", mock.Anything, &completed, &date).Return(0, errors.New("count error"))

		result, err := svc.ListTasks(context.Background(), &completed, &date, 10, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrOnServer, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error Getting ListTasks", func(t *testing.T) {
		mockRepo.Calls = nil
		mockRepo.ExpectedCalls = nil
		mockRepo.On("CountTasks", mock.Anything, &completed, &date).Return(2, nil)
		mockRepo.On("ListTasks", mock.Anything, &completed, &date, 10, 0).Return(([]*todo.Task)(nil), errors.New("list error"))

		result, err := svc.ListTasks(context.Background(), &completed, &date, 10, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrOnServer, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCountTasks(t *testing.T) {
	mockRepo := new(repositoryMock.MockTodoRepository)
	svc := NewTodoUsecase(mockRepo)

	date := time.Now()
	completed := true

	mockRepo.On("CountTasks", mock.Anything, &completed, &date).Return(1, nil)

	count, err := svc.CountTasks(context.Background(), &completed, &date)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	mockRepo.AssertExpectations(t)
}
