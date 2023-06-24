package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/justinas/nosurf"
	"golang.org/x/time/rate"
)

// Secure headers for http requests
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set security headers
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Log request based on method being requested at a URL
func (app *application) logRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// recoverPanic recovers the handler from any panics and
// sends a http.StatusInternalServerError response
func (app *application) recoverPanic(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		handler.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

// rateLimit a middleware that limits specific clients with a number of requests per
// second, with a rps and burst value determined by app.config.limiter
func (app *application) rateLimit(next http.Handler) http.Handler {

    // If rate limiting is disabled, skip this handler
    if !app.config.limiter.enabled {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            next.ServeHTTP(w, r)
        })
    }
    // Information saved for each client
    type client struct {
        limiter *rate.Limiter // Determines how many requests a client will be able to make
        lastSeen time.Time // Determines when the client will be removed from this map
    }

    var (
        mu sync.Mutex
        clients = make(map[string]*client)
    )


    // Run background process for removing expired clients
    go func() {
        for {
            time.Sleep(time.Minute)

            mu.Lock()

            for ip, client := range clients {
                if time.Since(client.lastSeen) > 3 * time.Minute {
                    delete(clients, ip)
                }
            }
            
            mu.Unlock()
        }
    }()

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get user ip address
        ip, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            app.serverErrorResponse(w, r, err)
            return
        }

        mu.Lock()

        // Assign new clients a limiter
        if _, found := clients[ip]; !found {
            clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), 
                app.config.limiter.burst)}
        }

        // Update lastSeen since the client is being seen
        clients[ip].lastSeen = time.Now()

        // Determine if the client has exceeded its limit, do not serve http if so
        if !clients[ip].limiter.Allow() {
            mu.Unlock()
            app.errorResponse(w, r, http.StatusTooManyRequests, "Rate limit exceeded")
            return
        }

        mu.Unlock()

        // Serve client if they have not exceeded their rate limit
        next.ServeHTTP(w, r)
    })

}
