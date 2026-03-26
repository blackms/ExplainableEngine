package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/api"
	"github.com/blackms/ExplainableEngine/internal/engine"
	"github.com/blackms/ExplainableEngine/internal/models"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

// newGraphTestServer creates a test server with a real orchestrator and in-memory store.
func newGraphTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	store := storage.NewInMemoryStore()
	orch := engine.NewOrchestrator()
	router := api.NewRouter(store, orch)
	return httptest.NewServer(router)
}

// createExplanation posts an explanation request that includes a graph and returns the ID.
func createExplanation(t *testing.T, ts *httptest.Server) string {
	t.Helper()
	req := models.ExplainRequest{
		Target: "market_regime_score",
		Value:  0.72,
		Components: []models.Component{
			{Name: "trend_strength", Value: 0.80, Weight: 0.40, Confidence: 0.90},
			{Name: "volatility", Value: 0.55, Weight: 0.60, Confidence: 0.75},
		},
		Options: &models.ExplainOptions{
			IncludeGraph: true,
		},
	}
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
		data, _ := io.ReadAll(resp.Body)
		t.Fatalf("POST expected 200, got %d: %s", resp.StatusCode, string(data))
	}

	var result models.ExplainResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result.ID == "" {
		t.Fatal("explanation ID should not be empty")
	}
	return result.ID
}

func TestGraphExport_JSON(t *testing.T) {
	ts := newGraphTestServer(t)
	defer ts.Close()

	id := createExplanation(t, ts)

	resp, err := http.Get(ts.URL + "/api/v1/explain/" + id + "/graph")
	if err != nil {
		t.Fatalf("GET graph: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(body))
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("content type: got %s, want application/json", ct)
	}

	var graph models.GraphResponse
	if err := json.NewDecoder(resp.Body).Decode(&graph); err != nil {
		t.Fatalf("decode graph JSON: %v", err)
	}
	if len(graph.Nodes) == 0 {
		t.Error("graph should have nodes")
	}
	if len(graph.Edges) == 0 {
		t.Error("graph should have edges")
	}
}

func TestGraphExport_DOT(t *testing.T) {
	ts := newGraphTestServer(t)
	defer ts.Close()

	id := createExplanation(t, ts)

	resp, err := http.Get(ts.URL + "/api/v1/explain/" + id + "/graph?format=dot")
	if err != nil {
		t.Fatalf("GET graph dot: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(body))
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "text/vnd.graphviz" {
		t.Errorf("content type: got %s, want text/vnd.graphviz", ct)
	}

	body, _ := io.ReadAll(resp.Body)
	content := string(body)
	if !strings.Contains(content, "digraph") {
		t.Error("DOT output should contain 'digraph'")
	}
	if !strings.Contains(content, "->") {
		t.Error("DOT output should contain edges with '->'")
	}
}

func TestGraphExport_Mermaid(t *testing.T) {
	ts := newGraphTestServer(t)
	defer ts.Close()

	id := createExplanation(t, ts)

	resp, err := http.Get(ts.URL + "/api/v1/explain/" + id + "/graph?format=mermaid")
	if err != nil {
		t.Fatalf("GET graph mermaid: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200, got %d: %s", resp.StatusCode, string(body))
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "text/plain" {
		t.Errorf("content type: got %s, want text/plain", ct)
	}

	body, _ := io.ReadAll(resp.Body)
	content := string(body)
	if !strings.Contains(content, "graph LR") {
		t.Error("Mermaid output should contain 'graph LR'")
	}
	if !strings.Contains(content, "-->|") {
		t.Error("Mermaid output should contain edges with '-->|'")
	}
}

func TestGraphExport_UnsupportedFormat(t *testing.T) {
	ts := newGraphTestServer(t)
	defer ts.Close()

	id := createExplanation(t, ts)

	resp, err := http.Get(ts.URL + "/api/v1/explain/" + id + "/graph?format=xml")
	if err != nil {
		t.Fatalf("GET graph xml: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for unsupported format, got %d", resp.StatusCode)
	}
}

func TestGraphExport_NotFound(t *testing.T) {
	ts := newGraphTestServer(t)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/v1/explain/nonexistent/graph")
	if err != nil {
		t.Fatalf("GET graph nonexistent: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 for nonexistent explanation, got %d", resp.StatusCode)
	}
}
