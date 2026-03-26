package models

// NodeType represents the type of a node in the explanation graph.
type NodeType string

const (
	NodeTypeInput    NodeType = "input"
	NodeTypeComputed NodeType = "computed"
	NodeTypeOutput   NodeType = "output"
	NodeTypeMissing  NodeType = "missing"
)

// Node represents a single node in the explanation graph.
type Node struct {
	ID         string            `json:"id"`
	Label      string            `json:"label"`
	Value      float64           `json:"value"`
	Confidence float64           `json:"confidence"`
	NodeType   NodeType          `json:"node_type"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}
