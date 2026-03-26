package engine

import (
	"testing"

	"github.com/blackms/ExplainableEngine/internal/models"
)

func makeTestRequest() models.ExplainRequest {
	return models.ExplainRequest{
		Target: "score",
		Value:  0.72,
		Components: []models.Component{
			{Name: "trend", Value: 0.8, Weight: 0.4, Confidence: 0.9},
			{Name: "momentum", Value: 0.6, Weight: 0.3, Confidence: 0.85},
			{Name: "volatility", Value: 0.5, Weight: 0.3, Confidence: 0.75},
		},
		Options: &models.ExplainOptions{
			IncludeGraph:   true,
			IncludeDrivers: true,
			MaxDrivers:     5,
			MaxDepth:       10,
		},
	}
}

func TestOrchestrator_FullPipeline(t *testing.T) {
	orch := NewOrchestrator()
	req := makeTestRequest()

	resp, err := orch.Explain(req)
	if err != nil {
		t.Fatalf("Explain failed: %v", err)
	}

	// Check basic fields.
	if resp.ID == "" {
		t.Error("ID should not be empty")
	}
	if resp.Target != "score" {
		t.Errorf("target: got %s, want score", resp.Target)
	}
	if resp.FinalValue != 0.72 {
		t.Errorf("final value: got %f, want 0.72", resp.FinalValue)
	}

	// Breakdown should have 3 items.
	if len(resp.Breakdown) != 3 {
		t.Errorf("breakdown items: got %d, want 3", len(resp.Breakdown))
	}

	// Top drivers should be present.
	if len(resp.TopDrivers) == 0 {
		t.Error("top drivers should not be empty")
	}
	// Drivers should be ranked.
	for i, d := range resp.TopDrivers {
		if d.Rank != i+1 {
			t.Errorf("driver rank: got %d, want %d", d.Rank, i+1)
		}
	}

	// Graph should be included.
	if resp.Graph == nil {
		t.Error("graph should be included")
	} else {
		if len(resp.Graph.Nodes) != 4 {
			t.Errorf("graph nodes: got %d, want 4", len(resp.Graph.Nodes))
		}
		if len(resp.Graph.Edges) != 3 {
			t.Errorf("graph edges: got %d, want 3", len(resp.Graph.Edges))
		}
	}

	// Dependency tree should be present.
	if resp.DependencyTree == nil {
		t.Error("dependency tree should be present")
	}

	// Confidence detail.
	if resp.ConfidenceDetail == nil {
		t.Error("confidence detail should be present")
	} else {
		if resp.ConfidenceDetail.Overall <= 0 || resp.ConfidenceDetail.Overall > 1.0 {
			t.Errorf("confidence out of range: %f", resp.ConfidenceDetail.Overall)
		}
	}

	// Confidence should match detail.
	if resp.Confidence != resp.ConfidenceDetail.Overall {
		t.Errorf("confidence mismatch: resp=%f, detail=%f", resp.Confidence, resp.ConfidenceDetail.Overall)
	}

	// Metadata.
	if resp.Metadata.Version != "1.0.0" {
		t.Errorf("version: got %s, want 1.0.0", resp.Metadata.Version)
	}
	if resp.Metadata.DeterministicHash == "" {
		t.Error("hash should not be empty")
	}
	if resp.Metadata.CreatedAt.IsZero() {
		t.Error("created_at should not be zero")
	}
}

func TestOrchestrator_Determinism(t *testing.T) {
	orch := NewOrchestrator()
	req := makeTestRequest()

	// Run 100 iterations and verify that the deterministic hash is always the same.
	resp0, err := orch.Explain(req)
	if err != nil {
		t.Fatalf("Explain failed: %v", err)
	}
	expectedHash := resp0.Metadata.DeterministicHash

	for i := 1; i < 100; i++ {
		resp, err := orch.Explain(req)
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		if resp.Metadata.DeterministicHash != expectedHash {
			t.Fatalf("iteration %d: hash mismatch: got %s, want %s",
				i, resp.Metadata.DeterministicHash, expectedHash)
		}
		// ID should be unique (UUID v4).
		if resp.ID == resp0.ID {
			t.Errorf("iteration %d: ID should be unique", i)
		}
	}
}

func TestOrchestrator_NoGraph(t *testing.T) {
	orch := NewOrchestrator()
	req := makeTestRequest()
	req.Options.IncludeGraph = false

	resp, err := orch.Explain(req)
	if err != nil {
		t.Fatalf("Explain failed: %v", err)
	}
	if resp.Graph != nil {
		t.Error("graph should be nil when not requested")
	}
}

func TestOrchestrator_NestedComponents(t *testing.T) {
	orch := NewOrchestrator()
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

	resp, err := orch.Explain(req)
	if err != nil {
		t.Fatalf("Explain failed: %v", err)
	}

	if len(resp.Breakdown) != 2 {
		t.Fatalf("expected 2 top-level breakdown items, got %d", len(resp.Breakdown))
	}

	// Find the "skills" breakdown item (should have children).
	var skillsItem *models.BreakdownItem
	for i := range resp.Breakdown {
		if resp.Breakdown[i].NodeID == "skills" {
			skillsItem = &resp.Breakdown[i]
			break
		}
	}
	if skillsItem == nil {
		t.Fatal("skills breakdown item not found")
	}
	if len(skillsItem.Children) != 2 {
		t.Errorf("skills should have 2 children, got %d", len(skillsItem.Children))
	}
}

func TestOrchestrator_DefaultOptions(t *testing.T) {
	orch := NewOrchestrator()
	req := models.ExplainRequest{
		Target: "score",
		Value:  50.0,
		Components: []models.Component{
			{Name: "a", Value: 50.0, Weight: 1.0, Confidence: 0.9},
		},
	}

	resp, err := orch.Explain(req)
	if err != nil {
		t.Fatalf("Explain failed: %v", err)
	}

	// With default options, graph should be included.
	if resp.Graph == nil {
		t.Error("graph should be included by default")
	}
}
