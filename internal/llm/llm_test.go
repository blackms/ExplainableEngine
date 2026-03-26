package llm

import (
	"context"
	"strings"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/models"
)

func sampleExplanation() *models.ExplainResponse {
	return &models.ExplainResponse{
		ID:         "test-id-123",
		Target:     "credit_score",
		FinalValue: 0.72,
		Confidence: 0.85,
		Breakdown: []models.BreakdownItem{
			{
				NodeID:               "comp-a",
				Label:                "income",
				Value:                0.8,
				Weight:               0.4,
				AbsoluteContribution: 0.32,
				Percentage:           44.4,
				Confidence:           0.9,
			},
			{
				NodeID:               "comp-b",
				Label:                "history",
				Value:                0.6,
				Weight:               0.3,
				AbsoluteContribution: 0.18,
				Percentage:           25.0,
				Confidence:           0.8,
			},
		},
		TopDrivers: []models.DriverItem{
			{Name: "income", Impact: 0.44, Rank: 1},
			{Name: "history", Impact: 0.25, Rank: 2},
		},
		MissingImpact: 0.05,
		Graph: &models.GraphResponse{
			Nodes: []models.GraphNodeResponse{
				{ID: "root", Label: "credit_score", Value: 0.72, Confidence: 0.85, NodeType: "root"},
				{ID: "comp-a", Label: "income", Value: 0.8, Confidence: 0.9, NodeType: "component"},
			},
			Edges: []models.GraphEdgeResponse{
				{Source: "comp-a", Target: "root", Weight: 0.4, TransformationType: "weighted_sum"},
			},
		},
		DependencyTree: &models.DependencyTree{
			Root:       models.DependencyNode{NodeID: "root", Label: "credit_score", Depth: 0},
			Depth:      2,
			TotalNodes: 3,
		},
	}
}

func TestBuildExplanationContext(t *testing.T) {
	ctx := BuildExplanationContext(sampleExplanation())

	checks := []string{
		"Target: credit_score",
		"Final Value: 0.7200",
		"Overall Confidence: 85.0%",
		"Missing Data Impact: 5.0%",
		"income",
		"history",
		"#1 income",
		"#2 history",
		"2 nodes",
		"1 edges",
		"depth=2",
	}

	for _, check := range checks {
		if !strings.Contains(ctx, check) {
			t.Errorf("context missing %q\nGot:\n%s", check, ctx)
		}
	}
}

func TestBuildExplanationContext_Empty(t *testing.T) {
	e := &models.ExplainResponse{
		Target:     "empty",
		FinalValue: 0,
		Confidence: 0,
	}
	ctx := BuildExplanationContext(e)
	if !strings.Contains(ctx, "Target: empty") {
		t.Errorf("context should contain target, got: %s", ctx)
	}
	// Should not contain breakdown or driver sections.
	if strings.Contains(ctx, "Breakdown:") {
		t.Error("context should not contain Breakdown section for empty response")
	}
	if strings.Contains(ctx, "Top Drivers:") {
		t.Error("context should not contain Top Drivers section for empty response")
	}
}

func TestBuildNarrativeSystemPrompt_Levels(t *testing.T) {
	tests := []struct {
		level    string
		lang     string
		contains string
	}{
		{"basic", "en", "non-technical stakeholder"},
		{"advanced", "en", "quantitative analyst"},
		{"executive", "en", "financial analyst"},
		{"basic", "it", "Rispondi in italiano"},
		{"advanced", "it", "Rispondi in italiano"},
		{"executive", "it", "Rispondi in italiano"},
		{"basic", "en", "Respond in English"},
		{"unknown", "en", "non-technical stakeholder"}, // defaults to basic
	}

	for _, tc := range tests {
		prompt := buildNarrativeSystemPrompt(tc.level, tc.lang)
		if !strings.Contains(prompt, tc.contains) {
			t.Errorf("buildNarrativeSystemPrompt(%q, %q) should contain %q, got: %s",
				tc.level, tc.lang, tc.contains, prompt)
		}
	}
}

