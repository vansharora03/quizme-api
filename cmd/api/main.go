package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"time"
	"vanshadhruvp/quizme-api/internal/data"

	_ "github.com/lib/pq"
)

type config struct {
	env  string
	port int
    db struct{
        dsn string
        MaxOpenConns int
        MaxIdleConns int
        MaxIdleTime string 
    }
    limiter struct {
        rps float64
        burst int
        enabled bool
    }
}

type application struct {
	config config
	logger *log.Logger
    models data.Models
}

const version = "1.0.0"

func main() {
	var cfg config

	// Parse command line into cfg
	flag.StringVar(
        &cfg.env, 
        "env", 
        "development", 
        "env is the current environment (development|production|staging)")
	flag.IntVar(&cfg.port, 
        "port", 
        8080, 
        "port is the address the server is listening on")
    flag.StringVar(
        &cfg.db.dsn, 
        "dsn", 
        "", 
        "data source name to connect to postgres")
    flag.IntVar(
        &cfg.db.MaxOpenConns, 
        "DB-MaxOpenConns", 
         25, 
        "maximum number of open connections allowed on db")
    flag.IntVar(
        &cfg.db.MaxIdleConns, 
        "DB-MaxIdleConns", 
        25, 
        "maximum number of idle connections allowed on db")
	flag.StringVar(
        &cfg.db.MaxIdleTime, 
        "DB-MaxIdleTime", 
        "15m", 
        "maximum time for idle connections to stay alive")
	flag.IntVar(
        &cfg.limiter.burst,
        "Limiter-burst", 
        4, 
        "Rate limiter maximum burst")
	flag.Float64Var(
        &cfg.limiter.rps, 
        "Limiter-rps", 
        2, 
        "Rate limiter maximum requests per second")
	flag.BoolVar(
        &cfg.limiter.enabled,
        "Limiter-enabled", 
        true, 
        "Enable rate limiter")
	flag.Parse()

	// Prepare dependencies for app
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

    // Establish connection to database
    db, err := openDB(cfg)
    if err != nil {
        logger.Fatal(err)
    }
    defer db.Close()

    logger.Println("Connection to database established.")

	// Prepare app
	app := &application{
		config: cfg,
		logger: logger,
        models: data.NewModels(db),
	}

    // Start server
    err = app.serve()
	if err != nil {
		logger.Fatal(err)
	}

}

// openDB attempts to open a connection to a database using the
// cfg dsn. It then pings this connection to validate the
// connection. If successful, the database connection pool will
// be returned.
func openDB(cfg config) (*sql.DB, error) {
    db, err := sql.Open("postgres", cfg.db.dsn)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(cfg.db.MaxOpenConns)
    db.SetMaxIdleConns(cfg.db.MaxIdleConns)
    duration, err := time.ParseDuration(cfg.db.MaxIdleTime)
    if err != nil {
        return nil, err
    }
    db.SetConnMaxIdleTime(duration)

    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    err = db.PingContext(ctx)
    if err != nil {
        return nil, err
    }

    return db, nil

}
