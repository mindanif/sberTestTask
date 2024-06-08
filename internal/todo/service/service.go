package service

import (
	"context"
	"errors"
	"log/slog"
	"sberTestTask/internal/todo"
	"sberTestTask/internal/todo/repository"
	"time"
)

var (
	ErrIdNotFound  = errors.New("id not found")
	ErrInvalidData = errors.New("invalid data")
	ErrOnServer    = errors.New("error on server")
)

type TodoUsecase interface {
	CreateTask(ctx context.Context, task *todo.Task) error
	GetTask(ctx context.Context, id int) (*todo.Task, error)
	UpdateTask(ctx context.Context, task *todo.Task) error
	DeleteTask(ctx context.Context, id int) error
	ListTasks(ctx context.Context, completed *bool, dueDate *time.Time, limit, page int) (*todo.Pages, error)
	CountTasks(ctx context.Context, completed *bool, dueDate *time.Time) (int, error)
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoUsecase(repo repository.TodoRepository) TodoUsecase {
	return &todoService{repo: repo}
}

func (u *todoService) CreateTask(ctx context.Context, task *todo.Task) error {

	if err := u.repo.CreateTask(ctx, task); err != nil {
		slog.Error("create error: ", slog.String("error", err.Error()))
		return ErrOnServer
	}
	return nil
}

func (u *todoService) GetTask(ctx context.Context, id int) (*todo.Task, error) {
	task, err := u.repo.GetTask(ctx, id)
	if err != nil {
		slog.Error("Task not found", slog.String("error", err.Error()))
		return nil, ErrIdNotFound
	}
	return task, nil
}

func (u *todoService) UpdateTask(ctx context.Context, task *todo.Task) error {
	if err := u.repo.UpdateTask(ctx, task); err != nil {
		slog.Error("update error: ", slog.String("error", err.Error()))
		return ErrOnServer
	}
	return nil
}

func (u *todoService) DeleteTask(ctx context.Context, id int) error {

	if err := u.repo.DeleteTask(ctx, id); err != nil {
		slog.Error("delete error: ", slog.String("error", err.Error()))
		return ErrOnServer
	}
	return nil
}

func (u *todoService) ListTasks(ctx context.Context, completed *bool, dueDate *time.Time, limit, page int) (*todo.Pages, error) {

	totalCount, err := u.CountTasks(ctx, completed, dueDate)
	if err != nil {
		slog.Error("error counting tasks", slog.String("error", err.Error()))
		return nil, ErrOnServer
	}

	countPage := totalCount / limit
	if totalCount%limit != 0 {
		countPage++
	}
	if page > countPage {
		page = countPage
	}

	offset := 0
	if page > 1 {
		offset = (page - 1) * limit
	}

	tasks, err := u.repo.ListTasks(ctx, completed, dueDate, limit, offset)
	if err != nil {
		slog.Error("error getting list ", slog.String("error", err.Error()))
		return nil, ErrOnServer
	}
	return &todo.Pages{
		CountPage: countPage,
		CurPage:   page,
		Tasks:     tasks,
	}, nil

}

func (u *todoService) CountTasks(ctx context.Context, completed *bool, dueDate *time.Time) (int, error) {
	return u.repo.CountTasks(ctx, completed, dueDate)
}
