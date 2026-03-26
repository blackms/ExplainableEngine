package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/blackms/ExplainableEngine/internal/api"
	"github.com/blackms/ExplainableEngine/internal/engine"
	"github.com/blackms/ExplainableEngine/internal/llm"
	"github.com/blackms/ExplainableEngine/internal/middleware"
	"github.com/blackms/ExplainableEngine/internal/storage"
	"github.com/blackms/ExplainableEngine/migrations"
)

func main() {
	migrateFlag := flag.String("migrate", "", "Run migrations and exit: up or down (requires STORAGE_BACKEND=postgresql)")
	flag.Parse()

	if *migrateFlag != "" {
		dir := storage.MigrationDirection(*migrateFlag)
		if dir != storage.MigrateUp && dir != storage.MigrateDown {
			log.Fatalf("invalid migration direction %q: must be 'up' or 'down'", *migrateFlag)
		}
		dsn := storage.BuildPostgresDSN()
		if err := storage.RunMigrations(dsn, migrations.FS, dir); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
		log.Printf("Migrations (%s) completed successfully", dir)
		return
	}

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

	// Initialize LLM service.
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	var llmService llm.Service
	if apiKey != "" {
		svc, err := llm.NewClaudeService(apiKey)
		if err != nil {
			log.Printf("LLM service unavailable: %v", err)
			llmService = llm.NewFallbackService()
		} else {
			llmService = svc
			log.Println("LLM service: Claude API enabled")
		}
	} else {
		llmService = llm.NewFallbackService()
		log.Println("LLM service: template fallback (no ANTHROPIC_API_KEY)")
	}

	router := api.NewRouter(store, orch, api.WithLLMService(llmService))

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
