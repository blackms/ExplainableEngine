package engine

import (
	"math"
	"sort"

	"github.com/blackms/ExplainableEngine/internal/models"
)

// AnalyzeDrivers computes the top N drivers ranked by impact score.
// Impact score = |absolute_contribution| * confidence, normalized to [0, 1].
// Drivers are sorted descending by impact, with tie-breaking by name (ascending)
// for determinism.
func AnalyzeDrivers(
	contributions []models.Contribution,
	confidences map[string]float64,
	topN int,
) []models.DriverItem {
	// Flatten all contributions (including recursive children).
	var candidates []driverCandidate
	collectAllDrivers(contributions, confidences, &candidates)

	if len(candidates) == 0 {
		return nil
	}

	// Sort descending by impact, tie-break by name ascending.
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].impact != candidates[j].impact {
			return candidates[i].impact > candidates[j].impact
		}
		return candidates[i].name < candidates[j].name
	})

	// Take top N.
	n := topN
	if n > len(candidates) || n <= 0 {
		n = len(candidates)
	}
	candidates = candidates[:n]

	// Normalize: divide all impacts by the max so the top driver has impact = 1.0.
	maxImpact := candidates[0].impact
	items := make([]models.DriverItem, n)
	for i, c := range candidates {
		normalized := 0.0
		if maxImpact > 0 {
			normalized = c.impact / maxImpact
		}
		items[i] = models.DriverItem{
			Name:   c.name,
			Impact: normalized,
			Rank:   i + 1,
		}
	}

	return items
}

// driverCandidate holds intermediate driver data for sorting.
type driverCandidate struct {
	name   string
	impact float64
}

// collectAllDrivers recursively flattens contributions into driver candidates.
func collectAllDrivers(contributions []models.Contribution, confidences map[string]float64, out *[]driverCandidate) {
	for _, c := range contributions {
		conf := confidences[c.NodeID]
		impact := math.Abs(c.AbsoluteContribution) * conf
		*out = append(*out, driverCandidate{name: c.Label, impact: impact})

		// Recurse into children for multi-level support.
		if len(c.Children) > 0 {
			collectAllDrivers(c.Children, confidences, out)
		}
	}
}
