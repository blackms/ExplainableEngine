package engine

import "github.com/blackms/ExplainableEngine/internal/models"

// OrchestratorInterface defines the contract for explanation orchestration.
// This interface exists so that API and test code can depend on it without
// importing the full engine implementation.
type OrchestratorInterface interface {
	Explain(req models.ExplainRequest) (*models.ExplainResponse, error)
}
