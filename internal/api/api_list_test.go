package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/models"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

func postExplainSimple(t *testing.T, router http.Handler, target string) models.ExplainResponse {
	t.Helper()
	req := models.ExplainRequest{
		Target: target,
		Value:  0.72,
		Components: []models.Component{
			{Name: "trend", Value: 0.8, Weight: 0.4, Confidence: 0.9},
		},
	}
	body, _ := json.Marshal(req)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/explain", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/v1/explain failed: %d %s", w.Code, w.Body.String())
	}

	var resp models.ExplainResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode POST response: %v", err)
	}
	return resp
}

func TestListExplain_ReturnsAll(t *testing.T) {
	router := newTestRouter()

	postExplainSimple(t, router, "score")
	postExplainSimple(t, router, "rank")
	postExplainSimple(t, router, "rating")

	r := httptest.NewRequest(http.MethodGet, "/api/v1/explain", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/v1/explain: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var result storage.ListResult
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("decode list response: %v", err)
	}

	if len(result.Items) != 3 {
		t.Errorf("items: expected 3, got %d", len(result.Items))
	}
	if result.Total != 3 {
		t.Errorf("total: expected 3, got %d", result.Total)
	}
}

func TestListExplain_LimitAndCursor(t *testing.T) {
	router := newTestRouter()

	postExplainSimple(t, router, "a")
	postExplainSimple(t, router, "b")
	postExplainSimple(t, router, "c")

	// Page 1: limit=1.
	r := httptest.NewRequest(http.MethodGet, "/api/v1/explain?limit=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("page1: expected 200, got %d", w.Code)
	}

	var page1 storage.ListResult
	json.NewDecoder(w.Body).Decode(&page1)

	if len(page1.Items) != 1 {
		t.Errorf("page1 items: expected 1, got %d", len(page1.Items))
	}
	if page1.NextCursor == "" {
		t.Fatal("page1 should have next_cursor")
	}

	// Page 2: use cursor.
	r = httptest.NewRequest(http.MethodGet, "/api/v1/explain?limit=1&cursor="+page1.NextCursor, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("page2: expected 200, got %d", w.Code)
	}

	var page2 storage.ListResult
	json.NewDecoder(w.Body).Decode(&page2)

	if len(page2.Items) != 1 {
		t.Errorf("page2 items: expected 1, got %d", len(page2.Items))
	}
	if page1.Items[0].ID == page2.Items[0].ID {
		t.Error("page2 returned same item as page1")
	}
}

func TestListExplain_TargetFilter(t *testing.T) {
	router := newTestRouter()

	postExplainSimple(t, router, "user_score")
	postExplainSimple(t, router, "rank_score")
	postExplainSimple(t, router, "category")

	r := httptest.NewRequest(http.MethodGet, "/api/v1/explain?target=score", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result storage.ListResult
	json.NewDecoder(w.Body).Decode(&result)

	if len(result.Items) != 2 {
		t.Errorf("items: expected 2 (score matches), got %d", len(result.Items))
	}
}

func TestListExplain_EmptyResult(t *testing.T) {
	router := newTestRouter()

	r := httptest.NewRequest(http.MethodGet, "/api/v1/explain", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result storage.ListResult
	json.NewDecoder(w.Body).Decode(&result)

	if len(result.Items) != 0 {
		t.Errorf("items: expected 0, got %d", len(result.Items))
	}
	if result.Total != 0 {
		t.Errorf("total: expected 0, got %d", result.Total)
	}
}

func TestStatsEndpoint(t *testing.T) {
	router := newTestRouter()

	postExplainSimple(t, router, "score")
	postExplainSimple(t, router, "rank")

	r := httptest.NewRequest(http.MethodGet, "/api/v1/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var body map[string]any
	json.NewDecoder(w.Body).Decode(&body)

	if body["status"] != "ok" {
		t.Errorf("status: expected ok, got %v", body["status"])
	}
	if body["total_explanations"] != float64(2) {
		t.Errorf("total_explanations: expected 2, got %v", body["total_explanations"])
	}
}

func TestListExplain_InvalidLimit(t *testing.T) {
	router := newTestRouter()

	r := httptest.NewRequest(http.MethodGet, "/api/v1/explain?limit=abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
