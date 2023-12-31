package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/rrebeiz/quicktasks/internal/validator"
	"time"
)

type Tasks interface {
	GetAllTasks(ctx context.Context) ([]*Task, error)
	GetTaskByID(ctx context.Context, id int64) (*Task, error)
	CreateTask(ctx context.Context, task *Task) error
	UpdateTask(ctx context.Context, task *Task) error
	DeleteTask(ctx context.Context, task *Task) error
}
type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Complete    bool      `json:"complete"`
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

func (t TaskModel) GetAllTasks(ctx context.Context) ([]*Task, error) {
	query := `select id, title, description, complete from tasks`
	var tasks []*Task
	rows, err := t.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Complete)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

func (t TaskModel) GetTaskByID(ctx context.Context, id int64) (*Task, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `select id, title, description, complete, version from tasks where id = $1`
	var task Task
	err := t.DB.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Complete, &task.Version)

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
	query := `insert into tasks (title, description, complete) VALUES ($1, $2, $3) returning id, version`
	args := []any{task.Title, task.Description, task.Complete}

	err := t.DB.QueryRowContext(ctx, query, args...).Scan(&task.ID, &task.Version)

	if err != nil {
		return err
	}
	return nil
}

func (t TaskModel) UpdateTask(ctx context.Context, task *Task) error {
	query := `update tasks set title = $1, description = $2, complete = $3, version = version + 1 where id = $4 and version = $5 returning version`
	args := []any{task.Title, task.Description, task.Complete, task.ID, task.Version}
	err := t.DB.QueryRowContext(ctx, query, args...).Scan(&task.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (t TaskModel) DeleteTask(ctx context.Context, task *Task) error {

	query := `delete from tasks where id = $1 and version = $2`
	args := []any{task.ID, task.Version}
	result, err := t.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}

func ValidateTask(v *validator.Validator, task *Task) {
	v.Check(task.Title != "", "title", "should not be empty")
	v.Check(task.Description != "", "description", "should not be empty")
}
