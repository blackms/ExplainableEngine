package engine

import (
	"errors"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/models"
)

func TestNewDAG_Empty(t *testing.T) {
	dag := NewDAG()
	if len(dag.Nodes()) != 0 {
		t.Errorf("new DAG should have 0 nodes, got %d", len(dag.Nodes()))
	}
	sorted, err := dag.TopologicalSort()
	if err != nil {
		t.Fatalf("topological sort on empty DAG: %v", err)
	}
	if len(sorted) != 0 {
		t.Errorf("topological sort should return empty slice, got %v", sorted)
	}
	if dag.HasCycle() {
		t.Error("empty DAG should not have a cycle")
	}
}

func TestAddNode_Duplicate(t *testing.T) {
	dag := NewDAG()
	node := models.Node{ID: "a", Label: "A", Value: 1.0}
	if err := dag.AddNode(node); err != nil {
		t.Fatalf("first add: %v", err)
	}
	if err := dag.AddNode(node); err == nil {
		t.Error("expected error for duplicate node")
	}
}

func TestAddEdge_MissingNodes(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "a", Label: "A"})

	// Target missing.
	err := dag.AddEdge(models.Edge{Source: "a", Target: "b", Weight: 0.5})
	var nfe *models.NodeNotFoundError
	if !errors.As(err, &nfe) {
		t.Errorf("expected NodeNotFoundError, got %v", err)
	}

	// Source missing.
	err = dag.AddEdge(models.Edge{Source: "c", Target: "a", Weight: 0.5})
	if !errors.As(err, &nfe) {
		t.Errorf("expected NodeNotFoundError, got %v", err)
	}
}

func TestGetNode_NotFound(t *testing.T) {
	dag := NewDAG()
	_, err := dag.GetNode("missing")
	var nfe *models.NodeNotFoundError
	if !errors.As(err, &nfe) {
		t.Errorf("expected NodeNotFoundError, got %v", err)
	}
}

func TestGetPredecessors(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "root", Label: "Root"})
	dag.AddNode(models.Node{ID: "b", Label: "B"})
	dag.AddNode(models.Node{ID: "a", Label: "A"})
	dag.AddEdge(models.Edge{Source: "b", Target: "root", Weight: 0.5})
	dag.AddEdge(models.Edge{Source: "a", Target: "root", Weight: 0.5})

	preds, err := dag.GetPredecessors("root")
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{"a", "b"}
	if len(preds) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, preds)
	}
	for i, p := range preds {
		if p != expected[i] {
			t.Errorf("predecessor[%d]: expected %s, got %s", i, expected[i], p)
		}
	}
}

func TestGetPredecessors_NotFound(t *testing.T) {
	dag := NewDAG()
	_, err := dag.GetPredecessors("missing")
	var nfe *models.NodeNotFoundError
	if !errors.As(err, &nfe) {
		t.Errorf("expected NodeNotFoundError, got %v", err)
	}
}

func TestGetAncestors(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "root", Label: "Root"})
	dag.AddNode(models.Node{ID: "mid", Label: "Mid"})
	dag.AddNode(models.Node{ID: "leaf", Label: "Leaf"})
	dag.AddEdge(models.Edge{Source: "mid", Target: "root", Weight: 1.0})
	dag.AddEdge(models.Edge{Source: "leaf", Target: "mid", Weight: 1.0})

	anc, err := dag.GetAncestors("root")
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{"leaf", "mid"}
	if len(anc) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, anc)
	}
	for i := range anc {
		if anc[i] != expected[i] {
			t.Errorf("ancestor[%d]: expected %s, got %s", i, expected[i], anc[i])
		}
	}
}

func TestGetDescendants(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "root", Label: "Root"})
	dag.AddNode(models.Node{ID: "mid", Label: "Mid"})
	dag.AddNode(models.Node{ID: "leaf", Label: "Leaf"})
	dag.AddEdge(models.Edge{Source: "leaf", Target: "mid", Weight: 1.0})
	dag.AddEdge(models.Edge{Source: "mid", Target: "root", Weight: 1.0})

	desc, err := dag.GetDescendants("leaf")
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{"mid", "root"}
	if len(desc) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, desc)
	}
	for i := range desc {
		if desc[i] != expected[i] {
			t.Errorf("descendant[%d]: expected %s, got %s", i, expected[i], desc[i])
		}
	}
}