func TestBuildSummarySystemPrompt_Audiences(t *testing.T) {
	tests := []struct {
		audience string
		lang     string
		contains string
	}{
		{"board", "en", "executive board"},
		{"technical", "en", "technical deep-dive"},
		{"client", "en", "client-facing"},
		{"board", "it", "Rispondi in italiano"},
		{"unknown", "en", "client-facing"}, // defaults to client
	}

	for _, tc := range tests {
		prompt := buildSummarySystemPrompt(tc.audience, tc.lang)
		if !strings.Contains(prompt, tc.contains) {
			t.Errorf("buildSummarySystemPrompt(%q, %q) should contain %q, got: %s",
				tc.audience, tc.lang, tc.contains, prompt)
		}
		// All prompts must include JSON schema instruction.
		if !strings.Contains(prompt, "JSON object") {
			t.Errorf("buildSummarySystemPrompt(%q, %q) should contain JSON schema instruction", tc.audience, tc.lang)
		}
	}
}

func TestNewClaudeService_EmptyKey(t *testing.T) {
	_, err := NewClaudeService("")
	if err == nil {
		t.Fatal("expected error for empty API key")
	}
	if !strings.Contains(err.Error(), "ANTHROPIC_API_KEY") {
		t.Errorf("error should mention ANTHROPIC_API_KEY, got: %v", err)
	}
}

func TestNewClaudeService_ValidKey(t *testing.T) {
	svc, err := NewClaudeService("sk-test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !svc.IsAvailable() {
		t.Error("service should be available")
	}
}

func TestFallbackService_IsAvailable(t *testing.T) {
	svc := NewFallbackService()
	if !svc.IsAvailable() {
		t.Error("fallback service should always be available")
	}
}

func TestFallbackService_GenerateNarrative(t *testing.T) {
	svc := NewFallbackService()
	narrative, err := svc.GenerateNarrative(context.Background(), sampleExplanation(), "basic", "en")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if narrative == "" {
		t.Error("narrative should not be empty")
	}
	if !strings.Contains(narrative, "score") {
		t.Errorf("narrative should contain target name, got: %s", narrative)
	}
}

func TestFallbackService_GenerateNarrative_Executive(t *testing.T) {
	// Executive level is not supported by template engine; fallback should
	// default to basic.
	svc := NewFallbackService()
	narrative, err := svc.GenerateNarrative(context.Background(), sampleExplanation(), "executive", "en")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if narrative == "" {
		t.Error("narrative should not be empty")
	}
}

func TestFallbackService_AnswerQuestion_Error(t *testing.T) {
	svc := NewFallbackService()
	_, err := svc.AnswerQuestion(context.Background(), sampleExplanation(), "why?", nil)
	if err == nil {
		t.Fatal("expected error from fallback AnswerQuestion")
	}
	if !strings.Contains(err.Error(), "ANTHROPIC_API_KEY") {
		t.Errorf("error should mention ANTHROPIC_API_KEY, got: %v", err)
	}
}

func TestFallbackService_GenerateSummary_Error(t *testing.T) {
	svc := NewFallbackService()
	_, err := svc.GenerateSummary(context.Background(), sampleExplanation(), "board", "en")
	if err == nil {
		t.Fatal("expected error from fallback GenerateSummary")
	}
	if !strings.Contains(err.Error(), "ANTHROPIC_API_KEY") {
		t.Errorf("error should mention ANTHROPIC_API_KEY, got: %v", err)
	}
}

func TestStripCodeFences(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{"key":"value"}`, `{"key":"value"}`},
		{"```json\n{\"key\":\"value\"}\n```", `{"key":"value"}`},
		{"```\n{\"key\":\"value\"}\n```", `{"key":"value"}`},
		{" {\"key\":\"value\"} ", `{"key":"value"}`},
	}

	for _, tc := range tests {
		got := stripCodeFences(tc.input)
		if got != tc.expected {
			t.Errorf("stripCodeFences(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}
