package engine

import (
	"fmt"
	"math"
	"sort"

	"github.com/blackms/ExplainableEngine/internal/models"
)

// Modification represents a single input change for what-if analysis.
type Modification struct {
	ComponentName string  `json:"component"`
	NewValue      float64 `json:"new_value"`
}

// ComponentDiff shows the before/after for one component.
type ComponentDiff struct {
	Name                 string  `json:"name"`
	OriginalValue        float64 `json:"original_value"`
	ModifiedValue        float64 `json:"modified_value"`
	DeltaValue           float64 `json:"delta_value"`
	DeltaPercentage      float64 `json:"delta_percentage"`
	OriginalContribution float64 `json:"original_contribution"`
	ModifiedContribution float64 `json:"modified_contribution"`
}

// SensitivityRanking ranks components by how much they changed the output.
type SensitivityRanking struct {
	Name   string  `json:"name"`
	Impact float64 `json:"impact"` // |delta| of contribution
	Rank   int     `json:"rank"`
}

// SensitivityResult is the output of a what-if analysis.
type SensitivityResult struct {
	OriginalValue   float64              `json:"original_value"`
	ModifiedValue   float64              `json:"modified_value"`
	DeltaValue      float64              `json:"delta_value"`
	DeltaPercentage float64              `json:"delta_percentage"`
	ComponentDiffs  []ComponentDiff      `json:"component_diffs"`
	Ranking         []SensitivityRanking `json:"sensitivity_ranking"`
}

// ComponentNotFoundError is returned when a modification references a non-existent component.
type ComponentNotFoundError struct {
	ComponentName string
}

func (e *ComponentNotFoundError) Error() string {
	return fmt.Sprintf("component not found: %s", e.ComponentName)
}

// AnalyzeSensitivity runs a what-if analysis by modifying inputs and recomputing.
func AnalyzeSensitivity(
	originalRequest *models.ExplainRequest,
	originalResponse *models.ExplainResponse,
	modifications []Modification,
	orchestrator OrchestratorInterface,
) (*SensitivityResult, error) {
	// Clone the original request.
	modifiedReq := cloneRequest(originalRequest)

	// Apply modifications.
	for _, mod := range modifications {
		if !applyModification(&modifiedReq, mod) {
			return nil, &ComponentNotFoundError{ComponentName: mod.ComponentName}
		}
	}

	// Recompute the target value as the weighted sum of (possibly modified) components
	// so that the orchestrator produces a FinalValue consistent with the new inputs.
	modifiedReq.Value = computeWeightedSum(modifiedReq.Components)

	// Recompute with modified request.
	modifiedResponse, err := orchestrator.Explain(modifiedReq)
	if err != nil {
		return nil, fmt.Errorf("recomputing with modifications: %w", err)
	}

	// Build component diffs by matching breakdown items by node_id.
	originalBreakdown := flattenBreakdown(originalResponse.Breakdown)
	modifiedBreakdown := flattenBreakdown(modifiedResponse.Breakdown)

	var diffs []ComponentDiff
	for nodeID, origItem := range originalBreakdown {
		modItem, exists := modifiedBreakdown[nodeID]
		if !exists {
			continue
		}
		deltaVal := modItem.Value - origItem.Value
		var deltaPct float64
		if origItem.Value != 0 {
			deltaPct = (deltaVal / origItem.Value) * 100
		}
		diffs = append(diffs, ComponentDiff{
			Name:                 origItem.Label,
			OriginalValue:        origItem.Value,
			ModifiedValue:        modItem.Value,
			DeltaValue:           deltaVal,
			DeltaPercentage:      deltaPct,
			OriginalContribution: origItem.AbsoluteContribution,
			ModifiedContribution: modItem.AbsoluteContribution,
		})
	}
	// Sort diffs by name for determinism.
	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].Name < diffs[j].Name
	})

	// Build ranking: sort by |delta of contribution| descending, tie-break by name.
	type rankEntry struct {
		name   string
		impact float64
	}
	var entries []rankEntry
	for _, d := range diffs {
		impact := math.Abs(d.ModifiedContribution - d.OriginalContribution)
		entries = append(entries, rankEntry{name: d.Name, impact: impact})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].impact != entries[j].impact {
			return entries[i].impact > entries[j].impact
		}
		return entries[i].name < entries[j].name
	})
	ranking := make([]SensitivityRanking, len(entries))
	for i, e := range entries {
		ranking[i] = SensitivityRanking{
			Name:   e.name,
			Impact: e.impact,
			Rank:   i + 1,
		}
	}

	// Overall deltas.
	overallDelta := modifiedResponse.FinalValue - originalResponse.FinalValue
	var overallDeltaPct float64
	if originalResponse.FinalValue != 0 {
		overallDeltaPct = (overallDelta / originalResponse.FinalValue) * 100
	}

	return &SensitivityResult{
		OriginalValue:   originalResponse.FinalValue,
		ModifiedValue:   modifiedResponse.FinalValue,
		DeltaValue:      overallDelta,
		DeltaPercentage: overallDeltaPct,
		ComponentDiffs:  diffs,
		Ranking:         ranking,
	}, nil
}

// computeWeightedSum computes the weighted sum of component values.
func computeWeightedSum(comps []models.Component) float64 {
	var sum float64
	for _, c := range comps {
		sum += c.Value * c.Weight
	}
	return sum
}

// cloneRequest creates a deep copy of an ExplainRequest.
func cloneRequest(req *models.ExplainRequest) models.ExplainRequest {
	clone := models.ExplainRequest{
		Target: req.Target,
		Value:  req.Value,
	}
	clone.Components = cloneComponents(req.Components)
	if req.Options != nil {
		optsCopy := *req.Options
		clone.Options = &optsCopy
	}
	if req.Metadata != nil {
		clone.Metadata = make(map[string]string, len(req.Metadata))
		for k, v := range req.Metadata {
			clone.Metadata[k] = v
		}
	}
	return clone
}

// cloneComponents deep-copies a slice of Component.
func cloneComponents(comps []models.Component) []models.Component {
	if comps == nil {
		return nil
	}
	result := make([]models.Component, len(comps))
	for i, c := range comps {
		result[i] = models.Component{
			ID:         c.ID,
			Name:       c.Name,
			Value:      c.Value,
			Weight:     c.Weight,
			Confidence: c.Confidence,
			Missing:    c.Missing,
			Components: cloneComponents(c.Components),
		}
	}
	return result
}

// applyModification finds a component by name (recursively) and sets its value.
// Returns false if the component is not found.
func applyModification(req *models.ExplainRequest, mod Modification) bool {
	return applyModToComponents(req.Components, mod)
}

func applyModToComponents(comps []models.Component, mod Modification) bool {
	for i := range comps {
		if comps[i].Name == mod.ComponentName {
			comps[i].Value = mod.NewValue
			return true
		}
		if applyModToComponents(comps[i].Components, mod) {
			return true
		}
	}
	return false
}

// flattenBreakdown creates a map from node_id to BreakdownItem for easy lookup.
func flattenBreakdown(items []models.BreakdownItem) map[string]models.BreakdownItem {
	result := make(map[string]models.BreakdownItem)
	for _, item := range items {
		result[item.NodeID] = item
		for k, v := range flattenBreakdown(item.Children) {
			result[k] = v
		}
	}
	return result
}
