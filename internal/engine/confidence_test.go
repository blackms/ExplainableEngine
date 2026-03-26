package engine

import (
	"math"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/models"
)

func TestConfidence_Linear(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "root", Label: "Root", NodeType: models.NodeTypeOutput})
	dag.AddNode(models.Node{ID: "a", Label: "A", Confidence: 0.8, NodeType: models.NodeTypeInput})
	dag.AddEdge(models.Edge{Source: "a", Target: "root", Weight: 1.0})

	result, err := PropagateConfidence(dag, "root")
	if err != nil {
		t.Fatal(err)
	}
	// Root confidence = (1.0 * 0.8) / 1.0 = 0.8
	assertConfFloat(t, "overall", result.OverallConfidence, 0.8)
	assertConfFloat(t, "a", result.NodeConfidences["a"], 0.8)
	assertConfFloat(t, "root", result.NodeConfidences["root"], 0.8)
}

func TestConfidence_Diamond(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "root", Label: "Root", NodeType: models.NodeTypeOutput})
	dag.AddNode(models.Node{ID: "a", Label: "A", Confidence: 0.9, NodeType: models.NodeTypeInput})
	dag.AddNode(models.Node{ID: "b", Label: "B", Confidence: 0.7, NodeType: models.NodeTypeInput})
	dag.AddEdge(models.Edge{Source: "a", Target: "root", Weight: 0.6})
	dag.AddEdge(models.Edge{Source: "b", Target: "root", Weight: 0.4})

	result, err := PropagateConfidence(dag, "root")
	if err != nil {
		t.Fatal(err)
	}
	// Root = (0.6*0.9 + 0.4*0.7) / (0.6 + 0.4) = (0.54 + 0.28) / 1.0 = 0.82
	assertConfFloat(t, "overall", result.OverallConfidence, 0.82)
}

func TestConfidence_AllOne(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "root", Label: "Root", NodeType: models.NodeTypeOutput})
	dag.AddNode(models.Node{ID: "a", Label: "A", Confidence: 1.0, NodeType: models.NodeTypeInput})
	dag.AddNode(models.Node{ID: "b", Label: "B", Confidence: 1.0, NodeType: models.NodeTypeInput})
	dag.AddEdge(models.Edge{Source: "a", Target: "root", Weight: 0.5})
	dag.AddEdge(models.Edge{Source: "b", Target: "root", Weight: 0.5})

	result, err := PropagateConfidence(dag, "root")
	if err != nil {
		t.Fatal(err)
	}
	assertConfFloat(t, "overall", result.OverallConfidence, 1.0)
}

func TestConfidence_AllZero(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "root", Label: "Root", NodeType: models.NodeTypeOutput})
	dag.AddNode(models.Node{ID: "a", Label: "A", Confidence: 0.0, NodeType: models.NodeTypeInput})
	dag.AddNode(models.Node{ID: "b", Label: "B", Confidence: 0.0, NodeType: models.NodeTypeInput})
	dag.AddEdge(models.Edge{Source: "a", Target: "root", Weight: 0.5})
	dag.AddEdge(models.Edge{Source: "b", Target: "root", Weight: 0.5})

	result, err := PropagateConfidence(dag, "root")
	if err != nil {
		t.Fatal(err)
	}
	assertConfFloat(t, "overall", result.OverallConfidence, 0.0)
}

func TestConfidence_MixedWeights(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "root", Label: "Root", NodeType: models.NodeTypeOutput})
	dag.AddNode(models.Node{ID: "a", Label: "A", Confidence: 1.0, NodeType: models.NodeTypeInput})
	dag.AddNode(models.Node{ID: "b", Label: "B", Confidence: 0.5, NodeType: models.NodeTypeInput})
	dag.AddEdge(models.Edge{Source: "a", Target: "root", Weight: 0.8})
	dag.AddEdge(models.Edge{Source: "b", Target: "root", Weight: 0.2})

	result, err := PropagateConfidence(dag, "root")
	if err != nil {
		t.Fatal(err)
	}
	// (0.8*1.0 + 0.2*0.5) / (0.8+0.2) = 0.9
	assertConfFloat(t, "overall", result.OverallConfidence, 0.9)
}

func TestConfidence_ZeroWeights(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "root", Label: "Root", NodeType: models.NodeTypeOutput})
	dag.AddNode(models.Node{ID: "a", Label: "A", Confidence: 0.9, NodeType: models.NodeTypeInput})
	dag.AddEdge(models.Edge{Source: "a", Target: "root", Weight: 0.0})

	result, err := PropagateConfidence(dag, "root")
	if err != nil {
		t.Fatal(err)
	}
	assertConfFloat(t, "overall", result.OverallConfidence, 0.0)
}

func TestConfidence_PropagationPath(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "root", Label: "Root", NodeType: models.NodeTypeOutput})
	dag.AddNode(models.Node{ID: "a", Label: "A", Confidence: 0.8, NodeType: models.NodeTypeInput})
	dag.AddEdge(models.Edge{Source: "a", Target: "root", Weight: 1.0})

	result, err := PropagateConfidence(dag, "root")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.PropagationPath) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(result.PropagationPath))
	}
}

func TestConfidence_NotFound(t *testing.T) {
	dag := NewDAG()
	_, err := PropagateConfidence(dag, "missing")
	if err == nil {
		t.Error("expected error for missing node")
	}
}

func assertConfFloat(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s: got %f, want %f", name, got, want)
	}
}
