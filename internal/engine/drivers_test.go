package engine

import (
	"testing"

	"github.com/blackms/ExplainableEngine/internal/models"
)

func TestAnalyzeDrivers_UniformContributions(t *testing.T) {
	contributions := []models.Contribution{
		{NodeID: "a", Label: "alpha", AbsoluteContribution: 10.0},
		{NodeID: "b", Label: "beta", AbsoluteContribution: 10.0},
		{NodeID: "c", Label: "gamma", AbsoluteContribution: 10.0},
	}
	confidences := map[string]float64{
		"a": 1.0,
		"b": 1.0,
		"c": 1.0,
	}

	drivers := AnalyzeDrivers(contributions, confidences, 5)

	if len(drivers) != 3 {
		t.Fatalf("expected 3 drivers, got %d", len(drivers))
	}
	// All impacts should be equal (normalized to 1.0).
	for _, d := range drivers {
		if d.Impact != 1.0 {
			t.Errorf("driver %s: expected impact 1.0, got %f", d.Name, d.Impact)
		}
	}
}

func TestAnalyzeDrivers_OneDominant(t *testing.T) {
	contributions := []models.Contribution{
		{NodeID: "a", Label: "alpha", AbsoluteContribution: 100.0},
		{NodeID: "b", Label: "beta", AbsoluteContribution: 10.0},
		{NodeID: "c", Label: "gamma", AbsoluteContribution: 5.0},
	}
	confidences := map[string]float64{
		"a": 1.0,
		"b": 1.0,
		"c": 1.0,
	}

	drivers := AnalyzeDrivers(contributions, confidences, 5)

	if len(drivers) != 3 {
		t.Fatalf("expected 3 drivers, got %d", len(drivers))
	}
	if drivers[0].Name != "alpha" {
		t.Errorf("expected rank 1 driver 'alpha', got %s", drivers[0].Name)
	}
	if drivers[0].Impact != 1.0 {
		t.Errorf("expected rank 1 impact 1.0, got %f", drivers[0].Impact)
	}
	if drivers[0].Rank != 1 {
		t.Errorf("expected rank 1, got %d", drivers[0].Rank)
	}

	// beta: impact = 10/100 = 0.1
	expectedBeta := 10.0 / 100.0
	if diff := drivers[1].Impact - expectedBeta; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("expected beta impact %f, got %f", expectedBeta, drivers[1].Impact)
	}
}

func TestAnalyzeDrivers_AllZeroContributions(t *testing.T) {
	contributions := []models.Contribution{
		{NodeID: "a", Label: "alpha", AbsoluteContribution: 0.0},
		{NodeID: "b", Label: "beta", AbsoluteContribution: 0.0},
	}
	confidences := map[string]float64{
		"a": 1.0,
		"b": 1.0,
	}

	drivers := AnalyzeDrivers(contributions, confidences, 5)

	if len(drivers) != 2 {
		t.Fatalf("expected 2 drivers, got %d", len(drivers))
	}
	for _, d := range drivers {
		if d.Impact != 0.0 {
			t.Errorf("driver %s: expected impact 0.0, got %f", d.Name, d.Impact)
		}
	}
}

func TestAnalyzeDrivers_TieBreakByName(t *testing.T) {
	contributions := []models.Contribution{
		{NodeID: "c", Label: "charlie", AbsoluteContribution: 10.0},
		{NodeID: "a", Label: "alpha", AbsoluteContribution: 10.0},
		{NodeID: "b", Label: "bravo", AbsoluteContribution: 10.0},
	}
	confidences := map[string]float64{
		"a": 1.0,
		"b": 1.0,
		"c": 1.0,
	}

	drivers := AnalyzeDrivers(contributions, confidences, 5)

	if len(drivers) != 3 {
		t.Fatalf("expected 3 drivers, got %d", len(drivers))
	}
	// All same impact, should be sorted by name ascending.
	expectedOrder := []string{"alpha", "bravo", "charlie"}
	for i, name := range expectedOrder {
		if drivers[i].Name != name {
			t.Errorf("position %d: expected %s, got %s", i, name, drivers[i].Name)
		}
		if drivers[i].Rank != i+1 {
			t.Errorf("position %d: expected rank %d, got %d", i, i+1, drivers[i].Rank)
		}
	}
}

