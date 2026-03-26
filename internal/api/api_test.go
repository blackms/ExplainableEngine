package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blackms/ExplainableEngine/internal/api"
	"github.com/blackms/ExplainableEngine/internal/models"
	"github.com/blackms/ExplainableEngine/internal/storage"
	"github.com/google/uuid"
)

// mockOrchestrator returns a fixed valid response for any request.
type mockOrchestrator struct{}

func (m *mockOrchestrator) Explain(req models.ExplainRequest) (*models.ExplainResponse, error) {
	return &models.ExplainResponse{
		ID:         uuid.New().String(),
		Target:     req.Target,
		FinalValue: 0.72,
		Confidence: 0.85,
		Breakdown: []models.BreakdownItem{
			{
				NodeID:               "comp-a",
				Label:                req.Components[0].Name,
				Value:                req.Components[0].Value,
				Weight:               req.Components[0].Weight,
				AbsoluteContribution: 0.32,
				Percentage:           44.4,
				Confidence:           0.9,
			},
		},
		TopDrivers: []models.DriverItem{
			{Name: req.Components[0].Name, Impact: 0.44, Rank: 1},
		},
		Metadata: models.ExplainMetadata{
			Version:   "0.1.0",
			CreatedAt: time.Now(),
		},
	}, nil
}

func newTestRouter() http.Handler {
	store := storage.NewInMemoryStore()
	orch := &mockOrchestrator{}
	return api.NewRouter(store, orch)
}

func validRequestBody() []byte {
	req := models.ExplainRequest{
		Target: "score",
		Value:  0.72,
		Components: []models.Component{
			{Name: "trend", Value: 0.8, Weight: 0.4, Confidence: 0.9},
		},
	}
	data, _ := json.Marshal(req)
	return data
}

func TestPostExplain_Valid(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain", bytes.NewReader(validRequestBody()))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp models.ExplainResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.ID == "" {
		t.Error("response ID should not be empty")
	}
	if resp.Target != "score" {
		t.Errorf("target: got %s, want score", resp.Target)
	}
	if resp.FinalValue != 0.72 {
		t.Errorf("final_value: got %f, want 0.72", resp.FinalValue)
	}
	if len(resp.Breakdown) != 1 {
		t.Errorf("breakdown length: got %d, want 1", len(resp.Breakdown))
	}
	if len(resp.TopDrivers) != 1 {
		t.Errorf("top_drivers length: got %d, want 1", len(resp.TopDrivers))
	}
}

func TestPostExplain_EmptyBody(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain", bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPostExplain_InvalidJSON(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain", bytes.NewReader([]byte("{not json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPostExplain_MissingTarget(t *testing.T) {
	body, _ := json.Marshal(models.ExplainRequest{
		Components: []models.Component{{Name: "a", Value: 1, Weight: 1, Confidence: 1}},
	})
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain", bytes.NewReader(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing target, got %d", w.Code)
	}
}

func TestPostExplain_EmptyComponents(t *testing.T) {
	body, _ := json.Marshal(models.ExplainRequest{
		Target: "score",
	})
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain", bytes.NewReader(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty components, got %d", w.Code)
	}
}

func TestGetExplain_AfterPost(t *testing.T) {
	router := newTestRouter()

	// POST first.
	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/explain", bytes.NewReader(validRequestBody()))
	postReq.Header.Set("Content-Type", "application/json")
	postW := httptest.NewRecorder()
	router.ServeHTTP(postW, postReq)

	if postW.Code != http.StatusOK {
		t.Fatalf("POST failed: %d %s", postW.Code, postW.Body.String())
	}

	var created models.ExplainResponse
	json.NewDecoder(postW.Body).Decode(&created)

	// GET the same ID.
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/explain/"+created.ID, nil)
	getW := httptest.NewRecorder()
	router.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Fatalf("GET failed: %d %s", getW.Code, getW.Body.String())
	}

	var fetched models.ExplainResponse
	json.NewDecoder(getW.Body).Decode(&fetched)

	if fetched.ID != created.ID {
		t.Errorf("ID mismatch: got %s, want %s", fetched.ID, created.ID)
	}
	if fetched.Target != created.Target {
		t.Errorf("target mismatch: got %s, want %s", fetched.Target, created.Target)
	}
}

func TestGetExplain_Nonexistent(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/explain/nonexistent-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHealthEndpoint(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)

	if body["status"] != "healthy" {
		t.Errorf("status: got %s, want healthy", body["status"])
	}
}

func TestRequestIDHeader(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	reqID := w.Header().Get("X-Request-Id")
	if reqID == "" {
		t.Fatal("X-Request-Id header is missing")
	}
	// Validate it's a valid UUID.
	if _, err := uuid.Parse(reqID); err != nil {
		t.Errorf("X-Request-Id is not a valid UUID: %s", reqID)
	}
}

func TestProcessingTimeHeader(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	timing := w.Header().Get("X-Processing-Time-Ms")
	if timing == "" {
		t.Fatal("X-Processing-Time-Ms header is missing")
	}
}
