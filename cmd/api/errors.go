package main

import "net/http"

func (app *application) logError(r *http.Request, err error) {
	app.errorLog.Printf("method: %s error: %s", r.Method, err.Error())
}

func (app *application) errorResponse(w http.ResponseWriter, status int, r *http.Request, message any) {
	env := envelope{"error": message}
	err := app.writeJSON(w, env, status, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, http.StatusInternalServerError, r, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, http.StatusNotFound, r, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	app.errorResponse(w, http.StatusBadRequest, r, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, http.StatusUnprocessableEntity, r, errors)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, http.StatusConflict, r, message)
}
