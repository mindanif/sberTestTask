package repositoryMock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"sberTestTask/internal/todo"
	"time"
)

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) CreateTask(ctx context.Context, task *todo.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTodoRepository) GetTask(ctx context.Context, id int) (*todo.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*todo.Task), args.Error(1)
}

func (m *MockTodoRepository) UpdateTask(ctx context.Context, task *todo.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTodoRepository) DeleteTask(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTodoRepository) ListTasks(ctx context.Context, completed *bool, dueDate *time.Time, limit, offset int) ([]*todo.Task, error) {
	args := m.Called(ctx, completed, dueDate, limit, offset)
	return args.Get(0).([]*todo.Task), args.Error(1)
}

func (m *MockTodoRepository) CountTasks(ctx context.Context, completed *bool, dueDate *time.Time) (int, error) {
	args := m.Called(ctx, completed, dueDate)
	return args.Int(0), args.Error(1)
}
