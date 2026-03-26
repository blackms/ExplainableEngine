package main

import (
	"log"
	"net/http"
	"os"

	"github.com/blackms/ExplainableEngine/internal/api"
	"github.com/blackms/ExplainableEngine/internal/engine"
	"github.com/blackms/ExplainableEngine/internal/middleware"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	backend := os.Getenv("STORAGE_BACKEND")
	if backend == "" {
		backend = "memory"
	}

	sqlitePath := os.Getenv("SQLITE_PATH")
	if sqlitePath == "" {
		sqlitePath = "explanations.db"
	}

	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "*"
	}

	store, err := storage.NewStore(backend, sqlitePath)
	if err != nil {
		log.Fatalf("failed to create store: %v", err)
	}

	orch := engine.NewOrchestrator()
	router := api.NewRouter(store, orch)

	// Wrap router with CORS and structured logging middleware.
	// Order: CORS (outermost) -> Logging -> router (innermost with recovery/requestID/timing).
	handler := middleware.CORSMiddleware(corsOrigins)(
		middleware.LoggingMiddleware(router),
	)

	log.Printf("Explainable Engine starting on :%s (storage=%s)", port, backend)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
