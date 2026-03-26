package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/api"
	"github.com/blackms/ExplainableEngine/internal/engine"
	"github.com/blackms/ExplainableEngine/internal/models"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

// newRealTestServerSprint3 creates a test server using the real orchestrator.
func newRealTestServerSprint3(t *testing.T) *httptest.Server {
	t.Helper()
	store := storage.NewInMemoryStore()
	orch := engine.NewOrchestrator()
	router := api.NewRouter(store, orch)
	return httptest.NewServer(router)
}

// sprint3ExplainRequest returns a request whose Value matches the weighted sum.
// revenue: 80 * 0.6 = 48, growth: 20 * 0.4 = 8, total = 56.
func sprint3ExplainRequest() models.ExplainRequest {
	return models.ExplainRequest{
		Target: "score",
		Value:  56.0,
		Components: []models.Component{
			{Name: "revenue", Value: 80.0, Weight: 0.6, Confidence: 0.9},
			{Name: "growth", Value: 20.0, Weight: 0.4, Confidence: 0.8},
		},
	}
}

// createWhatIfExplanation posts an explain request and returns the response.
func createWhatIfExplanation(t *testing.T, ts *httptest.Server, req models.ExplainRequest) models.ExplainResponse {
	t.Helper()
	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	resp, err := http.Post(ts.URL+"/api/v1/explain", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST /api/v1/explain: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST /api/v1/explain: expected 200, got %d", resp.StatusCode)
	}

	var result models.ExplainResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return result
}

func TestWhatIf_ValidModification(t *testing.T) {
	ts := newRealTestServerSprint3(t)
	defer ts.Close()

	// Create an explanation first.
	explainResp := createWhatIfExplanation(t, ts, sprint3ExplainRequest())

	// Now do a what-if analysis.
	whatIfReq := struct {
		Modifications []engine.Modification `json:"modifications"`
	}{
		Modifications: []engine.Modification{
			{ComponentName: "revenue", NewValue: 90.0},
		},
	}
	body, _ := json.Marshal(whatIfReq)

	resp, err := http.Post(
		ts.URL+"/api/v1/explain/"+explainResp.ID+"/what-if",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("POST what-if: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody map[string]string
		json.NewDecoder(resp.Body).Decode(&errBody)
		t.Fatalf("expected 200, got %d: %v", resp.StatusCode, errBody)
	}

	var result engine.SensitivityResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode result: %v", err)
	}

	// Original value should be 56 (80*0.6 + 20*0.4).
	if result.OriginalValue != 56.0 {
		t.Errorf("original value: got %f, want 56.0", result.OriginalValue)
	}

	// Modified value should be 90*0.6 + 20*0.4 = 54+8 = 62.
	if result.ModifiedValue != 62.0 {
		t.Errorf("modified value: got %f, want 62.0", result.ModifiedValue)
	}

	// Delta should be 6.
	if result.DeltaValue != 6.0 {
		t.Errorf("delta value: got %f, want 6.0", result.DeltaValue)
	}

	if len(result.ComponentDiffs) == 0 {
		t.Error("component diffs should not be empty")
	}

	if len(result.Ranking) == 0 {
		t.Error("ranking should not be empty")
	}
}

func TestWhatIf_NonexistentExplanation(t *testing.T) {
	ts := newRealTestServerSprint3(t)
	defer ts.Close()

	whatIfReq := struct {
		Modifications []engine.Modification `json:"modifications"`
	}{
		Modifications: []engine.Modification{
			{ComponentName: "revenue", NewValue: 90.0},
		},
	}
	body, _ := json.Marshal(whatIfReq)

	resp, err := http.Post(
		ts.URL+"/api/v1/explain/nonexistent-id/what-if",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("POST what-if: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestWhatIf_UnknownComponent(t *testing.T) {
	ts := newRealTestServerSprint3(t)
	defer ts.Close()

	explainResp := createWhatIfExplanation(t, ts, sprint3ExplainRequest())

	whatIfReq := struct {
		Modifications []engine.Modification `json:"modifications"`
	}{
		Modifications: []engine.Modification{
			{ComponentName: "nonexistent_component", NewValue: 50.0},
		},
	}
	body, _ := json.Marshal(whatIfReq)

	resp, err := http.Post(
		ts.URL+"/api/v1/explain/"+explainResp.ID+"/what-if",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("POST what-if: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", resp.StatusCode)
	}
}

func TestWhatIf_EmptyBody(t *testing.T) {
	ts := newRealTestServerSprint3(t)
	defer ts.Close()

	explainResp := createWhatIfExplanation(t, ts, sprint3ExplainRequest())

	resp, err := http.Post(
		ts.URL+"/api/v1/explain/"+explainResp.ID+"/what-if",
		"application/json",
		bytes.NewReader([]byte{}),
	)
	if err != nil {
		t.Fatalf("POST what-if: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestWhatIf_ProcessingTimeHeader(t *testing.T) {
	// Use httptest.NewRecorder to verify middleware headers (real HTTP servers
	// flush headers before the timing middleware can set them).
	store := storage.NewInMemoryStore()
	orch := engine.NewOrchestrator()
	router := api.NewRouter(store, orch)

	// First create an explanation.
	explainBody, _ := json.Marshal(sprint3ExplainRequest())
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/explain", bytes.NewReader(explainBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)
	if createW.Code != http.StatusOK {
		t.Fatalf("POST /api/v1/explain failed: %d %s", createW.Code, createW.Body.String())
	}
	var created models.ExplainResponse
	json.NewDecoder(createW.Body).Decode(&created)

	// Now do the what-if.
	whatIfReq := struct {
		Modifications []engine.Modification `json:"modifications"`
	}{
		Modifications: []engine.Modification{
			{ComponentName: "revenue", NewValue: 90.0},
		},
	}
	body, _ := json.Marshal(whatIfReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain/"+created.ID+"/what-if", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	timing := w.Header().Get("X-Processing-Time-Ms")
	if timing == "" {
		t.Error("X-Processing-Time-Ms header should be present")
	}
}
