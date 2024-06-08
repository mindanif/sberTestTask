package repository

import (
	"context"
	"sberTestTask/internal/todo"

	"time"
)

type TodoRepository interface {
	CreateTask(ctx context.Context, task *todo.Task) error
	GetTask(ctx context.Context, id int) (*todo.Task, error)
	UpdateTask(ctx context.Context, task *todo.Task) error
	DeleteTask(ctx context.Context, id int) error
	ListTasks(ctx context.Context, completed *bool, dueDate *time.Time, limit, offset int) ([]*todo.Task, error)
	CountTasks(ctx context.Context, completed *bool, dueDate *time.Time) (int, error)
}
