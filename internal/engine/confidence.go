package engine

import (
	"fmt"
	"sort"
	"strings"

	"github.com/blackms/ExplainableEngine/internal/models"
)

// PropagateConfidence computes confidence for all nodes using reverse topological order.
// Leaf nodes keep their declared confidence. Parent confidence is the weighted average
// of child confidences. Results are clamped to [0.0, 1.0].
func PropagateConfidence(dag *DAG, rootNodeID string) (*models.ConfidenceResult, error) {
	if _, err := dag.GetNode(rootNodeID); err != nil {
		return nil, err
	}

	order, err := dag.TopologicalSort()
	if err != nil {
		return nil, err
	}

	// Topological order naturally places leaves (no incoming edges) first,
	// which is what we need: process leaves before their dependents.
	confidences := make(map[string]float64, len(order))
	var path []models.PropagationStep

	for _, nodeID := range order {
		node, _ := dag.GetNode(nodeID)
		inEdges := dag.GetIncomingEdges(nodeID)

		if len(inEdges) == 0 {
			// Leaf node: keep declared confidence.
			conf := clamp(node.Confidence, 0.0, 1.0)
			confidences[nodeID] = conf
			path = append(path, models.PropagationStep{
				NodeID:             nodeID,
				ComputedConfidence: conf,
				SourceNodes:        nil,
				Formula:            "leaf_node",
			})
			continue
		}

		// Parent: weighted average of child confidences.
		weightSum := 0.0
		weightedConfSum := 0.0
		var sourceNodes []string

		for _, e := range inEdges {
			sourceNodes = append(sourceNodes, e.Source)
			weightSum += e.Weight
			weightedConfSum += e.Weight * confidences[e.Source]
		}
		sort.Strings(sourceNodes)

		var conf float64
		if weightSum == 0 {
			conf = 0.0
		} else {
			conf = weightedConfSum / weightSum
		}
		conf = clamp(conf, 0.0, 1.0)
		confidences[nodeID] = conf

		// Build formula string.
		var parts []string
		for _, src := range sourceNodes {
			// Find the weight for this source.
			for _, e := range inEdges {
				if e.Source == src {
					parts = append(parts, fmt.Sprintf("%.4f * %.4f", e.Weight, confidences[src]))
					break
				}
			}
		}
		formula := fmt.Sprintf("(%s) / %.4f", strings.Join(parts, " + "), weightSum)

		path = append(path, models.PropagationStep{
			NodeID:             nodeID,
			ComputedConfidence: conf,
			SourceNodes:        sourceNodes,
			Formula:            formula,
		})
	}

	return &models.ConfidenceResult{
		OverallConfidence: confidences[rootNodeID],
		NodeConfidences:   confidences,
		PropagationPath:   path,
	}, nil
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
