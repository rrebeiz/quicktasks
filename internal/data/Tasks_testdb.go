package data

import (
	"context"
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
			Completed:   false,
		}, nil
	}
	return nil, ErrRecordNotFound
}
