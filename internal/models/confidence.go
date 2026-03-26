package models

// PropagationStep records one step of confidence propagation.
type PropagationStep struct {
	NodeID             string   `json:"node_id"`
	ComputedConfidence float64  `json:"computed_confidence"`
	SourceNodes        []string `json:"source_nodes"`
	Formula            string   `json:"formula"`
}

// ConfidenceResult is the output of confidence propagation.
type ConfidenceResult struct {
	OverallConfidence float64            `json:"overall_confidence"`
	NodeConfidences   map[string]float64 `json:"node_confidences"`
	PropagationPath   []PropagationStep  `json:"propagation_path"`
}
