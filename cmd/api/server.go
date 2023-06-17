package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
	app.infoLog.Printf("started server on port %d\n", app.config.port)
	return srv.ListenAndServe()
}
