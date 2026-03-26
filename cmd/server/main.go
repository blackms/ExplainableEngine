package main

import (
	"log"
	"net/http"
	"os"

	"github.com/blackms/ExplainableEngine/internal/api"
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

	store, err := storage.NewStore(backend, sqlitePath)
	if err != nil {
		log.Fatalf("failed to create store: %v", err)
	}

	// The orchestrator is nil for now; it will be wired in when the engine
	// package provides a concrete implementation. The API will return 500
	// if POST /api/v1/explain is called without one.
	router := api.NewRouter(store, nil)

	log.Printf("Explainable Engine starting on :%s (storage=%s)", port, backend)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
