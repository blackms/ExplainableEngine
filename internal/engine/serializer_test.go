package engine

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/models"
)

func sampleGraph() *models.GraphResponse {
	return &models.GraphResponse{
		Nodes: []models.GraphNodeResponse{
			{ID: "trend_strength", Label: "trend_strength", Value: 0.80, Confidence: 0.90, NodeType: "input"},
			{ID: "market_regime_score", Label: "market_regime_score", Value: 0.72, Confidence: 0.83, NodeType: "output"},
			{ID: "volatility", Label: "volatility", Value: 0.55, Confidence: 0.75, NodeType: "computed"},
		},
		Edges: []models.GraphEdgeResponse{
			{Source: "trend_strength", Target: "market_regime_score", Weight: 0.40, TransformationType: "weighted_sum"},
			{Source: "volatility", Target: "market_regime_score", Weight: 0.60, TransformationType: "weighted_sum"},
		},
	}
}

func emptyGraph() *models.GraphResponse {
	return &models.GraphResponse{
		Nodes: []models.GraphNodeResponse{},
		Edges: []models.GraphEdgeResponse{},
	}
}

func TestSerializeGraph_JSON_Valid(t *testing.T) {
	graph := sampleGraph()
	content, contentType, err := SerializeGraph(graph, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contentType != "application/json" {
		t.Errorf("content type: got %s, want application/json", contentType)
	}

	// Verify it's valid JSON.
	var parsed models.GraphResponse
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	// Verify it contains nodes and edges.
	if len(parsed.Nodes) != 3 {
		t.Errorf("nodes count: got %d, want 3", len(parsed.Nodes))
	}
	if len(parsed.Edges) != 2 {
		t.Errorf("edges count: got %d, want 2", len(parsed.Edges))
	}
}

func TestSerializeGraph_DOT_ContainsExpectedElements(t *testing.T) {
	graph := sampleGraph()
	content, contentType, err := SerializeGraph(graph, FormatDOT)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contentType != "text/vnd.graphviz" {
		t.Errorf("content type: got %s, want text/vnd.graphviz", contentType)
	}

	// Must contain "digraph".
	if !strings.Contains(content, "digraph") {
		t.Error("DOT output should contain 'digraph'")
	}

	// Must contain node IDs.
	for _, id := range []string{"trend_strength", "market_regime_score", "volatility"} {
		if !strings.Contains(content, id) {
			t.Errorf("DOT output should contain node ID %q", id)
		}
	}

	// Must contain edges with arrows.
	if !strings.Contains(content, "->") {
		t.Error("DOT output should contain '->' for edges")
	}

	// Must contain fillcolors by type.
	if !strings.Contains(content, "#90EE90") {
		t.Error("DOT output should contain green color for input nodes")
	}
	if !strings.Contains(content, "#FFB6C1") {
		t.Error("DOT output should contain pink color for output nodes")
	}
	if !strings.Contains(content, "#ADD8E6") {
		t.Error("DOT output should contain blue color for computed nodes")
	}
}

func TestSerializeGraph_Mermaid_ContainsExpectedElements(t *testing.T) {
	graph := sampleGraph()
	content, contentType, err := SerializeGraph(graph, FormatMermaid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contentType != "text/plain" {
		t.Errorf("content type: got %s, want text/plain", contentType)
	}

	// Must contain "graph LR".
	if !strings.Contains(content, "graph LR") {
		t.Error("Mermaid output should contain 'graph LR'")
	}

	// Must contain node definitions with brackets.
	for _, id := range []string{"trend_strength", "market_regime_score", "volatility"} {
		if !strings.Contains(content, id+"[") {
			t.Errorf("Mermaid output should contain node definition for %q", id)
		}
	}

	// Must contain edges with -->|w=...|.
	if !strings.Contains(content, "-->|") {
		t.Error("Mermaid output should contain '-->|' for edges")
	}

	// Must contain style lines.
	if !strings.Contains(content, "style trend_strength fill:#90EE90") {
		t.Error("Mermaid output should contain style for input node")
	}
	if !strings.Contains(content, "style market_regime_score fill:#FFB6C1") {
		t.Error("Mermaid output should contain style for output node")
	}
}

func TestSerializeGraph_UnsupportedFormat(t *testing.T) {
	graph := sampleGraph()
	_, _, err := SerializeGraph(graph, GraphFormat("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("error message should mention unsupported format, got: %v", err)
	}
}

func TestSerializeGraph_EmptyGraph(t *testing.T) {
	graph := emptyGraph()

	for _, format := range []GraphFormat{FormatJSON, FormatDOT, FormatMermaid} {
		content, _, err := SerializeGraph(graph, format)
		if err != nil {
			t.Fatalf("format %s: unexpected error: %v", format, err)
		}
		if content == "" {
			t.Errorf("format %s: output should not be empty even for empty graph", format)
		}
	}
}

func TestSerializeGraph_Determinism(t *testing.T) {
	// Create a graph with nodes in non-sorted order.
	graph := &models.GraphResponse{
		Nodes: []models.GraphNodeResponse{
			{ID: "z_node", Label: "z_node", Value: 1.0, Confidence: 0.5, NodeType: "input"},
			{ID: "a_node", Label: "a_node", Value: 2.0, Confidence: 0.8, NodeType: "output"},
			{ID: "m_node", Label: "m_node", Value: 1.5, Confidence: 0.7, NodeType: "computed"},
		},
		Edges: []models.GraphEdgeResponse{
			{Source: "z_node", Target: "a_node", Weight: 0.3, TransformationType: "weighted_sum"},
			{Source: "m_node", Target: "a_node", Weight: 0.7, TransformationType: "weighted_sum"},
		},
	}

	for _, format := range []GraphFormat{FormatJSON, FormatDOT, FormatMermaid} {
		first, _, err := SerializeGraph(graph, format)
		if err != nil {
			t.Fatalf("format %s first call: %v", format, err)
		}
		second, _, err := SerializeGraph(graph, format)
		if err != nil {
			t.Fatalf("format %s second call: %v", format, err)
		}
		if first != second {
			t.Errorf("format %s: output is not deterministic\nfirst:\n%s\nsecond:\n%s", format, first, second)
		}
	}
}

func TestSerializeGraph_MissingNodeColor(t *testing.T) {
	graph := &models.GraphResponse{
		Nodes: []models.GraphNodeResponse{
			{ID: "missing_data", Label: "missing_data", Value: 0.0, Confidence: 0.0, NodeType: "missing"},
		},
		Edges: []models.GraphEdgeResponse{},
	}

	// DOT format should use grey for missing nodes.
	content, _, err := SerializeGraph(graph, FormatDOT)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(content, "#D3D3D3") {
		t.Error("DOT output should contain grey color for missing nodes")
	}

	// Mermaid format should use grey for missing nodes.
	content, _, err = SerializeGraph(graph, FormatMermaid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(content, "#D3D3D3") {
		t.Error("Mermaid output should contain grey color for missing nodes")
	}
}
