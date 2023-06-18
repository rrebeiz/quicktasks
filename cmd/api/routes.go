package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://localhost", "http://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Get("/healthcheck", app.healthcheckHandler)

	router.Get("/v1/tasks", app.getAllTasks)
	router.Get("/v1/tasks/{id}", app.getTaskHandler)
	router.Post("/v1/tasks", app.createTaskHandler)
	router.Patch("/v1/tasks/{id}", app.updateTaskHandler)
	router.Delete("/v1/tasks/{id}", app.deleteTaskHandler)

	return router
}
