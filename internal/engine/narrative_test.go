package engine

import (
	"strings"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/models"
)

func sampleResponse() *models.ExplainResponse {
	return &models.ExplainResponse{
		ID:         "test-id-123",
		Target:     "score",
		FinalValue: 0.72,
		Confidence: 0.85,
		Breakdown: []models.BreakdownItem{
			{Label: "trend", Percentage: 44.4},
			{Label: "seasonality", Percentage: 30.2},
			{Label: "volume", Percentage: 15.1},
		},
		TopDrivers: []models.DriverItem{
			{Name: "trend", Impact: 0.44, Rank: 1},
			{Name: "seasonality", Impact: 0.30, Rank: 2},
			{Name: "volume", Impact: 0.15, Rank: 3},
		},
		MissingImpact: 0.05,
	}
}

func TestGenerateNarrative_BasicEN(t *testing.T) {
	resp := sampleResponse()
	result, err := GenerateNarrative(resp, LevelBasic, LangEN)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result.Narrative, "score") {
		t.Error("basic EN narrative should contain target name")
	}
	if !strings.Contains(result.Narrative, "0.72") {
		t.Error("basic EN narrative should contain final value")
	}
	if !strings.Contains(result.Narrative, "trend") {
		t.Error("basic EN narrative should contain top driver")
	}
	if !strings.Contains(result.Narrative, "high") {
		t.Error("basic EN narrative should contain confidence level")
	}
	if result.ConfidenceLevel != "high" {
		t.Errorf("confidence level: got %q, want %q", result.ConfidenceLevel, "high")
	}
	if result.ExplanationID != "test-id-123" {
		t.Errorf("explanation ID: got %q, want %q", result.ExplanationID, "test-id-123")
	}
}

func TestGenerateNarrative_BasicIT(t *testing.T) {
	resp := sampleResponse()
	result, err := GenerateNarrative(resp, LevelBasic, LangIT)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result.Narrative, "guidato principalmente") {
		t.Error("basic IT narrative should contain 'guidato principalmente'")
	}
	if !strings.Contains(result.Narrative, "La confidenza è") {
		t.Error("basic IT narrative should contain 'La confidenza è'")
	}
	if result.Language != LangIT {
		t.Errorf("language: got %q, want %q", result.Language, LangIT)
	}
}

func TestGenerateNarrative_AdvancedEN(t *testing.T) {
	resp := sampleResponse()
	result, err := GenerateNarrative(resp, LevelAdvanced, LangEN)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result.Narrative, "Key drivers:") {
		t.Error("advanced EN narrative should contain 'Key drivers:'")
	}

	// Should have 3 driver lines
	lines := strings.Split(result.Narrative, "\n")
	driverLines := 0
	for _, l := range lines {
		if strings.HasPrefix(l, "- ") {
			driverLines++
		}
	}
	if driverLines != 3 {
		t.Errorf("expected 3 driver lines, got %d", driverLines)
	}

	if !strings.Contains(result.Narrative, "trend") {
		t.Error("advanced EN narrative should contain driver name 'trend'")
	}
	if !strings.Contains(result.Narrative, "seasonality") {
		t.Error("advanced EN narrative should contain driver name 'seasonality'")
	}
	if !strings.Contains(result.Narrative, "volume") {
		t.Error("advanced EN narrative should contain driver name 'volume'")
	}
}

func TestGenerateNarrative_AdvancedIT(t *testing.T) {
	resp := sampleResponse()
	result, err := GenerateNarrative(resp, LevelAdvanced, LangIT)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result.Narrative, "Fattori principali:") {
		t.Error("advanced IT narrative should contain 'Fattori principali:'")
	}
	if !strings.Contains(result.Narrative, "contributo del") {
		t.Error("advanced IT narrative should contain 'contributo del'")
	}
}

func TestGenerateNarrative_HighConfidence(t *testing.T) {
	resp := sampleResponse()
	resp.Confidence = 0.9

	result, err := GenerateNarrative(resp, LevelBasic, LangEN)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ConfidenceLevel != "high" {
		t.Errorf("expected 'high', got %q", result.ConfidenceLevel)
	}

	resultIT, err := GenerateNarrative(resp, LevelBasic, LangIT)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resultIT.ConfidenceLevel != "alta" {
		t.Errorf("expected 'alta', got %q", resultIT.ConfidenceLevel)
	}
}

func TestGenerateNarrative_LowConfidence(t *testing.T) {
	resp := sampleResponse()
	resp.Confidence = 0.3

	result, err := GenerateNarrative(resp, LevelBasic, LangEN)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ConfidenceLevel != "low" {
		t.Errorf("expected 'low', got %q", result.ConfidenceLevel)
	}

	resultIT, err := GenerateNarrative(resp, LevelBasic, LangIT)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resultIT.ConfidenceLevel != "bassa" {
		t.Errorf("expected 'bassa', got %q", resultIT.ConfidenceLevel)
	}
}

func TestGenerateNarrative_MissingDataWarning(t *testing.T) {
	resp := sampleResponse()
	resp.MissingImpact = 0.15

	result, err := GenerateNarrative(resp, LevelAdvanced, LangEN)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasMissingData {
		t.Error("HasMissingData should be true")
	}
	if !strings.Contains(result.Narrative, "15.0% of input data is missing") {
		t.Error("advanced EN narrative should contain missing data warning")
	}

	resultIT, err := GenerateNarrative(resp, LevelAdvanced, LangIT)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resultIT.Narrative, "15.0% dei dati di input è mancante") {
		t.Error("advanced IT narrative should contain missing data warning in Italian")
	}
}

func TestGenerateNarrative_NoMissingData(t *testing.T) {
	resp := sampleResponse()
	resp.MissingImpact = 0.05

	result, err := GenerateNarrative(resp, LevelAdvanced, LangEN)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasMissingData {
		t.Error("HasMissingData should be false")
	}
	if strings.Contains(result.Narrative, "missing") {
		t.Error("narrative should not contain missing data warning")
	}
}

func TestGenerateNarrative_UnsupportedLevel(t *testing.T) {
	resp := sampleResponse()
	_, err := GenerateNarrative(resp, "expert", LangEN)
	if err == nil {
		t.Error("expected error for unsupported level")
	}
}

func TestGenerateNarrative_UnsupportedLanguage(t *testing.T) {
	resp := sampleResponse()
	_, err := GenerateNarrative(resp, LevelBasic, "fr")
	if err == nil {
		t.Error("expected error for unsupported language")
	}
}

func TestGenerateNarrative_ZeroDrivers(t *testing.T) {
	resp := sampleResponse()
	resp.TopDrivers = nil
	resp.Breakdown = nil

	result, err := GenerateNarrative(resp, LevelBasic, LangEN)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result.Narrative, "score") {
		t.Error("narrative with zero drivers should still contain target")
	}
	if !strings.Contains(result.Narrative, "0.72") {
		t.Error("narrative with zero drivers should still contain value")
	}

	// Advanced with no drivers should not panic
	resultAdv, err := GenerateNarrative(resp, LevelAdvanced, LangEN)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(resultAdv.Narrative, "Key drivers:") {
		t.Error("advanced narrative with zero drivers should not contain 'Key drivers:'")
	}
}
