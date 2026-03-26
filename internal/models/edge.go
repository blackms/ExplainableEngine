package models

// TransformationType represents how values are transformed along an edge.
type TransformationType string

const (
	TransformationWeightedSum   TransformationType = "weighted_sum"
	TransformationNormalization TransformationType = "normalization"
	TransformationThreshold     TransformationType = "threshold"
	TransformationCustom        TransformationType = "custom"
)

// Edge represents a directed connection between two nodes.
type Edge struct {
	Source             string             `json:"source"`
	Target             string             `json:"target"`
	Weight             float64            `json:"weight"`
	TransformationType TransformationType `json:"transformation_type"`
	Metadata           map[string]string  `json:"metadata,omitempty"`
}
