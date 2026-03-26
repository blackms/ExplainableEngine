package engine

import (
	"math"
	"sort"

	"github.com/blackms/ExplainableEngine/internal/models"
)

// ComputeBreakdown computes the contribution breakdown for a given root node.
// It examines direct predecessors, calculates absolute contributions (weight * value),
// and derives percentages. Recurses into sub-components.
func ComputeBreakdown(dag *DAG, rootNodeID string) ([]models.Contribution, error) {
	if _, err := dag.GetNode(rootNodeID); err != nil {
		return nil, err
	}

	inEdges := dag.GetIncomingEdges(rootNodeID)
	if len(inEdges) == 0 {
		return nil, nil
	}

	// Build contributions from incoming edges.
	type entry struct {
		nodeID string
		edge   EdgeInfo
	}
	var entries []entry
	for _, e := range inEdges {
		entries = append(entries, entry{nodeID: e.Source, edge: e})
	}
	// Sort by node ID for determinism.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].nodeID < entries[j].nodeID
	})

	// Compute total of absolute contributions.
	total := 0.0
	type contribData struct {
		node   models.Node
		weight float64
		absC   float64
	}
	data := make([]contribData, len(entries))
	for i, ent := range entries {
		node, _ := dag.GetNode(ent.nodeID)
		absContrib := ent.edge.Weight * node.Value
		total += math.Abs(absContrib)
		data[i] = contribData{node: node, weight: ent.edge.Weight, absC: absContrib}
	}

	contributions := make([]models.Contribution, len(entries))
	for i, d := range data {
		pct := 0.0
		if total != 0 {
			pct = (math.Abs(d.absC) / total) * 100
		}

		// Recurse into children.
		children, err := ComputeBreakdown(dag, d.node.ID)
		if err != nil {
			return nil, err
		}

		contributions[i] = models.Contribution{
			NodeID:               d.node.ID,
			Label:                d.node.Label,
			Value:                d.node.Value,
			Weight:               d.weight,
			AbsoluteContribution: d.absC,
			Percentage:           pct,
			Confidence:           d.node.Confidence,
			Children:             children,
		}
	}

	return contributions, nil
}
