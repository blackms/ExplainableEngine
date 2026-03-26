package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/blackms/ExplainableEngine/internal/engine"
	"github.com/blackms/ExplainableEngine/internal/storage"
	"github.com/google/uuid"
)

// NewRouter creates the HTTP router with all routes and middleware.
func NewRouter(store storage.ExplanationStore, orch engine.OrchestratorInterface) http.Handler {
	handler := &ExplainHandler{orchestrator: orch, store: store}
	graphHandler := &GraphHandler{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /api/v1/explain", handler.Create)
	mux.HandleFunc("GET /api/v1/explain/{id}/graph", graphHandler.Export)
	mux.HandleFunc("GET /api/v1/explain/{id}", handler.Get)

	whatIfHandler := &WhatIfHandler{store: store, orchestrator: orch}
	mux.HandleFunc("POST /api/v1/explain/{id}/what-if", whatIfHandler.Analyze)

	return requestIDMiddleware(timingMiddleware(recoveryMiddleware(mux)))
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"version": "0.1.0",
	})
}

// --- Middleware ---

// recoveryMiddleware catches panics and returns a 500 response.
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic recovered: %v", rec)
				writeError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// requestIDMiddleware adds a unique X-Request-Id header to each response.
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		w.Header().Set("X-Request-Id", reqID)
		next.ServeHTTP(w, r)
	})
}

// timingMiddleware records the processing time and sets X-Processing-Time-Ms.
func timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start).Milliseconds()
		w.Header().Set("X-Processing-Time-Ms", fmt.Sprintf("%d", elapsed))
	})
}

// --- Response Helpers ---

// writeJSON encodes data as JSON and writes it with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}

// writeError writes a JSON error response with the given status code and message.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
