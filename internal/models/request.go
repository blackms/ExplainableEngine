package models

// Component represents an input component that contributes to the target value.
type Component struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name"`
	Value      float64     `json:"value"`
	Weight     float64     `json:"weight"`
	Confidence float64     `json:"confidence"`
	Components []Component `json:"components,omitempty"`
}

// ExplainOptions configures what to include in the explanation response.
type ExplainOptions struct {
	IncludeGraph   bool `json:"include_graph"`
	IncludeDrivers bool `json:"include_drivers"`
	MaxDrivers     int  `json:"max_drivers"`
	MaxDepth       int  `json:"max_depth"`
}

// DefaultExplainOptions returns sensible defaults.
func DefaultExplainOptions() ExplainOptions {
	return ExplainOptions{
		IncludeGraph:   true,
		IncludeDrivers: true,
		MaxDrivers:     5,
		MaxDepth:       10,
	}
}

// ExplainRequest is the input payload for creating an explanation.
type ExplainRequest struct {
	Target     string            `json:"target"`
	Value      float64           `json:"value"`
	Components []Component       `json:"components"`
	Options    *ExplainOptions   `json:"options,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// GetOptions returns the options, using defaults if not specified.
func (r *ExplainRequest) GetOptions() ExplainOptions {
	if r.Options != nil {
		opts := *r.Options
		if opts.MaxDrivers == 0 {
			opts.MaxDrivers = 5
		}
		if opts.MaxDepth == 0 {
			opts.MaxDepth = 10
		}
		return opts
	}
	return DefaultExplainOptions()
}
