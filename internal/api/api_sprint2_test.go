package api_test

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/api"
	"github.com/blackms/ExplainableEngine/internal/engine"
	"github.com/blackms/ExplainableEngine/internal/models"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

// newRealTestServer creates a test server using the real orchestrator (no mocks).
func newRealTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	store := storage.NewInMemoryStore()
	orch := engine.NewOrchestrator()
	router := api.NewRouter(store, orch)
	return httptest.NewServer(router)
}

// postExplain sends a POST /api/v1/explain request and returns the decoded response.
func postExplain(t *testing.T, ts *httptest.Server, req models.ExplainRequest) (*http.Response, models.ExplainResponse) {
	t.Helper()
	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp, err := http.Post(ts.URL+"/api/v1/explain", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST /api/v1/explain: %v", err)
	}

	var result models.ExplainResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	resp.Body.Close()

	return resp, result
}

// TestPostExplainWithMissingNode verifies that a request with a missing component
// produces a response with missing_impact > 0.
func TestPostExplainWithMissingNode(t *testing.T) {
	ts := newRealTestServer(t)
	defer ts.Close()

	req := models.ExplainRequest{
		Target: "score",
		Value:  100.0,
		Components: []models.Component{
			{Name: "revenue", Value: 80.0, Weight: 0.6, Confidence: 0.9},
			{Name: "cost", Value: 0.0, Weight: 0.4, Confidence: 0.0, Missing: true},
		},
	}

	resp, result := postExplain(t, ts, req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if result.ID == "" {
		t.Error("response ID should not be empty")
	}
	if result.Target != "score" {
		t.Errorf("target: got %s, want score", result.Target)
	}
	if result.MissingImpact <= 0 {
		t.Errorf("missing_impact should be > 0 when a component is missing, got %f", result.MissingImpact)
	}
}

// TestPostExplainWithMissingThreshold verifies that a request with a missing node
// and a custom missing_threshold in options produces a valid response with impact data.
func TestPostExplainWithMissingThreshold(t *testing.T) {
	ts := newRealTestServer(t)
	defer ts.Close()

	req := models.ExplainRequest{
		Target: "score",
		Value:  100.0,
		Components: []models.Component{
			{Name: "revenue", Value: 80.0, Weight: 0.6, Confidence: 0.9},
			{Name: "cost", Value: 0.0, Weight: 0.4, Confidence: 0.0, Missing: true},
		},
		Options: &models.ExplainOptions{
			MissingThreshold: 0.1,
		},
	}

	resp, result := postExplain(t, ts, req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if result.MissingImpact <= 0 {
		t.Errorf("missing_impact should be > 0 with missing component, got %f", result.MissingImpact)
	}
	// The missing component has weight 0.4 out of total 1.0, so impact should be 0.4.
	expectedImpact := 0.4
	if math.Abs(result.MissingImpact-expectedImpact) > 0.01 {
		t.Errorf("missing_impact: got %f, want ~%f", result.MissingImpact, expectedImpact)
	}
}

// TestPostExplainDriversNormalized verifies that the top driver has a normalized
// impact of 1.0 (the highest-impact driver is always normalized to 1.0).
func TestPostExplainDriversNormalized(t *testing.T) {
	ts := newRealTestServer(t)
	defer ts.Close()

	req := models.ExplainRequest{
		Target: "score",
		Value:  100.0,
		Components: []models.Component{
			{Name: "revenue", Value: 80.0, Weight: 0.6, Confidence: 0.9},
			{Name: "growth", Value: 20.0, Weight: 0.4, Confidence: 0.8},
		},
	}

	resp, result := postExplain(t, ts, req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if len(result.TopDrivers) == 0 {
		t.Fatal("expected at least one top driver")
	}
	topDriver := result.TopDrivers[0]
	if topDriver.Impact != 1.0 {
		t.Errorf("top driver impact should be normalized to 1.0, got %f", topDriver.Impact)
	}
	if topDriver.Rank != 1 {
		t.Errorf("top driver rank: got %d, want 1", topDriver.Rank)
	}

	// If there are multiple drivers, subsequent ones should have impact <= 1.0.
	for i, d := range result.TopDrivers {
		if d.Impact > 1.0 || d.Impact < 0.0 {
			t.Errorf("driver[%d] impact %f out of [0, 1] range", i, d.Impact)
		}
	}
}

// TestPostExplainNoMissingNodes verifies that when no components are marked as
// missing, the missing_impact is 0.0.
func TestPostExplainNoMissingNodes(t *testing.T) {
	ts := newRealTestServer(t)
	defer ts.Close()

	req := models.ExplainRequest{
		Target: "score",
		Value:  100.0,
		Components: []models.Component{
			{Name: "revenue", Value: 80.0, Weight: 0.6, Confidence: 0.9},
			{Name: "growth", Value: 20.0, Weight: 0.4, Confidence: 0.8},
		},
	}

	resp, result := postExplain(t, ts, req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if result.MissingImpact != 0.0 {
		t.Errorf("missing_impact should be 0.0 when no nodes are missing, got %f", result.MissingImpact)
	}
}

// TestPostExplainBackwardCompat verifies that a Sprint 1 style request (no missing
// flags, no options) still produces a valid response with the expected structure.
func TestPostExplainBackwardCompat(t *testing.T) {
	ts := newRealTestServer(t)
	defer ts.Close()

	// Sprint 1 style request: no Missing field, no Options.
	req := models.ExplainRequest{
		Target: "score",
		Value:  0.72,
		Components: []models.Component{
			{Name: "trend", Value: 0.8, Weight: 0.4, Confidence: 0.9},
			{Name: "sentiment", Value: 0.6, Weight: 0.6, Confidence: 0.85},
		},
	}

	resp, result := postExplain(t, ts, req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Verify all Sprint 1 fields are present and valid.
	if result.ID == "" {
		t.Error("response ID should not be empty")
	}
	if result.Target != "score" {
		t.Errorf("target: got %s, want score", result.Target)
	}
	if result.FinalValue != 0.72 {
		t.Errorf("final_value: got %f, want 0.72", result.FinalValue)
	}
	if len(result.Breakdown) == 0 {
		t.Error("breakdown should not be empty")
	}
	if len(result.TopDrivers) == 0 {
		t.Error("top_drivers should not be empty")
	}
	if result.Confidence <= 0 {
		t.Errorf("confidence should be > 0, got %f", result.Confidence)
	}
	if result.MissingImpact != 0.0 {
		t.Errorf("missing_impact should be 0.0 for backward compat request, got %f", result.MissingImpact)
	}
	if result.Metadata.Version == "" {
		t.Error("metadata.version should not be empty")
	}
	if result.Metadata.DeterministicHash == "" {
		t.Error("metadata.deterministic_hash should not be empty")
	}
}
