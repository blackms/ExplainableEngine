package llm

import (
	"context"
	"fmt"

	"github.com/blackms/ExplainableEngine/internal/engine"
	"github.com/blackms/ExplainableEngine/internal/models"
)

// FallbackService uses the existing template-based narrative engine when the
// Claude API is unavailable. Q&A and summary operations are not supported.
type FallbackService struct{}

// NewFallbackService creates a new FallbackService.
func NewFallbackService() *FallbackService {
	return &FallbackService{}
}

// IsAvailable always returns true for the fallback service.
func (f *FallbackService) IsAvailable() bool { return true }

// GenerateNarrative delegates to the existing template-based engine.
func (f *FallbackService) GenerateNarrative(_ context.Context, e *models.ExplainResponse, level, lang string) (string, error) {
	narrativeLevel := engine.NarrativeLevel(level)
	narrativeLang := engine.NarrativeLanguage(lang)

	// The template engine only supports basic and advanced levels.
	if narrativeLevel != engine.LevelBasic && narrativeLevel != engine.LevelAdvanced {
		narrativeLevel = engine.LevelBasic
	}
	if narrativeLang != engine.LangEN && narrativeLang != engine.LangIT {
		narrativeLang = engine.LangEN
	}

	result, err := engine.GenerateNarrative(e, narrativeLevel, narrativeLang)
	if err != nil {
		return "", err
	}
	return result.Narrative, nil
}

// AnswerQuestion is not supported without an LLM.
func (f *FallbackService) AnswerQuestion(_ context.Context, _ *models.ExplainResponse, _ string, _ []Message) (string, error) {
	return "", fmt.Errorf("Q&A requires LLM — set ANTHROPIC_API_KEY")
}

// GenerateSummary is not supported without an LLM.
func (f *FallbackService) GenerateSummary(_ context.Context, _ *models.ExplainResponse, _ string, _ string) (*SummaryResult, error) {
	return nil, fmt.Errorf("executive summary requires LLM — set ANTHROPIC_API_KEY")
}
