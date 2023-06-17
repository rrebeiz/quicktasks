package main

import (
	"errors"
	"fmt"
	"github.com/rrebeiz/quicktasks/internal/data"
	"github.com/rrebeiz/quicktasks/internal/validator"
	"net/http"
)

func (app *application) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.getParamID(r)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidIDParam):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	task, err := app.models.Tasks.GetTaskByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, envelope{"task": task}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.NewValidator()

	task := &data.Task{
		Title:       input.Title,
		Description: input.Description,
	}

	data.ValidateTask(v, task)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tasks.CreateTask(r.Context(), task)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)

	location := fmt.Sprintf("/v1/tasks/%d", task.ID)

	headers.Set("Location", location)

	err = app.writeJSON(w, envelope{"task": task}, http.StatusCreated, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.getParamID(r)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidIDParam):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Completed   *bool   `json:"completed"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.NewValidator()

	task, err := app.models.Tasks.GetTaskByID(r.Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if input.Title != nil {
		v.Check(*input.Title != "", "title", "should not be empty")
		task.Title = *input.Title
	}
	if input.Description != nil {
		v.Check(*input.Description != "", "description", "should not be empty")
		task.Description = *input.Description
	}

	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tasks.UpdateTask(r.Context(), task)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, envelope{"task": task}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
