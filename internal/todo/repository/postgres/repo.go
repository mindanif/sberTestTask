package postgres

import (
	"context"
	"database/sql"
	"sberTestTask/internal/todo"
	"sberTestTask/internal/todo/repository"
	"strconv"
	"time"
)

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) repository.TodoRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateTask(ctx context.Context, task *todo.Task) error {
	query := `INSERT INTO tasks (title, description, due_date, completed) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, task.Title, task.Description, task.DueDate, task.Completed).Scan(&task.ID)
	return err
}

func (r *postgresRepository) GetTask(ctx context.Context, id int) (*todo.Task, error) {
	task := &todo.Task{}
	err := r.db.QueryRowContext(ctx, "SELECT id, title, description, due_date, completed FROM tasks WHERE id = $1", id).Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Completed)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *postgresRepository) UpdateTask(ctx context.Context, task *todo.Task) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tasks SET title = $1, description = $2, due_date = $3, completed = $4 WHERE id = $5", task.Title, task.Description, task.DueDate, task.Completed, task.ID)
	return err
}

func (r *postgresRepository) DeleteTask(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", id)
	return err
}

func (r *postgresRepository) ListTasks(ctx context.Context, completed *bool, dueDate *time.Time, limit, offset int) ([]*todo.Task, error) {
	var tasks []*todo.Task
	var rows *sql.Rows
	var err error

	query := "SELECT id, title, description, due_date, completed FROM tasks WHERE 1=1"
	args := []interface{}{}

	if completed != nil {
		query += " AND completed = $1"
		args = append(args, *completed)
	}

	if dueDate != nil {
		query += " AND DATE(due_date) = DATE($" + strconv.Itoa(len(args)+1) + ")"
		args = append(args, *dueDate)
	}

	args = append(args, limit)
	args = append(args, offset)

	query += " ORDER BY due_date LIMIT $" + strconv.Itoa(len(args)-1) + " OFFSET $" + strconv.Itoa(len(args))
	rows, err = r.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := new(todo.Task)
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Completed); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
func (r *postgresRepository) CountTasks(ctx context.Context, completed *bool, dueDate *time.Time) (int, error) {
	var err error

	query := "SELECT COUNT(id) FROM tasks WHERE 1=1"
	args := []interface{}{}

	if completed != nil {
		query += " AND completed = $1"
		args = append(args, *completed)
	}

	if dueDate != nil {
		query += " AND DATE(due_date) = DATE($" + strconv.Itoa(len(args)+1) + ")"
		args = append(args, *dueDate)
	}

	var count int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
