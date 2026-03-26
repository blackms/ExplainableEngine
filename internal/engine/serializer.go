package engine

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/blackms/ExplainableEngine/internal/models"
)

// GraphFormat represents a supported export format.
type GraphFormat string

const (
	FormatJSON    GraphFormat = "json"
	FormatDOT     GraphFormat = "dot"
	FormatMermaid GraphFormat = "mermaid"
)

// nodeColor returns the fill color for a given node type.
func nodeColor(nodeType string) string {
	switch models.NodeType(nodeType) {
	case models.NodeTypeInput:
		return "#90EE90"
	case models.NodeTypeOutput:
		return "#FFB6C1"
	case models.NodeTypeComputed:
		return "#ADD8E6"
	case models.NodeTypeMissing:
		return "#D3D3D3"
	default:
		return "#FFFFFF"
	}
}

// SerializeGraph converts a GraphResponse to the requested format.
// Returns (content, contentType, error).
func SerializeGraph(graph *models.GraphResponse, format GraphFormat) (string, string, error) {
	switch format {
	case FormatJSON:
		return serializeJSON(graph)
	case FormatDOT:
		return serializeDOT(graph)
	case FormatMermaid:
		return serializeMermaid(graph)
	default:
		return "", "", fmt.Errorf("unsupported format: %s", string(format))
	}
}

func serializeJSON(graph *models.GraphResponse) (string, string, error) {
	// Sort nodes and edges for determinism.
	g := sortedGraph(graph)
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("marshaling graph to JSON: %w", err)
	}
	return string(data), "application/json", nil
}

func serializeDOT(graph *models.GraphResponse) (string, string, error) {
	g := sortedGraph(graph)
	var b strings.Builder

	b.WriteString("digraph explanation {\n")
	b.WriteString("    rankdir=LR;\n")
	b.WriteString("    node [shape=box, style=filled];\n")

	// Nodes
	for _, n := range g.Nodes {
		label := fmt.Sprintf("%s\\nvalue=%.2f\\nconf=%.2f", n.ID, n.Value, n.Confidence)
		color := nodeColor(n.NodeType)
		b.WriteString(fmt.Sprintf("    %q [label=%q, fillcolor=%q];\n", n.ID, label, color))
	}

	// Edges
	for _, e := range g.Edges {
		b.WriteString(fmt.Sprintf("    %q -> %q [label=\"w=%.2f\"];\n", e.Source, e.Target, e.Weight))
	}

	b.WriteString("}\n")
	return b.String(), "text/vnd.graphviz", nil
}

func serializeMermaid(graph *models.GraphResponse) (string, string, error) {
	g := sortedGraph(graph)
	var b strings.Builder

	b.WriteString("graph LR\n")

	// Node definitions
	for _, n := range g.Nodes {
		label := fmt.Sprintf("%s<br/>value=%.2f<br/>conf=%.2f", n.ID, n.Value, n.Confidence)
		b.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", n.ID, label))
	}

	// Edges
	for _, e := range g.Edges {
		b.WriteString(fmt.Sprintf("    %s -->|w=%.2f| %s\n", e.Source, e.Weight, e.Target))
	}

	// Styles
	for _, n := range g.Nodes {
		color := nodeColor(n.NodeType)
		b.WriteString(fmt.Sprintf("    style %s fill:%s\n", n.ID, color))
	}

	return b.String(), "text/plain", nil
}

// sortedGraph returns a copy of the graph with nodes sorted by ID and edges
// sorted by (source, target) for deterministic output.
func sortedGraph(graph *models.GraphResponse) *models.GraphResponse {
	nodes := make([]models.GraphNodeResponse, len(graph.Nodes))
	copy(nodes, graph.Nodes)
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].ID < nodes[j].ID
	})

	edges := make([]models.GraphEdgeResponse, len(graph.Edges))
	copy(edges, graph.Edges)
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].Source != edges[j].Source {
			return edges[i].Source < edges[j].Source
		}
		return edges[i].Target < edges[j].Target
	})

	return &models.GraphResponse{
		Nodes: nodes,
		Edges: edges,
	}
}
