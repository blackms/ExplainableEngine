package engine

import (
	"fmt"
	"sort"

	"github.com/blackms/ExplainableEngine/internal/models"
)

// EdgeInfo stores edge metadata in the adjacency lists.
type EdgeInfo struct {
	Source             string
	Target             string
	Weight             float64
	TransformationType models.TransformationType
}

// DAG is a directed acyclic graph backed by adjacency lists.
type DAG struct {
	nodes    map[string]*models.Node
	edges    map[string][]EdgeInfo // target -> incoming edges
	outEdges map[string][]EdgeInfo // source -> outgoing edges
}

// NewDAG creates an empty DAG.
func NewDAG() *DAG {
	return &DAG{
		nodes:    make(map[string]*models.Node),
		edges:    make(map[string][]EdgeInfo),
		outEdges: make(map[string][]EdgeInfo),
	}
}

// AddNode adds a node to the graph. Returns an error if the node ID already exists.
func (d *DAG) AddNode(node models.Node) error {
	if _, exists := d.nodes[node.ID]; exists {
		return fmt.Errorf("duplicate node: %s", node.ID)
	}
	n := node // copy
	d.nodes[node.ID] = &n
	return nil
}

// AddEdge adds a directed edge. Returns an error if source or target nodes do not exist.
func (d *DAG) AddEdge(edge models.Edge) error {
	if _, exists := d.nodes[edge.Source]; !exists {
		return &models.NodeNotFoundError{NodeID: edge.Source}
	}
	if _, exists := d.nodes[edge.Target]; !exists {
		return &models.NodeNotFoundError{NodeID: edge.Target}
	}
	info := EdgeInfo{
		Source:             edge.Source,
		Target:             edge.Target,
		Weight:             edge.Weight,
		TransformationType: edge.TransformationType,
	}
	d.edges[edge.Target] = append(d.edges[edge.Target], info)
	d.outEdges[edge.Source] = append(d.outEdges[edge.Source], info)
	return nil
}

// GetNode retrieves a node by ID.
func (d *DAG) GetNode(id string) (models.Node, error) {
	n, exists := d.nodes[id]
	if !exists {
		return models.Node{}, &models.NodeNotFoundError{NodeID: id}
	}
	return *n, nil
}

// GetPredecessors returns the direct parents of a node, sorted by ID.
func (d *DAG) GetPredecessors(nodeID string) ([]string, error) {
	if _, exists := d.nodes[nodeID]; !exists {
		return nil, &models.NodeNotFoundError{NodeID: nodeID}
	}
	seen := make(map[string]bool)
	for _, e := range d.edges[nodeID] {
		seen[e.Source] = true
	}
	result := make([]string, 0, len(seen))
	for id := range seen {
		result = append(result, id)
	}
	sort.Strings(result)
	return result, nil
}

// GetAncestors returns all ancestors of a node (transitive predecessors), sorted by ID.
func (d *DAG) GetAncestors(nodeID string) ([]string, error) {
	if _, exists := d.nodes[nodeID]; !exists {
		return nil, &models.NodeNotFoundError{NodeID: nodeID}
	}
	visited := make(map[string]bool)
	d.collectAncestors(nodeID, visited)
	result := make([]string, 0, len(visited))
	for id := range visited {
		result = append(result, id)
	}
	sort.Strings(result)
	return result, nil
}

func (d *DAG) collectAncestors(nodeID string, visited map[string]bool) {
	for _, e := range d.edges[nodeID] {
		if !visited[e.Source] {
			visited[e.Source] = true
			d.collectAncestors(e.Source, visited)
		}
	}
}

// GetDescendants returns all descendants of a node (transitive successors), sorted by ID.
func (d *DAG) GetDescendants(nodeID string) ([]string, error) {
	if _, exists := d.nodes[nodeID]; !exists {
		return nil, &models.NodeNotFoundError{NodeID: nodeID}
	}
	visited := make(map[string]bool)
	d.collectDescendants(nodeID, visited)
	result := make([]string, 0, len(visited))
	for id := range visited {
		result = append(result, id)
	}
	sort.Strings(result)
	return result, nil
}

