package main

import (
	"testing"
	"time"

	"github.com/hellofresh/health-go/v5"
	healthHttp "github.com/hellofresh/health-go/v5/checks/http"
)

func TestServerHealth(t *testing.T) {
	h, _ := health.New(health.WithSystemInfo())
	h.Register(health.Config{
		Name:      "http-check",
		Timeout:   time.Second * 5,
		SkipOnErr: true,
		Check: healthHttp.New(healthHttp.Config{
			URL: "http://localhost:8080",
		}),
	})

}
