package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/engine"
	"github.com/blackms/ExplainableEngine/internal/models"
)

// postExplainAndGetID creates an explanation via POST and returns the ID.
func postExplainAndGetID(t *testing.T, router http.Handler) string {
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

func TestNarrative_BasicEN(t *testing.T) {
	router := newTestRouter()
	id := postExplainAndGetID(t, router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/explain/"+id+"/narrative?level=basic&lang=en", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var result engine.NarrativeResult
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if result.Level != "basic" {
		t.Errorf("level: got %q, want %q", result.Level, "basic")
	}
	if result.Language != "en" {
		t.Errorf("language: got %q, want %q", result.Language, "en")
	}
	if result.ExplanationID != id {
		t.Errorf("explanation_id: got %q, want %q", result.ExplanationID, id)
	}
	if !strings.Contains(result.Narrative, "score") {
		t.Error("narrative should contain target name")
	}
}

func TestNarrative_AdvancedIT(t *testing.T) {
	router := newTestRouter()
	id := postExplainAndGetID(t, router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/explain/"+id+"/narrative?level=advanced&lang=it", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var result engine.NarrativeResult
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if result.Level != "advanced" {
		t.Errorf("level: got %q, want %q", result.Level, "advanced")
	}
	if result.Language != "it" {
		t.Errorf("language: got %q, want %q", result.Language, "it")
	}
	if !strings.Contains(result.Narrative, "Il punteggio") {
		t.Error("Italian narrative should contain 'Il punteggio'")
	}
}

func TestNarrative_InvalidLevel(t *testing.T) {
	router := newTestRouter()
	id := postExplainAndGetID(t, router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/explain/"+id+"/narrative?level=expert&lang=en", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNarrative_InvalidLang(t *testing.T) {
	router := newTestRouter()
	id := postExplainAndGetID(t, router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/explain/"+id+"/narrative?level=basic&lang=fr", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNarrative_NonexistentID(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/explain/nonexistent-id/narrative", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestNarrative_DefaultParams(t *testing.T) {
	router := newTestRouter()
	id := postExplainAndGetID(t, router)

	// No level or lang params — should default to basic/en
	req := httptest.NewRequest(http.MethodGet, "/api/v1/explain/"+id+"/narrative", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var result engine.NarrativeResult
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if result.Level != "basic" {
		t.Errorf("default level: got %q, want %q", result.Level, "basic")
	}
	if result.Language != "en" {
		t.Errorf("default language: got %q, want %q", result.Language, "en")
	}
}
