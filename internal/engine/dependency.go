package engine

import (
	"github.com/blackms/ExplainableEngine/internal/models"
)

// ResolveDependencies builds a dependency tree from the DAG starting at nodeID.
func ResolveDependencies(dag *DAG, nodeID string) (*models.DependencyTree, error) {
	node, err := dag.GetNode(nodeID)
	if err != nil {
		return nil, err
	}

	totalNodes := 0
	maxDepth := 0
	root := buildDependencyNode(dag, node, 0, &totalNodes, &maxDepth)

	return &models.DependencyTree{
		Root:       root,
		Depth:      maxDepth,
		TotalNodes: totalNodes,
	}, nil
}

func buildDependencyNode(dag *DAG, node models.Node, depth int, totalNodes *int, maxDepth *int) models.DependencyNode {
	*totalNodes++
	if depth > *maxDepth {
		*maxDepth = depth
	}

	preds, _ := dag.GetPredecessors(node.ID) // already sorted

	var children []models.DependencyNode
	for _, predID := range preds {
		predNode, _ := dag.GetNode(predID)
		child := buildDependencyNode(dag, predNode, depth+1, totalNodes, maxDepth)
		child.Relation = "contributes_to"
		children = append(children, child)
	}

	return models.DependencyNode{
		NodeID:   node.ID,
		Label:    node.Label,
		Depth:    depth,
		Children: children,
	}
}

// TraverseDFS performs a depth-first traversal of the dependency tree.
func TraverseDFS(tree *models.DependencyTree) []models.DependencyNode {
	var result []models.DependencyNode
	dfsVisit(tree.Root, &result)
	return result
}

func dfsVisit(node models.DependencyNode, result *[]models.DependencyNode) {
	*result = append(*result, node)
	for _, child := range node.Children {
		dfsVisit(child, result)
	}
}

// TraverseBFS performs a breadth-first traversal of the dependency tree.
func TraverseBFS(tree *models.DependencyTree) []models.DependencyNode {
	var result []models.DependencyNode
	queue := []models.DependencyNode{tree.Root}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)
		queue = append(queue, current.Children...)
	}
	return result
}
