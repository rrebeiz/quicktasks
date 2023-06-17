package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/rrebeiz/quicktasks/internal/validator"
	"time"
)

type Tasks interface {
	GetTaskByID(ctx context.Context, id int64) (*Task, error)
	CreateTask(ctx context.Context, task *Task) error
}
type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Version     int       `json:"-"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type TaskModel struct {
	DB *sql.DB
}

func NewTaskModel(db *sql.DB) TaskModel {
	return TaskModel{DB: db}
}

func (t TaskModel) GetTaskByID(ctx context.Context, id int64) (*Task, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `select id, title, description, completed from tasks where id = $1`
	var task Task
	err := t.DB.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Completed)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &task, nil
}

func (t TaskModel) CreateTask(ctx context.Context, task *Task) error {
	query := `insert into tasks (title, description, completed) VALUES ($1, $2, $3) returning id, version`
	args := []any{task.Title, task.Description, task.Completed}

	err := t.DB.QueryRowContext(ctx, query, args...).Scan(&task.ID, &task.Version)

	if err != nil {
		return err
	}
	return nil
}

func ValidateTask(v *validator.Validator, task *Task) {
	v.Check(task.Title != "", "title", "should not be empty")
	v.Check(task.Description != "", "description", "should not be empty")
}
