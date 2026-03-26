package storage

import "github.com/blackms/ExplainableEngine/internal/models"

// ExplanationStore defines the persistence interface for explanation responses.
type ExplanationStore interface {
	// Save persists an ExplainResponse. The response's ID field is used as the key.
	Save(resp *models.ExplainResponse) error

	// Get retrieves an ExplainResponse by ID. Returns nil, nil if not found.
	Get(id string) (*models.ExplainResponse, error)

	// Exists checks whether an explanation with the given ID is stored.
	Exists(id string) (bool, error)
}
