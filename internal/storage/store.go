package storage

import "github.com/blackms/ExplainableEngine/internal/models"

// ListOptions controls filtering and pagination for listing explanations.
type ListOptions struct {
	Cursor        string  // last seen ID for cursor pagination
	Limit         int     // max items (default 20, max 100)
	Target        string  // filter by target name (substring match)
	MinConfidence float64 // filter: confidence >= this
	MaxConfidence float64 // filter: confidence <= this (0 = no filter)
	FromTime      string  // ISO8601 timestamp filter
	ToTime        string  // ISO8601 timestamp filter
}

// ListResult is the paginated response for listing explanations.
type ListResult struct {
	Items      []*models.ExplainResponse `json:"items"`
	NextCursor string                    `json:"next_cursor,omitempty"`
	Total      int                       `json:"total"`
}

// ExplanationStore defines the persistence interface for explanation responses.
type ExplanationStore interface {
	// Save persists an ExplainResponse. The response's ID field is used as the key.
	Save(resp *models.ExplainResponse) error

	// Get retrieves an ExplainResponse by ID. Returns nil, nil if not found.
	Get(id string) (*models.ExplainResponse, error)

	// Exists checks whether an explanation with the given ID is stored.
	Exists(id string) (bool, error)

	// List returns a paginated, filtered list of explanations.
	List(opts ListOptions) (*ListResult, error)

	// Count returns the total number of stored explanations.
	Count() (int, error)
}
