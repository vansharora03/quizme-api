package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/hellofresh/health-go/v5"
	healthHttp "github.com/hellofresh/health-go/v5/checks/http"
)

func TestServerHealth(t *testing.T) {
	// Create a new health instance
	h, _ := health.New(health.WithSystemInfo())

	// Register a new health check
	h.Register(health.Config{
		Name:      "http-check",
		Timeout:   time.Second * 5,
		SkipOnErr: true,
		Check: healthHttp.New(healthHttp.Config{
			URL: "http://localhost:8080",
		}),
	})

	// Run the health checks
	http.Handle("/status", h.Handler())
	http.ListenAndServe(":8080", nil)
}
