package engine

import (
	"testing"

	"github.com/blackms/ExplainableEngine/internal/models"
)

func TestDependency_Leaf(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "leaf", Label: "Leaf"})

	tree, err := ResolveDependencies(dag, "leaf")
	if err != nil {
		t.Fatal(err)
	}
	if tree.Root.NodeID != "leaf" {
		t.Errorf("root should be leaf, got %s", tree.Root.NodeID)
	}
	if tree.Depth != 0 {
		t.Errorf("depth should be 0, got %d", tree.Depth)
	}
	if tree.TotalNodes != 1 {
		t.Errorf("total nodes should be 1, got %d", tree.TotalNodes)
	}
	if len(tree.Root.Children) != 0 {
		t.Errorf("leaf should have no children, got %d", len(tree.Root.Children))
	}
}

func TestDependency_Linear(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "root", Label: "Root"})
	dag.AddNode(models.Node{ID: "mid", Label: "Mid"})
	dag.AddNode(models.Node{ID: "leaf", Label: "Leaf"})
	dag.AddEdge(models.Edge{Source: "mid", Target: "root", Weight: 1.0})
	dag.AddEdge(models.Edge{Source: "leaf", Target: "mid", Weight: 1.0})

	tree, err := ResolveDependencies(dag, "root")
	if err != nil {
		t.Fatal(err)
	}
	if tree.Depth != 2 {
		t.Errorf("depth should be 2, got %d", tree.Depth)
	}
	if tree.TotalNodes != 3 {
		t.Errorf("total nodes should be 3, got %d", tree.TotalNodes)
	}
	if len(tree.Root.Children) != 1 {
		t.Fatalf("root should have 1 child, got %d", len(tree.Root.Children))
	}
	if tree.Root.Children[0].NodeID != "mid" {
		t.Errorf("expected mid, got %s", tree.Root.Children[0].NodeID)
	}
	if len(tree.Root.Children[0].Children) != 1 {
		t.Fatalf("mid should have 1 child, got %d", len(tree.Root.Children[0].Children))
	}
}

func TestDependency_Diamond(t *testing.T) {
	// root -> a, root -> b, a -> leaf, b -> leaf
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "root", Label: "Root"})
	dag.AddNode(models.Node{ID: "a", Label: "A"})
	dag.AddNode(models.Node{ID: "b", Label: "B"})
	dag.AddNode(models.Node{ID: "leaf", Label: "Leaf"})
	dag.AddEdge(models.Edge{Source: "a", Target: "root", Weight: 0.5})
	dag.AddEdge(models.Edge{Source: "b", Target: "root", Weight: 0.5})
	dag.AddEdge(models.Edge{Source: "leaf", Target: "a", Weight: 1.0})
	dag.AddEdge(models.Edge{Source: "leaf", Target: "b", Weight: 1.0})

	tree, err := ResolveDependencies(dag, "root")
	if err != nil {
		t.Fatal(err)
	}
	// Root has children a and b. Both a and b have child leaf.
	// So total nodes = root + a + b + leaf(under a) + leaf(under b) = 5
	if tree.TotalNodes != 5 {
		t.Errorf("total nodes should be 5, got %d", tree.TotalNodes)
	}
	if tree.Depth != 2 {
		t.Errorf("depth should be 2, got %d", tree.Depth)
	}
	if len(tree.Root.Children) != 2 {
		t.Fatalf("root should have 2 children, got %d", len(tree.Root.Children))
	}
}

func TestDFS_Traversal(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "root", Label: "Root"})
	dag.AddNode(models.Node{ID: "a", Label: "A"})
	dag.AddNode(models.Node{ID: "b", Label: "B"})
	dag.AddEdge(models.Edge{Source: "a", Target: "root", Weight: 0.5})
	dag.AddEdge(models.Edge{Source: "b", Target: "root", Weight: 0.5})

	tree, _ := ResolveDependencies(dag, "root")
	visited := TraverseDFS(tree)

	if len(visited) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(visited))
	}
	// DFS: root -> a -> b (a and b are sorted children of root)
	if visited[0].NodeID != "root" {
		t.Errorf("first should be root, got %s", visited[0].NodeID)
	}
	if visited[1].NodeID != "a" {
		t.Errorf("second should be a, got %s", visited[1].NodeID)
	}
	if visited[2].NodeID != "b" {
		t.Errorf("third should be b, got %s", visited[2].NodeID)
	}
}

func TestBFS_Traversal(t *testing.T) {
	dag := NewDAG()
	dag.AddNode(models.Node{ID: "root", Label: "Root"})
	dag.AddNode(models.Node{ID: "a", Label: "A"})
	dag.AddNode(models.Node{ID: "b", Label: "B"})
	dag.AddNode(models.Node{ID: "leaf", Label: "Leaf"})
	dag.AddEdge(models.Edge{Source: "a", Target: "root", Weight: 0.5})
	dag.AddEdge(models.Edge{Source: "b", Target: "root", Weight: 0.5})
	dag.AddEdge(models.Edge{Source: "leaf", Target: "a", Weight: 1.0})

	tree, _ := ResolveDependencies(dag, "root")
	visited := TraverseBFS(tree)

	if len(visited) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(visited))
	}
	// BFS: root -> a, b -> leaf
	if visited[0].NodeID != "root" {
		t.Errorf("first should be root, got %s", visited[0].NodeID)
	}
	// a and b at depth 1
	if visited[1].NodeID != "a" || visited[2].NodeID != "b" {
		t.Errorf("expected a, b at depth 1, got %s, %s", visited[1].NodeID, visited[2].NodeID)
	}
	if visited[3].NodeID != "leaf" {
		t.Errorf("last should be leaf, got %s", visited[3].NodeID)
	}
}

func TestDependency_NotFound(t *testing.T) {
	dag := NewDAG()
	_, err := ResolveDependencies(dag, "missing")
	if err == nil {
		t.Error("expected error for missing node")
	}
}
