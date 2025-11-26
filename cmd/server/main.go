package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	defaultPort     = 8080
	serviceName     = "go-platform-starter-service"
	serviceVersion  = "0.1.0"
	shutdownTimeout = 5 * time.Second
)

// struct for health endpoint response
type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// struct for version endpoint response
type versionResponse struct {
	Service string `json:"service"`
	Version string `json:"version"`
}

func main() {
	// determine which port to run on
	port := getPort()

	// create new http router
	mux := http.NewServeMux()
	// register health and version endpoints for logging
	mux.Handle("/health", loggingMiddleware(http.HandlerFunc(handleHealth)))
	mux.Handle("/version", loggingMiddleware(http.HandlerFunc(handleVersion)))

	// http server config
	server := &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// start server in separate goroutine so main can listen for shutdown signals
	go func() {
		log.Printf("starting %s on port %d", serviceName, port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// wait for shutdown signals and handle graceful shutdown
	waitForShutdown(server)
}

// reads PORT env var, falls back to default if invalid or missing
func getPort() int {
	if portStr := os.Getenv("PORT"); portStr != "" {
		if portVal, err := strconv.Atoi(portStr); err == nil && portVal > 0 {
			return portVal
		}
		log.Printf("invalid PORT value %q, falling back to default %d", portStr, defaultPort)
	}
	return defaultPort
}

// returns a simple ok status and service name
func handleHealth(w http.ResponseWriter, r *http.Request) {
	resp := healthResponse{
		Status:  "ok",
		Service: serviceName,
	}
	writeJSON(w, http.StatusOK, resp)
}

// returns service name and version
func handleVersion(w http.ResponseWriter, r *http.Request) {
	resp := versionResponse{
		Service: serviceName,
		Version: serviceVersion,
	}
	writeJSON(w, http.StatusOK, resp)
}

// encodes payload as json and writes it to response
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write json response: %v", err)
	}
}

// wraps handlers to log request start and completion
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("completed %s %s in %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// listens for SIGINT/SIGTERM and shuts down server gracefully
func waitForShutdown(server *http.Server) {
	// channel to receive os signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// block until a signal is received
	<-stop
	log.Printf("shutdown signal received, shutting down %s...", serviceName)

	// create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		// fallback to forced close
		if err := server.Close(); err != nil {
			log.Printf("forced close failed: %v", err)
		}
	} else {
		log.Printf("%s shut down cleanly", serviceName)
	}
}
