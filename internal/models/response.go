package models

import "time"

// BreakdownItem represents one component's contribution to the result.
type BreakdownItem struct {
	NodeID               string          `json:"node_id"`
	Label                string          `json:"label"`
	Value                float64         `json:"value"`
	Weight               float64         `json:"weight"`
	AbsoluteContribution float64         `json:"absolute_contribution"`
	Percentage           float64         `json:"percentage"`
	Confidence           float64         `json:"confidence"`
	Children             []BreakdownItem `json:"children,omitempty"`
}

// DriverItem represents a top driver of the result.
type DriverItem struct {
	Name   string  `json:"name"`
	Impact float64 `json:"impact"`
	Rank   int     `json:"rank"`
}

// GraphNodeResponse is the API representation of a graph node.
type GraphNodeResponse struct {
	ID         string  `json:"id"`
	Label      string  `json:"label"`
	Value      float64 `json:"value"`
	Confidence float64 `json:"confidence"`
	NodeType   string  `json:"node_type"`
}

// GraphEdgeResponse is the API representation of a graph edge.
type GraphEdgeResponse struct {
	Source             string  `json:"source"`
	Target             string  `json:"target"`
	Weight             float64 `json:"weight"`
	TransformationType string  `json:"transformation_type"`
}

// GraphResponse contains the full graph representation.
type GraphResponse struct {
	Nodes []GraphNodeResponse `json:"nodes"`
	Edges []GraphEdgeResponse `json:"edges"`
}

// DependencyNode represents a node in the dependency tree.
type DependencyNode struct {
	NodeID   string           `json:"node_id"`
	Label    string           `json:"label"`
	Depth    int              `json:"depth"`
	Relation string           `json:"relation,omitempty"`
	Children []DependencyNode `json:"children,omitempty"`
}

// DependencyTree represents the full dependency tree.
type DependencyTree struct {
	Root       DependencyNode `json:"root"`
	Depth      int            `json:"depth"`
	TotalNodes int            `json:"total_nodes"`
}

// ConfidenceDetail contains per-node confidence information.
type ConfidenceDetail struct {
	Overall float64            `json:"overall"`
	PerNode map[string]float64 `json:"per_node"`
}

// ExplainMetadata contains metadata about the explanation.
type ExplainMetadata struct {
	Version           string    `json:"version"`
	CreatedAt         time.Time `json:"created_at"`
	DeterministicHash string    `json:"deterministic_hash"`
	ComputationType   string    `json:"computation_type"`
}

// ExplainResponse is the full explanation output.
type ExplainResponse struct {
	ID               string            `json:"id"`
	Target           string            `json:"target"`
	FinalValue       float64           `json:"final_value"`
	Confidence       float64           `json:"confidence"`
	Breakdown        []BreakdownItem   `json:"breakdown"`
	TopDrivers       []DriverItem      `json:"top_drivers"`
	MissingImpact    float64           `json:"missing_impact"`
	Graph            *GraphResponse    `json:"graph,omitempty"`
	DependencyTree   *DependencyTree   `json:"dependency_tree,omitempty"`
	ConfidenceDetail *ConfidenceDetail `json:"confidence_detail,omitempty"`
	Metadata         ExplainMetadata   `json:"metadata"`
	OriginalRequest  *ExplainRequest   `json:"original_request,omitempty"`
}