func (d *DAG) collectDescendants(nodeID string, visited map[string]bool) {
	for _, e := range d.outEdges[nodeID] {
		if !visited[e.Target] {
			visited[e.Target] = true
			d.collectDescendants(e.Target, visited)
		}
	}
}

// TopologicalSort returns nodes in topological order using Kahn's algorithm.
// Returns a CyclicGraphError if the graph contains a cycle.
func (d *DAG) TopologicalSort() ([]string, error) {
	// Compute in-degrees.
	inDegree := make(map[string]int, len(d.nodes))
	for id := range d.nodes {
		inDegree[id] = 0
	}
	for _, edges := range d.outEdges {
		for _, e := range edges {
			inDegree[e.Target]++
		}
	}

	// Collect nodes with zero in-degree, sorted for determinism.
	var queue []string
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}
	sort.Strings(queue)

	var result []string
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		// Collect neighbors, sort for determinism.
		var neighbors []string
		for _, e := range d.outEdges[node] {
			inDegree[e.Target]--
			if inDegree[e.Target] == 0 {
				neighbors = append(neighbors, e.Target)
			}
		}
		sort.Strings(neighbors)
		queue = append(queue, neighbors...)
	}

	if len(result) != len(d.nodes) {
		return nil, &models.CyclicGraphError{Message: "graph contains a cycle"}
	}
	return result, nil
}

// HasCycle returns true if the graph contains a cycle.
func (d *DAG) HasCycle() bool {
	_, err := d.TopologicalSort()
	return err != nil
}

// BuildFromRequest constructs a DAG from an ExplainRequest.
// It creates a root OUTPUT node and child INPUT nodes, with edges carrying weights.
// Nested components are handled recursively.
func (d *DAG) BuildFromRequest(req models.ExplainRequest) (*DAG, error) {
	dag := NewDAG()

	// Create root output node.
	rootID := req.Target
	rootNode := models.Node{
		ID:         rootID,
		Label:      req.Target,
		Value:      req.Value,
		Confidence: 0, // will be computed later
		NodeType:   models.NodeTypeOutput,
	}
	if err := dag.AddNode(rootNode); err != nil {
		return nil, fmt.Errorf("adding root node: %w", err)
	}

	// Add components recursively.
	if err := addComponents(dag, rootID, req.Components); err != nil {
		return nil, err
	}

	return dag, nil
}

// addComponents recursively adds component nodes and edges to the DAG.
func addComponents(dag *DAG, parentID string, components []models.Component) error {
	for _, comp := range components {
		nodeID := comp.Name
		if comp.ID != "" {
			nodeID = comp.ID
		}

		nodeType := models.NodeTypeInput
		if comp.Missing {
			nodeType = models.NodeTypeMissing
		} else if len(comp.Components) > 0 {
			nodeType = models.NodeTypeComputed
		}

		node := models.Node{
			ID:         nodeID,
			Label:      comp.Name,
			Value:      comp.Value,
			Confidence: comp.Confidence,
			NodeType:   nodeType,
		}
		if err := dag.AddNode(node); err != nil {
			return fmt.Errorf("adding component node %s: %w", nodeID, err)
		}

		edge := models.Edge{
			Source:             nodeID,
			Target:             parentID,
			Weight:             comp.Weight,
			TransformationType: models.TransformationWeightedSum,
		}
		if err := dag.AddEdge(edge); err != nil {
			return fmt.Errorf("adding edge %s -> %s: %w", nodeID, parentID, err)
		}

		// Recurse into sub-components.
		if len(comp.Components) > 0 {
			if err := addComponents(dag, nodeID, comp.Components); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetIncomingEdges returns incoming edges for a node.
func (d *DAG) GetIncomingEdges(nodeID string) []EdgeInfo {
	return d.edges[nodeID]
}

// Nodes returns all node IDs sorted.
func (d *DAG) Nodes() []string {
	result := make([]string, 0, len(d.nodes))
	for id := range d.nodes {
		result = append(result, id)
	}
	sort.Strings(result)
	return result
}

// AllEdges returns all edges in the graph.
func (d *DAG) AllEdges() []EdgeInfo {
	var result []EdgeInfo
	for _, edges := range d.outEdges {
		result = append(result, edges...)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Source != result[j].Source {
			return result[i].Source < result[j].Source
		}
		return result[i].Target < result[j].Target
	})
	return result
}
