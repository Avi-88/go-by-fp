package main

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"
)

type rateLimitEntry struct {
	count int 
	startTime time.Time
}

type ClientStats struct {
	clients map[string]*rateLimitEntry
	mu sync.Mutex
}

var cl ClientStats = ClientStats{
	clients: make(map[string]*rateLimitEntry),
	mu: sync.Mutex{},
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleProtectedRoute)
	// auth := APIKeyAuth([]string{"secret123", "dev-key"})
	rateLimited := RateLimit(5)
	protectedMux := rateLimited(mux)

	err := http.ListenAndServe(":8080", protectedMux)
	if err != nil {
		fmt.Println("There was an error starting up the server")
	}
}

func APIKeyAuth(validKeys []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenLine := r.Header.Get("Authorization")
			if tokenLine == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Invalid request - needs authorization"))
				return
			}
			fields := strings.SplitN(tokenLine, " ", 2)
			if len(fields) < 2 {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Invalid token passed"))
				return
			}
			token := fields[1]
			if !slices.Contains(validKeys, token) {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Invalid token passed"))
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), "api-key", token))
			next.ServeHTTP(w, r)
		})
	}
}

func RateLimit(requestsPerSecond int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			source := r.RemoteAddr
			ip := strings.Split(source, ":")[0]
			cl.mu.Lock()
			defer cl.mu.Unlock()
			entry, ok := cl.clients[ip]
			ct := time.Now()
			if !ok || time.Since(entry.startTime).Seconds() >= 1 {
				cl.clients[ip] = &rateLimitEntry{
					count: 0,	
					startTime: ct,			
				}
			} 
			cl.clients[ip].count++
			if cl.clients[ip].count > requestsPerSecond {
				w.Header().Add("Retry-After",  "1")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Request limit exceeded, try after 1s"))
				return
			}else{
				next.ServeHTTP(w, r)
			}
		})
	}
}



func handleProtectedRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Shhh! here's the secret - %s", r.Context().Value("api-key"))))
}
