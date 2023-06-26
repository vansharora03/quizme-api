package main

import (
    "fmt"
    "time"
    "net/http"
)


// serve initializes and starts the server with app
func (app *application) serve() error {

	// Prepare server
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start server
	app.logger.Printf("Starting %s server on port %s", app.config.env, srv.Addr)

	return srv.ListenAndServe()

}
