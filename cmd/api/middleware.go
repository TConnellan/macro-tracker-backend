package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {

	if !app.config.limiter.enabled {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	type client struct {
		limiter      *rate.Limiter
		lastSeenTime time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeenTime) > 3*time.Minute {
					delete(clients, ip)

				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				time.Now(),
			}
		}

		if !clients[ip].limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			mu.Unlock()
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
