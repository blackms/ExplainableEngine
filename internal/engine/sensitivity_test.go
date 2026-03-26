package engine

import (
	"errors"
	"math"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/models"
)

func makeSensitivityRequest() models.ExplainRequest {
	// Value = 80*0.6 + 20*0.4 = 56.0 (matches weighted sum of components).
	return models.ExplainRequest{
		Target: "score",
		Value:  56.0,
		Components: []models.Component{
			{Name: "revenue", Value: 80.0, Weight: 0.6, Confidence: 0.9},
			{Name: "growth", Value: 20.0, Weight: 0.4, Confidence: 0.8},
		},
	}
}

func TestAnalyzeSensitivity_SingleModification(t *testing.T) {
	orch := NewOrchestrator()
	req := makeSensitivityRequest()

	origResp, err := orch.Explain(req)
	if err != nil {
		t.Fatalf("Explain failed: %v", err)
	}

	mods := []Modification{
		{ComponentName: "revenue", NewValue: 90.0},
	}

	result, err := AnalyzeSensitivity(&req, origResp, mods, orch)
	if err != nil {
		t.Fatalf("AnalyzeSensitivity failed: %v", err)
	}

	if result.OriginalValue != origResp.FinalValue {
		t.Errorf("original value: got %f, want %f", result.OriginalValue, origResp.FinalValue)
	}
	if result.DeltaValue == 0 {
		t.Error("delta value should not be 0 after modification")
	}
	if result.ModifiedValue != result.OriginalValue+result.DeltaValue {
		t.Error("modified value should equal original + delta")
	}
	if len(result.ComponentDiffs) == 0 {
		t.Error("component diffs should not be empty")
	}

	// Find the revenue diff.
	var revenueDiff *ComponentDiff
	for i, d := range result.ComponentDiffs {
		if d.Name == "revenue" {
			revenueDiff = &result.ComponentDiffs[i]
			break
		}
	}
	if revenueDiff == nil {
		t.Fatal("revenue diff not found")
	}
	if revenueDiff.OriginalValue != 80.0 {
		t.Errorf("revenue original value: got %f, want 80.0", revenueDiff.OriginalValue)
	}
	if revenueDiff.ModifiedValue != 90.0 {
		t.Errorf("revenue modified value: got %f, want 90.0", revenueDiff.ModifiedValue)
	}
	if revenueDiff.DeltaValue != 10.0 {
		t.Errorf("revenue delta value: got %f, want 10.0", revenueDiff.DeltaValue)
	}
}

func TestAnalyzeSensitivity_MultipleModifications(t *testing.T) {
	orch := NewOrchestrator()
	req := makeSensitivityRequest()

	origResp, err := orch.Explain(req)
	if err != nil {
		t.Fatalf("Explain failed: %v", err)
	}

	mods := []Modification{
		{ComponentName: "revenue", NewValue: 90.0},
		{ComponentName: "growth", NewValue: 30.0},
	}

	result, err := AnalyzeSensitivity(&req, origResp, mods, orch)
	if err != nil {
		t.Fatalf("AnalyzeSensitivity failed: %v", err)
	}

	if result.DeltaValue == 0 {
		t.Error("delta value should not be 0 after multiple modifications")
	}

	// Both components should appear in diffs.
	names := make(map[string]bool)
	for _, d := range result.ComponentDiffs {
		names[d.Name] = true
	}
	if !names["revenue"] {
		t.Error("revenue should appear in component diffs")
	}
	if !names["growth"] {
		t.Error("growth should appear in component diffs")
	}
}

func TestAnalyzeSensitivity_NoModification(t *testing.T) {
	orch := NewOrchestrator()
	req := makeSensitivityRequest()

	origResp, err := orch.Explain(req)
	if err != nil {
		t.Fatalf("Explain failed: %v", err)
	}

	result, err := AnalyzeSensitivity(&req, origResp, []Modification{}, orch)
	if err != nil {
		t.Fatalf("AnalyzeSensitivity failed: %v", err)
	}

	if result.DeltaValue != 0 {
		t.Errorf("delta value should be 0 with no modifications, got %f", result.DeltaValue)
	}
	if result.DeltaPercentage != 0 {
		t.Errorf("delta percentage should be 0 with no modifications, got %f", result.DeltaPercentage)
	}
	if result.OriginalValue != result.ModifiedValue {
		t.Errorf("original and modified values should be equal with no modifications")
	}
}

func TestAnalyzeSensitivity_InvalidComponent(t *testing.T) {
	orch := NewOrchestrator()
	req := makeSensitivityRequest()

	origResp, err := orch.Explain(req)
	if err != nil {
		t.Fatalf("Explain failed: %v", err)
	}

	mods := []Modification{
		{ComponentName: "nonexistent", NewValue: 50.0},
	}

	_, err = AnalyzeSensitivity(&req, origResp, mods, orch)
	if err == nil {
		t.Fatal("expected error for invalid component name")
	}

	var cnfErr *ComponentNotFoundError
	if !errors.As(err, &cnfErr) {
		t.Fatalf("expected ComponentNotFoundError, got %T: %v", err, err)
	}
	if cnfErr.ComponentName != "nonexistent" {
		t.Errorf("component name: got %s, want nonexistent", cnfErr.ComponentName)
	}
}

