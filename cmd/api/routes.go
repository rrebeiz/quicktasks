package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	router.Get("/healthcheck", app.healthcheckHandler)

	router.Get("/v1/tasks/{id}", app.getTaskHandler)
	router.Post("/v1/tasks", app.createTaskHandler)

	return router
}
