package llm

import (
	"context"

	"github.com/blackms/ExplainableEngine/internal/models"
)

// Service defines the LLM operations.
type Service interface {
	GenerateNarrative(ctx context.Context, explanation *models.ExplainResponse, level, lang string) (string, error)
	AnswerQuestion(ctx context.Context, explanation *models.ExplainResponse, question string, history []Message) (string, error)
	GenerateSummary(ctx context.Context, explanation *models.ExplainResponse, audience, lang string) (*SummaryResult, error)
	IsAvailable() bool
}

// Message represents a single message in a conversation history.
type Message struct {
	Role    string `json:"role"`    // "user" or "assistant"
	Content string `json:"content"`
}

// SummaryResult holds the structured output of an executive summary.
type SummaryResult struct {
	Title           string   `json:"title"`
	Summary         string   `json:"summary"`
	KeyFindings     []string `json:"key_findings"`
	Risks           []string `json:"risks"`
	Recommendations []string `json:"recommendations"`
	Audience        string   `json:"audience"`
	Language        string   `json:"language"`
}