func TestAnalyzeSensitivity_Determinism(t *testing.T) {
	orch := NewOrchestrator()
	req := makeSensitivityRequest()

	origResp, err := orch.Explain(req)
	if err != nil {
		t.Fatalf("Explain failed: %v", err)
	}

	mods := []Modification{
		{ComponentName: "revenue", NewValue: 90.0},
	}

	first, err := AnalyzeSensitivity(&req, origResp, mods, orch)
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}

	for i := 0; i < 10; i++ {
		result, err := AnalyzeSensitivity(&req, origResp, mods, orch)
		if err != nil {
			t.Fatalf("iteration %d failed: %v", i, err)
		}
		if result.DeltaValue != first.DeltaValue {
			t.Errorf("iteration %d: delta value mismatch: got %f, want %f", i, result.DeltaValue, first.DeltaValue)
		}
		if result.DeltaPercentage != first.DeltaPercentage {
			t.Errorf("iteration %d: delta percentage mismatch: got %f, want %f", i, result.DeltaPercentage, first.DeltaPercentage)
		}
		if result.ModifiedValue != first.ModifiedValue {
			t.Errorf("iteration %d: modified value mismatch: got %f, want %f", i, result.ModifiedValue, first.ModifiedValue)
		}
		if len(result.Ranking) != len(first.Ranking) {
			t.Fatalf("iteration %d: ranking length mismatch", i)
		}
		for j := range result.Ranking {
			if result.Ranking[j] != first.Ranking[j] {
				t.Errorf("iteration %d: ranking[%d] mismatch", i, j)
			}
		}
	}
}

func TestAnalyzeSensitivity_RankingOrder(t *testing.T) {
	orch := NewOrchestrator()
	req := makeSensitivityRequest()

	origResp, err := orch.Explain(req)
	if err != nil {
		t.Fatalf("Explain failed: %v", err)
	}

	// Make a big change to revenue (weight 0.6) and small change to growth (weight 0.4).
	mods := []Modification{
		{ComponentName: "revenue", NewValue: 180.0},
		{ComponentName: "growth", NewValue: 21.0},
	}

	result, err := AnalyzeSensitivity(&req, origResp, mods, orch)
	if err != nil {
		t.Fatalf("AnalyzeSensitivity failed: %v", err)
	}

	if len(result.Ranking) < 2 {
		t.Fatalf("expected at least 2 ranking entries, got %d", len(result.Ranking))
	}

	// Verify ranking is sorted by impact descending.
	for i := 1; i < len(result.Ranking); i++ {
		if result.Ranking[i].Impact > result.Ranking[i-1].Impact {
			t.Errorf("ranking not sorted: [%d].Impact=%f > [%d].Impact=%f",
				i, result.Ranking[i].Impact, i-1, result.Ranking[i-1].Impact)
		}
	}

	// Verify rank numbers are sequential.
	for i, r := range result.Ranking {
		if r.Rank != i+1 {
			t.Errorf("rank: got %d, want %d", r.Rank, i+1)
		}
	}

	// Revenue should have higher impact (larger change * larger weight).
	if result.Ranking[0].Name != "revenue" {
		t.Errorf("top ranked component: got %s, want revenue", result.Ranking[0].Name)
	}
}

func TestAnalyzeSensitivity_DeltaPercentage(t *testing.T) {
	orch := NewOrchestrator()
	req := makeSensitivityRequest()

	origResp, err := orch.Explain(req)
	if err != nil {
		t.Fatalf("Explain failed: %v", err)
	}

	mods := []Modification{
		{ComponentName: "revenue", NewValue: 90.0},
	}

	result, err := AnalyzeSensitivity(&req, origResp, mods, orch)
	if err != nil {
		t.Fatalf("AnalyzeSensitivity failed: %v", err)
	}

	// Verify delta percentage calculation.
	expectedPct := (result.DeltaValue / result.OriginalValue) * 100
	if math.Abs(result.DeltaPercentage-expectedPct) > 0.001 {
		t.Errorf("delta percentage: got %f, want %f", result.DeltaPercentage, expectedPct)
	}
}

func TestAnalyzeSensitivity_OriginalRequestUnchanged(t *testing.T) {
	orch := NewOrchestrator()
	req := makeSensitivityRequest()
	origValue := req.Components[0].Value

	origResp, err := orch.Explain(req)
	if err != nil {
		t.Fatalf("Explain failed: %v", err)
	}

	mods := []Modification{
		{ComponentName: "revenue", NewValue: 999.0},
	}

	_, err = AnalyzeSensitivity(&req, origResp, mods, orch)
	if err != nil {
		t.Fatalf("AnalyzeSensitivity failed: %v", err)
	}

	// Original request should not be mutated.
	if req.Components[0].Value != origValue {
		t.Errorf("original request was mutated: got %f, want %f", req.Components[0].Value, origValue)
	}
}
