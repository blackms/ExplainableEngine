package engine

import (
	"testing"

	"github.com/blackms/ExplainableEngine/internal/models"
)

func buildMissingTestDAG(t *testing.T, components []models.Component) (*DAG, string) {
	t.Helper()
	dag := NewDAG()
	rootID := "root"
	if err := dag.AddNode(models.Node{
		ID:       rootID,
		Label:    "root",
		Value:    100.0,
		NodeType: models.NodeTypeOutput,
	}); err != nil {
		t.Fatalf("adding root: %v", err)
	}
	for _, c := range components {
		nt := models.NodeTypeInput
		if c.Missing {
			nt = models.NodeTypeMissing
		}
		if err := dag.AddNode(models.Node{
			ID:       c.Name,
			Label:    c.Name,
			Value:    c.Value,
			NodeType: nt,
		}); err != nil {
			t.Fatalf("adding node %s: %v", c.Name, err)
		}
		if err := dag.AddEdge(models.Edge{
			Source: c.Name,
			Target: rootID,
			Weight: c.Weight,
		}); err != nil {
			t.Fatalf("adding edge %s->root: %v", c.Name, err)
		}
	}
	return dag, rootID
}

func TestAnalyzeMissingData_NoMissing(t *testing.T) {
	dag, rootID := buildMissingTestDAG(t, []models.Component{
		{Name: "a", Value: 10, Weight: 0.5},
		{Name: "b", Value: 20, Weight: 0.5},
	})

	result, err := AnalyzeMissingData(dag, rootID, 0.2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.MissingNodes) != 0 {
		t.Errorf("expected 0 missing nodes, got %d", len(result.MissingNodes))
	}
	if result.TotalImpact != 0.0 {
		t.Errorf("expected impact 0.0, got %f", result.TotalImpact)
	}
	if len(result.Warnings) != 0 {
		t.Errorf("expected no warnings, got %v", result.Warnings)
	}
}

func TestAnalyzeMissingData_OneMissing(t *testing.T) {
	dag, rootID := buildMissingTestDAG(t, []models.Component{
		{Name: "a", Value: 10, Weight: 0.3},
		{Name: "b", Value: 20, Weight: 0.7, Missing: true},
	})

	result, err := AnalyzeMissingData(dag, rootID, 0.2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.MissingNodes) != 1 {
		t.Fatalf("expected 1 missing node, got %d", len(result.MissingNodes))
	}
	if result.MissingNodes[0].NodeID != "b" {
		t.Errorf("expected missing node 'b', got %s", result.MissingNodes[0].NodeID)
	}

	// Impact = 0.7 / (0.3 + 0.7) = 0.7
	expectedImpact := 0.7
	if diff := result.TotalImpact - expectedImpact; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("expected impact %f, got %f", expectedImpact, result.TotalImpact)
	}

	// 0.7 > 0.2 threshold, so warning should be present.
	if len(result.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(result.Warnings))
	}
}

func TestAnalyzeMissingData_AllMissing(t *testing.T) {
	dag, rootID := buildMissingTestDAG(t, []models.Component{
		{Name: "a", Value: 10, Weight: 0.4, Missing: true},
		{Name: "b", Value: 20, Weight: 0.6, Missing: true},
	})

	result, err := AnalyzeMissingData(dag, rootID, 0.2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.MissingNodes) != 2 {
		t.Fatalf("expected 2 missing nodes, got %d", len(result.MissingNodes))
	}
	if result.TotalImpact != 1.0 {
		t.Errorf("expected impact 1.0, got %f", result.TotalImpact)
	}
	if len(result.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(result.Warnings))
	}
}

func TestAnalyzeMissingData_ThresholdWarningTriggered(t *testing.T) {
	dag, rootID := buildMissingTestDAG(t, []models.Component{
		{Name: "a", Value: 10, Weight: 0.5},
		{Name: "b", Value: 20, Weight: 0.5, Missing: true},
	})

	// Impact = 0.5, threshold = 0.3 => warning.
	result, err := AnalyzeMissingData(dag, rootID, 0.3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Warnings) == 0 {
		t.Error("expected a warning when impact exceeds threshold")
	}
}

