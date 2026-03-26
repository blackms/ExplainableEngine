package models

// Contribution represents a component's contribution to a computed value.
type Contribution struct {
	NodeID               string         `json:"node_id"`
	Label                string         `json:"label"`
	Value                float64        `json:"value"`
	Weight               float64        `json:"weight"`
	AbsoluteContribution float64        `json:"absolute_contribution"`
	Percentage           float64        `json:"percentage"`
	Confidence           float64        `json:"confidence"`
	Children             []Contribution `json:"children,omitempty"`
}
