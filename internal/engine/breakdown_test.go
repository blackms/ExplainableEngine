package engine

import (
	"math"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/models"
)

func buildSimpleDAG(components []models.Component) *DAG {
	req := models.ExplainRequest{
		Target:     "root",
		Value:      100.0,
		Components: components,
	}
	dag := NewDAG()
	dag, _ = dag.BuildFromRequest(req)
	return dag
}

func TestBreakdown_TwoComponents(t *testing.T) {
	dag := buildSimpleDAG([]models.Component{
		{Name: "a", Value: 80.0, Weight: 0.6, Confidence: 0.9},
		{Name: "b", Value: 60.0, Weight: 0.4, Confidence: 0.8},
	})

	contribs, err := ComputeBreakdown(dag, "root")
	if err != nil {
		t.Fatal(err)
	}
	if len(contribs) != 2 {
		t.Fatalf("expected 2 contributions, got %d", len(contribs))
	}

	// a: 0.6 * 80 = 48, b: 0.4 * 60 = 24, total = 72
	// a pct: 48/72*100 = 66.67, b pct: 24/72*100 = 33.33
	assertFloat(t, "a absolute", contribs[0].AbsoluteContribution, 48.0)
	assertFloat(t, "b absolute", contribs[1].AbsoluteContribution, 24.0)
	assertFloat(t, "a pct", contribs[0].Percentage, 100.0*48.0/72.0)
	assertFloat(t, "b pct", contribs[1].Percentage, 100.0*24.0/72.0)
}

func TestBreakdown_ThreeComponents(t *testing.T) {
	dag := buildSimpleDAG([]models.Component{
		{Name: "a", Value: 10.0, Weight: 0.5, Confidence: 1.0},
		{Name: "b", Value: 20.0, Weight: 0.3, Confidence: 1.0},
		{Name: "c", Value: 30.0, Weight: 0.2, Confidence: 1.0},
	})

	contribs, err := ComputeBreakdown(dag, "root")
	if err != nil {
		t.Fatal(err)
	}
	if len(contribs) != 3 {
		t.Fatalf("expected 3, got %d", len(contribs))
	}

	// a: 5, b: 6, c: 6, total: 17
	assertFloat(t, "a abs", contribs[0].AbsoluteContribution, 5.0)
	assertFloat(t, "b abs", contribs[1].AbsoluteContribution, 6.0)
	assertFloat(t, "c abs", contribs[2].AbsoluteContribution, 6.0)
}

func TestBreakdown_FiveComponents(t *testing.T) {
	comps := []models.Component{
		{Name: "a", Value: 10.0, Weight: 0.2, Confidence: 1.0},
		{Name: "b", Value: 20.0, Weight: 0.2, Confidence: 1.0},
		{Name: "c", Value: 30.0, Weight: 0.2, Confidence: 1.0},
		{Name: "d", Value: 40.0, Weight: 0.2, Confidence: 1.0},
		{Name: "e", Value: 50.0, Weight: 0.2, Confidence: 1.0},
	}
	dag := buildSimpleDAG(comps)
	contribs, err := ComputeBreakdown(dag, "root")
	if err != nil {
		t.Fatal(err)
	}
	if len(contribs) != 5 {
		t.Fatalf("expected 5, got %d", len(contribs))
	}

	// Sorted by node_id: a, b, c, d, e.
	if contribs[0].NodeID != "a" || contribs[4].NodeID != "e" {
		t.Errorf("wrong order: first=%s, last=%s", contribs[0].NodeID, contribs[4].NodeID)
	}
}

func TestBreakdown_ZeroWeight(t *testing.T) {
	dag := buildSimpleDAG([]models.Component{
		{Name: "a", Value: 80.0, Weight: 0.0, Confidence: 0.9},
		{Name: "b", Value: 60.0, Weight: 0.0, Confidence: 0.8},
	})

	contribs, err := ComputeBreakdown(dag, "root")
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range contribs {
		if c.Percentage != 0 {
			t.Errorf("expected 0 percentage for zero weight, got %f", c.Percentage)
		}
	}
}

func TestBreakdown_Recursive(t *testing.T) {
	req := models.ExplainRequest{
		Target: "root",
		Value:  100.0,
		Components: []models.Component{
			{
				Name:   "parent",
				Value:  80.0,
				Weight: 1.0,
				Components: []models.Component{
					{Name: "child_a", Value: 50.0, Weight: 0.6, Confidence: 0.9},
					{Name: "child_b", Value: 30.0, Weight: 0.4, Confidence: 0.8},
				},
			},
		},
	}
	dag := NewDAG()
	dag, _ = dag.BuildFromRequest(req)

	contribs, err := ComputeBreakdown(dag, "root")
	if err != nil {
		t.Fatal(err)
	}
	if len(contribs) != 1 {
		t.Fatalf("expected 1 top-level, got %d", len(contribs))
	}
	if len(contribs[0].Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(contribs[0].Children))
	}

	// child_a: 0.6 * 50 = 30, child_b: 0.4 * 30 = 12
	assertFloat(t, "child_a abs", contribs[0].Children[0].AbsoluteContribution, 30.0)
	assertFloat(t, "child_b abs", contribs[0].Children[1].AbsoluteContribution, 12.0)
}

func assertFloat(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s: got %f, want %f", name, got, want)
	}
}
