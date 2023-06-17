package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("the requested resource could not be found")
)

type Models struct {
	Tasks Tasks
}

func NewModels(db *sql.DB) Models {
	return Models{
		Tasks: NewTaskModel(db),
	}
}
