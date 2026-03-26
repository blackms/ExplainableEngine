package storage

import "fmt"

// NewStore creates an ExplanationStore for the given backend.
// Supported backends: "memory", "sqlite".
// For "sqlite", sqlitePath must be a valid file path (or ":memory:" for in-process).
func NewStore(backend string, sqlitePath string) (ExplanationStore, error) {
	switch backend {
	case "memory":
		return NewInMemoryStore(), nil
	case "sqlite":
		if sqlitePath == "" {
			return nil, fmt.Errorf("sqlite backend requires a non-empty path")
		}
		return NewSQLiteStore(sqlitePath)
	default:
		return nil, fmt.Errorf("unsupported storage backend: %q", backend)
	}
}
