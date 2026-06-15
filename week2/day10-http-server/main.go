package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Stats struct {
	startTime time.Time
	requestCount int
	mu sync.RWMutex
}

var serverStats Stats = Stats{
	startTime: time.Time{},
	requestCount: 0,
	mu: sync.RWMutex{},
}

type ServerStats struct {
	Uptime string `json:"uptime"`
	RequestCount int `json:"request_count"`
}

type StatusRecorder struct {
    http.ResponseWriter
    status int
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRootPath)
	mux.HandleFunc("/hello", handleHelloPath)
	mux.HandleFunc("/echo", handleEchoPath)
	mux.HandleFunc("/status", handleStatusPath)

	mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("Panicking for no apparent reason!!!")
	})

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))



	wMux := RequestCount(Logger(Recovery(RequestID(mux))))
	serverStats.startTime = time.Now()
	err := http.ListenAndServe(":8080", wMux)
	if err != nil {
		fmt.Printf("There was an error starting the server - %v", err)
	}
}

func RequestCount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverStats.mu.Lock()
		serverStats.requestCount++
		serverStats.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

// Logger logs: METHOD PATH → STATUS (duration)
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sr := &StatusRecorder{
			status: 200,
			ResponseWriter: w,
		}
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(sr, r)
		fmt.Printf("%s %s -> %d (%v)\n", r.Method, r.URL.Path, sr.status, time.Since(start).String())
	})
}

// Recovery catches panics and returns 500 instead of crashing
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Caught panic")
				http.Error(w, "Internal Server Error", 500)
			}
		}()
		next.ServeHTTP(w,r)
	})
}

// RequestID adds a unique request ID to every request's context
// and includes it in the response as X-Request-ID header
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := fmt.Sprintf("%d-%s", time.Now().UnixNano(), time.Now().Weekday())
		r = r.WithContext(context.WithValue(r.Context(), "request-id", reqId))
		w.Header().Set("X-Request-ID", reqId)
		next.ServeHTTP(w,r)
	})
}

func handleRootPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		w.Write([]byte("Welcome\n"))
	default:	
		http.Error(w, "Error processing you request", http.StatusMethodNotAllowed)
	}
}

func handleHelloPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := r.URL.Query().Get("name")
	switch r.Method {
	case "GET":
		if name == "" {
			http.Error(w, "Error processing you request", http.StatusBadRequest)
			return
		}
		w.Write([]byte(fmt.Sprintf("Hello, %s!\n", name)))
	default:	
		http.Error(w, "Error processing you request", http.StatusMethodNotAllowed)
	}
}

func handleEchoPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(r.Header)
	case "POST":
		body, _ := io.ReadAll(r.Body)
		w.Write(body)
	default:	
		http.Error(w, "Error processing you request", http.StatusMethodNotAllowed)
	}
}

func handleStatusPath(w http.ResponseWriter, r *http.Request) {
	serverStats.mu.RLock()
	response := ServerStats{
		Uptime: time.Since(serverStats.startTime).String(),
		RequestCount: serverStats.requestCount,
	}
	serverStats.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(response)
	default:	
		http.Error(w, "Error processing you request", http.StatusMethodNotAllowed)
	}
}

func (sr  *StatusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader((code))
}