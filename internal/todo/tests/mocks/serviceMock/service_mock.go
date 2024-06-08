package serviceMock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"sberTestTask/internal/todo"
	"time"
)

// MockTodoUsecase is a mock type for the TodoUsecase interface
type MockTodoUsecase struct {
	mock.Mock
	Tasks       map[int]*todo.Task
	NextID      int
	CreateErr   error
	UpdateErr   error
	DeleteErr   error
	GetErr      error
	ListErr     error
	CountErr    error
	CountResult int
}

func NewMockTodoUsecase() *MockTodoUsecase {
	return &MockTodoUsecase{
		Tasks:  make(map[int]*todo.Task),
		NextID: 1,
	}
}
func (m *MockTodoUsecase) CreateTask(ctx context.Context, task *todo.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTodoUsecase) GetTask(ctx context.Context, id int) (*todo.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*todo.Task), args.Error(1)
}

func (m *MockTodoUsecase) UpdateTask(ctx context.Context, task *todo.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTodoUsecase) DeleteTask(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTodoUsecase) ListTasks(ctx context.Context, completed *bool, dueDate *time.Time, limit, page int) (*todo.Pages, error) {
	args := m.Called(ctx, completed, dueDate, limit, page)
	return args.Get(0).(*todo.Pages), args.Error(1)
}

func (m *MockTodoUsecase) CountTasks(ctx context.Context, completed *bool, dueDate *time.Time) (int, error) {
	args := m.Called(ctx, completed, dueDate)
	return args.Int(0), args.Error(1)
}
