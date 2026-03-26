package main

import (
	"log"
	"net/http"
	"os"

	"github.com/blackms/ExplainableEngine/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := api.NewRouter()

	log.Printf("Explainable Engine starting on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
