package main

import "net/http"

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"system_info": struct {
			Port        int
			Environment string
			Status      string
		}{
			app.config.port,
			app.config.environment,
			"Healthy",
		},
	}
	err := app.writeJSON(w, env, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
