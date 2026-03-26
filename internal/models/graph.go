package models

// ExplanationGraph represents the full causal graph for an explanation.
type ExplanationGraph struct {
	Nodes      []Node `json:"nodes"`
	Edges      []Edge `json:"edges"`
	RootNodeID string `json:"root_node_id"`
}
