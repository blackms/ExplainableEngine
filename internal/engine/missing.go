package engine

import (
	"fmt"
	"sort"

	"github.com/blackms/ExplainableEngine/internal/models"
)

// MissingDataResult contains the analysis of missing data in the graph.
type MissingDataResult struct {
	MissingNodes []MissingNodeInfo
	TotalImpact  float64 // 0.0 - 1.0
	Warnings     []string
}

// MissingNodeInfo describes a single missing node and its impact.
type MissingNodeInfo struct {
	NodeID string
	Label  string
	Weight float64
	Impact float64 // this node's share of the total impact
}

// DefaultMissingThreshold is used when no threshold is specified.
const DefaultMissingThreshold = 0.2

// AnalyzeMissingData examines the DAG for MISSING nodes and computes their impact.
// It walks the entire graph starting from rootNodeID, collecting missing nodes at
// every level. Impact is computed as the ratio of missing-node edge weights to
// total edge weights among siblings. If the aggregate impact exceeds threshold,
// a warning is added to the result.
func AnalyzeMissingData(dag *DAG, rootNodeID string, threshold float64) (*MissingDataResult, error) {
	if _, err := dag.GetNode(rootNodeID); err != nil {
		return nil, err
	}

	result := &MissingDataResult{}
	visited := make(map[string]bool)

	if err := analyzeMissingRecursive(dag, rootNodeID, visited, result); err != nil {
		return nil, err
	}

	// Sort missing nodes by NodeID for determinism.
	sort.Slice(result.MissingNodes, func(i, j int) bool {
		return result.MissingNodes[i].NodeID < result.MissingNodes[j].NodeID
	})

	// Compute total impact across all levels: sum of all individual impacts.
	// Each node's impact is relative to its own parent's edge-weight total,
	// so we aggregate by averaging weighted contributions. However, the spec
	// says Impact = sum(missing weights) / sum(all weights) across the whole
	// graph walk. We compute it globally from collected data.
	totalImpact := computeGlobalImpact(dag, rootNodeID, visited)
	result.TotalImpact = totalImpact

	if totalImpact > threshold {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("missing data impact %.2f exceeds threshold %.2f", totalImpact, threshold))
	}

	return result, nil
}

// analyzeMissingRecursive walks the graph collecting missing nodes at every level.
func analyzeMissingRecursive(dag *DAG, nodeID string, visited map[string]bool, result *MissingDataResult) error {
	if visited[nodeID] {
		return nil
	}
	visited[nodeID] = true

	inEdges := dag.GetIncomingEdges(nodeID)
	if len(inEdges) == 0 {
		return nil
	}

	// Compute total weight of all direct predecessors at this level.
	totalWeight := 0.0
	for _, e := range inEdges {
		totalWeight += e.Weight
	}

	// Check each predecessor.
	for _, e := range inEdges {
		node, err := dag.GetNode(e.Source)
		if err != nil {
			return err
		}

		if node.NodeType == models.NodeTypeMissing {
			impact := 0.0
			if totalWeight > 0 {
				impact = e.Weight / totalWeight
			}
			result.MissingNodes = append(result.MissingNodes, MissingNodeInfo{
				NodeID: node.ID,
				Label:  node.Label,
				Weight: e.Weight,
				Impact: impact,
			})
		}

		// Recurse into sub-nodes.
		if err := analyzeMissingRecursive(dag, e.Source, visited, result); err != nil {
			return err
		}
	}

	return nil
}

// computeGlobalImpact computes the overall missing impact as:
// sum(weight of edges from MISSING nodes) / sum(weight of ALL edges in the reachable graph).
func computeGlobalImpact(dag *DAG, rootNodeID string, visited map[string]bool) float64 {
	totalWeight := 0.0
	missingWeight := 0.0

	// Walk all visited nodes and sum their incoming edge weights.
	for nodeID := range visited {
		for _, e := range dag.GetIncomingEdges(nodeID) {
			totalWeight += e.Weight
			node, err := dag.GetNode(e.Source)
			if err == nil && node.NodeType == models.NodeTypeMissing {
				missingWeight += e.Weight
			}
		}
	}

	if totalWeight == 0 {
		return 0.0
	}
	return missingWeight / totalWeight
}
