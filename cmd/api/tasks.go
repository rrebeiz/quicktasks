package main

import (
	"errors"
	"github.com/rrebeiz/quicktasks/internal/data"
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
	// get the id from the DB
	// for now just print it out
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
