//go:build integration

package storage_test

import (
	"os"
	"testing"

	"github.com/blackms/ExplainableEngine/internal/storage"
)

// These tests require a running PostgreSQL instance.
// Run with: go test -tags=integration ./internal/storage/ -run TestPostgres

func postgresStore(t *testing.T) *storage.PostgresStore {
	t.Helper()
	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		dsn = storage.BuildPostgresDSN()
	}
	store, err := storage.NewPostgresStore(dsn)
	if err != nil {
		t.Fatalf("NewPostgresStore: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestPostgres_SaveGetRoundtrip(t *testing.T) {
	store := postgresStore(t)
	resp := sampleResponse("pg-1")

	if err := store.Save(resp); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Get("pg-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got == nil {
		t.Fatal("Get returned nil")
	}
	if got.ID != "pg-1" {
		t.Errorf("ID mismatch: got %s, want pg-1", got.ID)
	}
	if got.FinalValue != 0.72 {
		t.Errorf("FinalValue mismatch: got %f, want 0.72", got.FinalValue)
	}
	if len(got.Breakdown) != 1 {
		t.Errorf("Breakdown length: got %d, want 1", len(got.Breakdown))
	}
}

func TestPostgres_Exists(t *testing.T) {
	store := postgresStore(t)
	_ = store.Save(sampleResponse("pg-2"))

	exists, err := store.Exists("pg-2")
	if err != nil {
		t.Fatalf("Exists: %v", err)
	}
	if !exists {
		t.Error("Exists returned false for saved entry")
	}

	exists, err = store.Exists("nonexistent")
	if err != nil {
		t.Fatalf("Exists: %v", err)
	}
	if exists {
		t.Error("Exists returned true for nonexistent entry")
	}
}

func TestPostgres_GetNonexistent(t *testing.T) {
	store := postgresStore(t)
	got, err := store.Get("nope")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestPostgres_SaveIdempotent(t *testing.T) {
	store := postgresStore(t)
	resp := sampleResponse("pg-idempotent")

	if err := store.Save(resp); err != nil {
		t.Fatalf("Save (1st): %v", err)
	}
	// Second save with same ID should not error (ON CONFLICT DO NOTHING).
	if err := store.Save(resp); err != nil {
		t.Fatalf("Save (2nd): %v", err)
	}

	got, err := store.Get("pg-idempotent")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got == nil {
		t.Fatal("Get returned nil after idempotent save")
	}
}