func TestAnalyzeDrivers_TopNLimiting(t *testing.T) {
	contributions := []models.Contribution{
		{NodeID: "a", Label: "alpha", AbsoluteContribution: 50.0},
		{NodeID: "b", Label: "beta", AbsoluteContribution: 30.0},
		{NodeID: "c", Label: "gamma", AbsoluteContribution: 20.0},
		{NodeID: "d", Label: "delta", AbsoluteContribution: 10.0},
	}
	confidences := map[string]float64{
		"a": 1.0,
		"b": 1.0,
		"c": 1.0,
		"d": 1.0,
	}

	drivers := AnalyzeDrivers(contributions, confidences, 2)

	if len(drivers) != 2 {
		t.Fatalf("expected 2 drivers (topN=2), got %d", len(drivers))
	}
	if drivers[0].Name != "alpha" {
		t.Errorf("expected rank 1 'alpha', got %s", drivers[0].Name)
	}
	if drivers[1].Name != "beta" {
		t.Errorf("expected rank 2 'beta', got %s", drivers[1].Name)
	}
}

func TestAnalyzeDrivers_MultiLevel(t *testing.T) {
	contributions := []models.Contribution{
		{
			NodeID:               "parent",
			Label:                "parent",
			AbsoluteContribution: 40.0,
			Children: []models.Contribution{
				{NodeID: "child1", Label: "child1", AbsoluteContribution: 25.0},
				{NodeID: "child2", Label: "child2", AbsoluteContribution: 15.0},
			},
		},
		{NodeID: "sibling", Label: "sibling", AbsoluteContribution: 20.0},
	}
	confidences := map[string]float64{
		"parent":  0.9,
		"child1":  1.0,
		"child2":  0.8,
		"sibling": 1.0,
	}

	drivers := AnalyzeDrivers(contributions, confidences, 10)

	// Should have 4 drivers: parent, child1, child2, sibling.
	if len(drivers) != 4 {
		t.Fatalf("expected 4 drivers (flattened), got %d", len(drivers))
	}

	// Verify that children appear as drivers.
	names := make(map[string]bool)
	for _, d := range drivers {
		names[d.Name] = true
	}
	for _, expected := range []string{"parent", "child1", "child2", "sibling"} {
		if !names[expected] {
			t.Errorf("expected driver %s to be present", expected)
		}
	}

	// parent: |40| * 0.9 = 36
	// child1: |25| * 1.0 = 25
	// sibling: |20| * 1.0 = 20
	// child2: |15| * 0.8 = 12
	// Rank 1 should be parent (36 is max, normalized to 1.0).
	if drivers[0].Name != "parent" {
		t.Errorf("expected rank 1 'parent', got %s", drivers[0].Name)
	}
	if drivers[0].Impact != 1.0 {
		t.Errorf("expected rank 1 impact 1.0, got %f", drivers[0].Impact)
	}
}

func TestAnalyzeDrivers_EmptyContributions(t *testing.T) {
	drivers := AnalyzeDrivers(nil, nil, 5)
	if drivers != nil {
		t.Errorf("expected nil for empty contributions, got %v", drivers)
	}
}

func TestAnalyzeDrivers_NegativeContributions(t *testing.T) {
	contributions := []models.Contribution{
		{NodeID: "a", Label: "alpha", AbsoluteContribution: -50.0},
		{NodeID: "b", Label: "beta", AbsoluteContribution: 30.0},
	}
	confidences := map[string]float64{
		"a": 1.0,
		"b": 1.0,
	}

	drivers := AnalyzeDrivers(contributions, confidences, 5)

	if len(drivers) != 2 {
		t.Fatalf("expected 2 drivers, got %d", len(drivers))
	}
	// |alpha| = 50, |beta| = 30. Alpha should be rank 1.
	if drivers[0].Name != "alpha" {
		t.Errorf("expected rank 1 'alpha' (negative contribution uses abs), got %s", drivers[0].Name)
	}
}
