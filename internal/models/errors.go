package models

import "fmt"

// CyclicGraphError is returned when a cycle is detected in the graph.
type CyclicGraphError struct {
	Message string
}

func (e *CyclicGraphError) Error() string {
	return fmt.Sprintf("cyclic graph detected: %s", e.Message)
}

// NodeNotFoundError is returned when a referenced node does not exist.
type NodeNotFoundError struct {
	NodeID string
}

func (e *NodeNotFoundError) Error() string {
	return fmt.Sprintf("node not found: %s", e.NodeID)
}

// ValidationError is returned for invalid input.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}
