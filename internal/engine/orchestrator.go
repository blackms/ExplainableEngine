package engine

import (
	crypto_rand "crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/blackms/ExplainableEngine/internal/models"
)

// Orchestrator coordinates the full explanation pipeline.
type Orchestrator struct{}

// NewOrchestrator creates a new Orchestrator.
func NewOrchestrator() *Orchestrator {
	return &Orchestrator{}
}

// Explain runs the full explanation pipeline for the given request.
func (o *Orchestrator) Explain(req models.ExplainRequest) (*models.ExplainResponse, error) {
	opts := req.GetOptions()

	// 1. Build DAG from request.
	dag := NewDAG()
	dag, err := dag.BuildFromRequest(req)
	if err != nil {
		return nil, fmt.Errorf("building DAG: %w", err)
	}

	rootNodeID := req.Target

	// 2. Compute breakdown.
	contributions, err := ComputeBreakdown(dag, rootNodeID)
	if err != nil {
		return nil, fmt.Errorf("computing breakdown: %w", err)
	}

	// 3. Resolve dependencies.
	depTree, err := ResolveDependencies(dag, rootNodeID)
	if err != nil {
		return nil, fmt.Errorf("resolving dependencies: %w", err)
	}

	// 4. Propagate confidence.
	confResult, err := PropagateConfidence(dag, rootNodeID)
	if err != nil {
		return nil, fmt.Errorf("propagating confidence: %w", err)
	}

	// 5. Analyze missing data.
	missingResult, err := AnalyzeMissingData(dag, rootNodeID, opts.MissingThreshold)
	if err != nil {
		return nil, fmt.Errorf("analyzing missing data: %w", err)
	}

	// 6. Build breakdown items from contributions.
	breakdown := contributionsToBreakdownItems(contributions)

	// 7. Compute top drivers using the dedicated analyzer.
	driverItems := AnalyzeDrivers(contributions, confResult.NodeConfidences, opts.MaxDrivers)
	if driverItems == nil {
		driverItems = []models.DriverItem{}
	}

	// 8. Build graph response if requested.
	var graphResp *models.GraphResponse
	if opts.IncludeGraph {
		graphResp = buildGraphResponse(dag)
	}

	// 9. Build confidence detail.
	confDetail := &models.ConfidenceDetail{
		Overall: confResult.OverallConfidence,
		PerNode: confResult.NodeConfidences,
	}

	// 10. Compute deterministic hash (of breakdown + confidence, excluding id/timestamp).
	hashInput := struct {
		Breakdown  []models.BreakdownItem   `json:"breakdown"`
		Confidence *models.ConfidenceDetail  `json:"confidence"`
	}{
		Breakdown:  breakdown,
		Confidence: confDetail,
	}
	hashBytes, _ := json.Marshal(hashInput)
	hash := fmt.Sprintf("%x", sha256.Sum256(hashBytes))

	// 11. Store original request for what-if analysis.
	reqCopy := req

	// 12. Assemble response.
	return &models.ExplainResponse{
		ID:               newUUID(),
		Target:           req.Target,
		FinalValue:       req.Value,
		Confidence:       confResult.OverallConfidence,
		Breakdown:        breakdown,
		TopDrivers:       driverItems,
		MissingImpact:    missingResult.TotalImpact,
		Graph:            graphResp,
		DependencyTree:   depTree,
		ConfidenceDetail: confDetail,
		Metadata: models.ExplainMetadata{
			Version:           "1.0.0",
			CreatedAt:         time.Now(),
			DeterministicHash: hash,
			ComputationType:   "weighted_sum",
		},
		OriginalRequest: &reqCopy,
	}, nil
}

func contributionsToBreakdownItems(contributions []models.Contribution) []models.BreakdownItem {
	if contributions == nil {
		return nil
	}
	items := make([]models.BreakdownItem, len(contributions))
	for i, c := range contributions {
		items[i] = models.BreakdownItem{
			NodeID:               c.NodeID,
			Label:                c.Label,
			Value:                c.Value,
			Weight:               c.Weight,
			AbsoluteContribution: c.AbsoluteContribution,
			Percentage:           c.Percentage,
			Confidence:           c.Confidence,
			Children:             contributionsToBreakdownItems(c.Children),
		}
	}
	return items
}

func buildGraphResponse(dag *DAG) *models.GraphResponse {
	nodeIDs := dag.Nodes()
	nodes := make([]models.GraphNodeResponse, len(nodeIDs))
	for i, id := range nodeIDs {
		n, _ := dag.GetNode(id)
		nodes[i] = models.GraphNodeResponse{
			ID:         n.ID,
			Label:      n.Label,
			Value:      n.Value,
			Confidence: n.Confidence,
			NodeType:   string(n.NodeType),
		}
	}

	allEdges := dag.AllEdges()
	edges := make([]models.GraphEdgeResponse, len(allEdges))
	for i, e := range allEdges {
		edges[i] = models.GraphEdgeResponse{
			Source:             e.Source,
			Target:             e.Target,
			Weight:             e.Weight,
			TransformationType: string(e.TransformationType),
		}
	}

	return &models.GraphResponse{
		Nodes: nodes,
		Edges: edges,
	}
}

// newUUID generates a v4 UUID using crypto/rand.
func newUUID() string {
	b := make([]byte, 16)
	crypto_rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