func TestAnalyzeMissingData_ThresholdNotTriggered(t *testing.T) {
	dag, rootID := buildMissingTestDAG(t, []models.Component{
		{Name: "a", Value: 10, Weight: 0.9},
		{Name: "b", Value: 0, Weight: 0.1, Missing: true},
	})

	// Impact = 0.1, threshold = 0.2 => no warning.
	result, err := AnalyzeMissingData(dag, rootID, 0.2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Warnings) != 0 {
		t.Errorf("expected no warnings when impact is below threshold, got %v", result.Warnings)
	}
	if result.TotalImpact != 0.1 {
		t.Errorf("expected impact 0.1, got %f", result.TotalImpact)
	}
}

func TestAnalyzeMissingData_Recursive(t *testing.T) {
	// Build a DAG with nested missing node:
	// root <- parent <- child_missing
	dag := NewDAG()
	if err := dag.AddNode(models.Node{ID: "root", Label: "root", Value: 100, NodeType: models.NodeTypeOutput}); err != nil {
		t.Fatal(err)
	}
	if err := dag.AddNode(models.Node{ID: "parent", Label: "parent", Value: 50, NodeType: models.NodeTypeComputed}); err != nil {
		t.Fatal(err)
	}
	if err := dag.AddNode(models.Node{ID: "child_ok", Label: "child_ok", Value: 30, NodeType: models.NodeTypeInput}); err != nil {
		t.Fatal(err)
	}
	if err := dag.AddNode(models.Node{ID: "child_missing", Label: "child_missing", Value: 0, NodeType: models.NodeTypeMissing}); err != nil {
		t.Fatal(err)
	}

	// root <- parent
	if err := dag.AddEdge(models.Edge{Source: "parent", Target: "root", Weight: 1.0}); err != nil {
		t.Fatal(err)
	}
	// parent <- child_ok, child_missing
	if err := dag.AddEdge(models.Edge{Source: "child_ok", Target: "parent", Weight: 0.6}); err != nil {
		t.Fatal(err)
	}
	if err := dag.AddEdge(models.Edge{Source: "child_missing", Target: "parent", Weight: 0.4}); err != nil {
		t.Fatal(err)
	}

	result, err := AnalyzeMissingData(dag, "root", 0.1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.MissingNodes) != 1 {
		t.Fatalf("expected 1 missing node in recursive walk, got %d", len(result.MissingNodes))
	}
	if result.MissingNodes[0].NodeID != "child_missing" {
		t.Errorf("expected missing node 'child_missing', got %s", result.MissingNodes[0].NodeID)
	}

	// Global impact: missing weight = 0.4, total weight = 1.0 + 0.6 + 0.4 = 2.0
	// So impact = 0.4 / 2.0 = 0.2
	expectedImpact := 0.2
	if diff := result.TotalImpact - expectedImpact; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("expected impact %f, got %f", expectedImpact, result.TotalImpact)
	}

	// 0.2 > 0.1 threshold => warning
	if len(result.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(result.Warnings))
	}
}

func TestAnalyzeMissingData_SortedByNodeID(t *testing.T) {
	dag, rootID := buildMissingTestDAG(t, []models.Component{
		{Name: "z_node", Value: 0, Weight: 0.3, Missing: true},
		{Name: "a_node", Value: 0, Weight: 0.3, Missing: true},
		{Name: "m_node", Value: 10, Weight: 0.4},
	})

	result, err := AnalyzeMissingData(dag, rootID, 0.2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.MissingNodes) != 2 {
		t.Fatalf("expected 2 missing nodes, got %d", len(result.MissingNodes))
	}
	if result.MissingNodes[0].NodeID != "a_node" {
		t.Errorf("expected first missing node 'a_node', got %s", result.MissingNodes[0].NodeID)
	}
	if result.MissingNodes[1].NodeID != "z_node" {
		t.Errorf("expected second missing node 'z_node', got %s", result.MissingNodes[1].NodeID)
	}
}