func TestTopologicalSort_Linear(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "a"})
	dag.AddNode(models.Node{ID: "b"})
	dag.AddNode(models.Node{ID: "c"})
	dag.AddEdge(models.Edge{Source: "a", Target: "b"})
	dag.AddEdge(models.Edge{Source: "b", Target: "c"})

	sorted, err := dag.TopologicalSort()
	if err != nil {
		t.Fatal(err)
	}
	if len(sorted) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(sorted))
	}
	// a must come before b, b before c.
	indexOf := func(s string) int {
		for i, v := range sorted {
			if v == s {
				return i
			}
		}
		return -1
	}
	if indexOf("a") >= indexOf("b") || indexOf("b") >= indexOf("c") {
		t.Errorf("bad topological order: %v", sorted)
	}
}

func TestTopologicalSort_Cycle(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "a"})
	dag.AddNode(models.Node{ID: "b"})
	dag.AddEdge(models.Edge{Source: "a", Target: "b"})
	dag.AddEdge(models.Edge{Source: "b", Target: "a"})

	_, err := dag.TopologicalSort()
	var ce *models.CyclicGraphError
	if !errors.As(err, &ce) {
		t.Errorf("expected CyclicGraphError, got %v", err)
	}
}

func TestHasCycle(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "a"})
	dag.AddNode(models.Node{ID: "b"})
	dag.AddEdge(models.Edge{Source: "a", Target: "b"})
	dag.AddEdge(models.Edge{Source: "b", Target: "a"})

	if !dag.HasCycle() {
		t.Error("expected cycle to be detected")
	}
}

func TestBuildFromRequest_Simple(t *testing.T) {
	req := models.ExplainRequest{
		Target: "score",
		Value:  0.72,
		Components: []models.Component{
			{Name: "trend", Value: 0.8, Weight: 0.4, Confidence: 0.9},
			{Name: "momentum", Value: 0.6, Weight: 0.6, Confidence: 0.85},
		},
	}

	dag := NewDAG()
	dag, err := dag.BuildFromRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	nodes := dag.Nodes()
	if len(nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d: %v", len(nodes), nodes)
	}

	root, err := dag.GetNode("score")
	if err != nil {
		t.Fatal(err)
	}
	if root.NodeType != models.NodeTypeOutput {
		t.Errorf("root should be output, got %s", root.NodeType)
	}

	preds, _ := dag.GetPredecessors("score")
	if len(preds) != 2 {
		t.Errorf("expected 2 predecessors, got %d", len(preds))
	}
}

func TestBuildFromRequest_Nested(t *testing.T) {
	req := models.ExplainRequest{
		Target: "score",
		Value:  82.0,
		Components: []models.Component{
			{
				Name:   "skills",
				Value:  80.0,
				Weight: 0.5,
				Components: []models.Component{
					{Name: "python", Value: 95.0, Weight: 0.5, Confidence: 0.92},
					{Name: "go", Value: 70.0, Weight: 0.5, Confidence: 0.85},
				},
			},
			{Name: "experience", Value: 84.0, Weight: 0.5, Confidence: 0.88},
		},
	}

	dag := NewDAG()
	dag, err := dag.BuildFromRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	// score, skills, python, go, experience = 5 nodes
	nodes := dag.Nodes()
	if len(nodes) != 5 {
		t.Errorf("expected 5 nodes, got %d: %v", len(nodes), nodes)
	}

	// Skills should be computed type.
	skills, _ := dag.GetNode("skills")
	if skills.NodeType != models.NodeTypeComputed {
		t.Errorf("skills should be computed, got %s", skills.NodeType)
	}

	// Python and go are predecessors of skills.
	preds, _ := dag.GetPredecessors("skills")
	if len(preds) != 2 {
		t.Errorf("expected 2 predecessors of skills, got %d", len(preds))
	}
}
