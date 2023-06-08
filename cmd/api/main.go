package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	env  string
	port int
}

type application struct {
	config config
	logger *log.Logger
}

const version = "1.0.0"

func main() {
	var cfg config

	// Parse command line into cfg
	flag.StringVar(&cfg.env, "env", "development", "env is the current environment (development|production|staging)")
	flag.IntVar(&cfg.port, "port", 8080, "port is the address the server is listening on")
	flag.Parse()

	// Prepare dependencies for app
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Prepare app
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Prepare server
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start server
	logger.Printf("Starting %s server on port %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}

}
