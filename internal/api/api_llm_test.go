package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/api"
	"github.com/blackms/ExplainableEngine/internal/llm"
	"github.com/blackms/ExplainableEngine/internal/models"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

// newTestRouterWithLLM creates a test router with the fallback LLM service.
func newTestRouterWithLLM() http.Handler {
	store := storage.NewInMemoryStore()
	orch := &mockOrchestrator{}
	return api.NewRouter(store, orch, api.WithLLMService(llm.NewFallbackService()))
}

// postExplainAndGetIDForLLM creates an explanation via POST and returns the ID.
func postExplainAndGetIDForLLM(t *testing.T, router http.Handler) string {
	t.Helper()
	body := validRequestBody()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/v1/explain failed: %d %s", w.Code, w.Body.String())
	}

	var resp models.ExplainResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode POST response: %v", err)
	}
	return resp.ID
}

func TestLLMNarrative_FallbackReturns200(t *testing.T) {
	router := newTestRouterWithLLM()
	id := postExplainAndGetIDForLLM(t, router)

	body, _ := json.Marshal(map[string]string{"level": "basic", "lang": "en"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain/"+id+"/narrative/llm", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if result["narrative"] == nil || result["narrative"] == "" {
		t.Error("narrative should not be empty")
	}
	if result["source"] != "template" {
		t.Errorf("source should be 'template' for fallback, got %v", result["source"])
	}
}

func TestLLMNarrative_DefaultParams(t *testing.T) {
	router := newTestRouterWithLLM()
	id := postExplainAndGetIDForLLM(t, router)

	// Empty body should use defaults (basic/en).
	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain/"+id+"/narrative/llm", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestLLMNarrative_InvalidLevel(t *testing.T) {
	router := newTestRouterWithLLM()
	id := postExplainAndGetIDForLLM(t, router)

	body, _ := json.Marshal(map[string]string{"level": "expert"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain/"+id+"/narrative/llm", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestLLMNarrative_NonexistentID(t *testing.T) {
	router := newTestRouterWithLLM()

	body, _ := json.Marshal(map[string]string{"level": "basic"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain/nonexistent-id/narrative/llm", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAskQuestion_FallbackReturns503(t *testing.T) {
	router := newTestRouterWithLLM()
	id := postExplainAndGetIDForLLM(t, router)

	body, _ := json.Marshal(map[string]string{"question": "Why is the score low?"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain/"+id+"/ask", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAskQuestion_EmptyQuestion(t *testing.T) {
	router := newTestRouterWithLLM()
	id := postExplainAndGetIDForLLM(t, router)

	body, _ := json.Marshal(map[string]string{"question": ""})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain/"+id+"/ask", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAskQuestion_NonexistentID(t *testing.T) {
	router := newTestRouterWithLLM()

	body, _ := json.Marshal(map[string]string{"question": "Why?"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain/nonexistent-id/ask", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGenerateSummary_FallbackReturns503(t *testing.T) {
	router := newTestRouterWithLLM()
	id := postExplainAndGetIDForLLM(t, router)

	body, _ := json.Marshal(map[string]string{"audience": "board", "lang": "en"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain/"+id+"/summary", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGenerateSummary_InvalidAudience(t *testing.T) {
	router := newTestRouterWithLLM()
	id := postExplainAndGetIDForLLM(t, router)

	body, _ := json.Marshal(map[string]string{"audience": "aliens"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain/"+id+"/summary", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGenerateSummary_NonexistentID(t *testing.T) {
	router := newTestRouterWithLLM()

	body, _ := json.Marshal(map[string]string{"audience": "board"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/explain/nonexistent-id/summary", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

// TestExistingRoutesStillWork verifies that adding LLM routes does not break
// the original router (backward compatibility with the variadic option).
func TestExistingRoutesStillWork(t *testing.T) {
	// Create router WITHOUT the LLM option — should still work.
	store := storage.NewInMemoryStore()
	orch := &mockOrchestrator{}
	router := api.NewRouter(store, orch)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
