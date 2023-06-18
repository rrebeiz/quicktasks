package data

import (
	"context"
	"errors"
)

type TaskMockModel struct {
}

func NewTaskMockModel() TaskMockModel {
	return TaskMockModel{}
}

func (t TaskMockModel) GetTaskByID(ctx context.Context, id int64) (*Task, error) {
	if id == 1 {
		return &Task{
			ID:          1,
			Title:       "Test Title",
			Description: "Test Description",
			Complete:    false,
		}, nil
	}
	return nil, ErrRecordNotFound
}

func (t TaskMockModel) CreateTask(ctx context.Context, task *Task) error {
	if task.Title == "Test Task" {
		task.ID = 1
		task.Version = 1
		return nil
	}
	return errors.New("something went wrong")
}

func (t TaskMockModel) UpdateTask(ctx context.Context, task *Task) error {
	if task.ID == 1 {
		return nil
	}
	return ErrEditConflict
}

func (t TaskMockModel) DeleteTask(ctx context.Context, task *Task) error {
	if task.ID == 1 {
		return nil
	}
	return ErrRecordNotFound
}

func (t TaskMockModel) GetAllTasks(ctx context.Context) ([]*Task, error) {
	var tasks []*Task
	tasks = append(tasks, &Task{
		ID:          1,
		Title:       "Test Task",
		Description: "Test Description",
		Complete:    false,
		Version:     1,
	})
	return tasks, nil
}
