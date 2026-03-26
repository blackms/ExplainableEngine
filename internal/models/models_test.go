package models_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/models"
)

func TestNodeRoundTrip(t *testing.T) {
	node := models.Node{
		ID:         "trend",
		Label:      "Trend Strength",
		Value:      0.8,
		Confidence: 0.9,
		NodeType:   models.NodeTypeInput,
	}
	data, err := json.Marshal(node)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var restored models.Node
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(restored, node) {
		t.Errorf("roundtrip mismatch: got %+v, want %+v", restored, node)
	}
}

func TestEdgeRoundTrip(t *testing.T) {
	edge := models.Edge{
		Source:             "a",
		Target:             "b",
		Weight:             0.4,
		TransformationType: models.TransformationWeightedSum,
	}
	data, err := json.Marshal(edge)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var restored models.Edge
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(restored, edge) {
		t.Errorf("roundtrip mismatch: got %+v, want %+v", restored, edge)
	}
}

func TestExplainRequestDefaults(t *testing.T) {
	req := models.ExplainRequest{
		Target: "score",
		Value:  0.72,
		Components: []models.Component{
			{Name: "a", Value: 0.8, Weight: 0.4, Confidence: 0.9},
		},
	}
	opts := req.GetOptions()
	if opts.MaxDrivers != 5 {
		t.Errorf("default MaxDrivers: got %d, want 5", opts.MaxDrivers)
	}
	if opts.MaxDepth != 10 {
		t.Errorf("default MaxDepth: got %d, want 10", opts.MaxDepth)
	}
	if !opts.IncludeGraph {
		t.Error("default IncludeGraph should be true")
	}
}

func TestComponentNested(t *testing.T) {
	comp := models.Component{
		Name:   "skills",
		Value:  82.0,
		Weight: 0.3,
		Components: []models.Component{
			{Name: "python", Value: 95.0, Weight: 0.5, Confidence: 0.92},
			{Name: "go", Value: 70.0, Weight: 0.5, Confidence: 0.85},
		},
	}
	data, err := json.Marshal(comp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var restored models.Component
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(restored.Components) != 2 {
		t.Errorf("nested components: got %d, want 2", len(restored.Components))
	}
}

func TestExplainResponseJSON(t *testing.T) {
	resp := models.ExplainResponse{
		ID:         "test-123",
		Target:     "score",
		FinalValue: 0.72,
		Confidence: 0.828,
		Breakdown: []models.BreakdownItem{
			{
				NodeID:               "a",
				Label:                "A",
				Value:                0.8,
				Weight:               0.4,
				AbsoluteContribution: 0.32,
				Percentage:           44.4,
				Confidence:           0.9,
			},
		},
		TopDrivers: []models.DriverItem{
			{Name: "a", Impact: 0.44, Rank: 1},
		},
	}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var restored models.ExplainResponse
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if restored.ID != resp.ID {
		t.Errorf("id mismatch: got %s, want %s", restored.ID, resp.ID)
	}
	if restored.Confidence != resp.Confidence {
		t.Errorf("confidence mismatch: got %f, want %f", restored.Confidence, resp.Confidence)
	}
}
