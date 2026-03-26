package storage_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/blackms/ExplainableEngine/internal/models"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

func sampleResponse(id string) *models.ExplainResponse {
	return &models.ExplainResponse{
		ID:         id,
		Target:     "score",
		FinalValue: 0.72,
		Confidence: 0.85,
		Breakdown: []models.BreakdownItem{
			{
				NodeID:               "a",
				Label:                "Component A",
				Value:                0.8,
				Weight:               0.4,
				AbsoluteContribution: 0.32,
				Percentage:           44.4,
				Confidence:           0.9,
			},
		},
		TopDrivers: []models.DriverItem{
			{Name: "Component A", Impact: 0.44, Rank: 1},
		},
		Metadata: models.ExplainMetadata{
			Version:   "0.1.0",
			CreatedAt: time.Now(),
		},
	}
}

// ---- InMemoryStore Tests ----

func TestInMemory_SaveGetRoundtrip(t *testing.T) {
	store := storage.NewInMemoryStore()
	resp := sampleResponse("test-1")

	if err := store.Save(resp); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Get("test-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got == nil {
		t.Fatal("Get returned nil")
	}
	if got.ID != "test-1" {
		t.Errorf("ID mismatch: got %s, want test-1", got.ID)
	}
	if got.FinalValue != 0.72 {
		t.Errorf("FinalValue mismatch: got %f, want 0.72", got.FinalValue)
	}
}

func TestInMemory_Exists(t *testing.T) {
	store := storage.NewInMemoryStore()
	resp := sampleResponse("test-2")
	_ = store.Save(resp)

	exists, err := store.Exists("test-2")
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

func TestInMemory_GetNonexistent(t *testing.T) {
	store := storage.NewInMemoryStore()
	got, err := store.Get("does-not-exist")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for nonexistent key, got %+v", got)
	}
}

func TestInMemory_LRUEviction(t *testing.T) {
	store := storage.NewInMemoryStore(storage.WithMaxSize(3))

	for i := 0; i < 4; i++ {
		id := "item-" + string(rune('a'+i))
		if err := store.Save(sampleResponse(id)); err != nil {
			t.Fatalf("Save %s: %v", id, err)
		}
	}

	// "item-a" (oldest) should have been evicted.
	got, err := store.Get("item-a")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != nil {
		t.Error("oldest entry should have been evicted")
	}

	// "item-b", "item-c", "item-d" should still exist.
	for _, id := range []string{"item-b", "item-c", "item-d"} {
		got, err := store.Get(id)
		if err != nil {
			t.Fatalf("Get %s: %v", id, err)
		}
		if got == nil {
			t.Errorf("entry %s should still exist after eviction", id)
		}
	}
}

// ---- SQLiteStore Tests ----

func TestSQLite_SaveGetRoundtrip(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := storage.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer store.Close()

	resp := sampleResponse("sqlite-1")
	if err := store.Save(resp); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Get("sqlite-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got == nil {
		t.Fatal("Get returned nil")
	}
	if got.ID != "sqlite-1" {
		t.Errorf("ID mismatch: got %s, want sqlite-1", got.ID)
	}
	if got.FinalValue != 0.72 {
		t.Errorf("FinalValue mismatch: got %f, want 0.72", got.FinalValue)
	}
	if len(got.Breakdown) != 1 {
		t.Errorf("Breakdown length: got %d, want 1", len(got.Breakdown))
	}
}

func TestSQLite_Exists(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer store.Close()

	_ = store.Save(sampleResponse("sqlite-2"))

	exists, err := store.Exists("sqlite-2")
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

func TestSQLite_GetNonexistent(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer store.Close()

	got, err := store.Get("nope")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestSQLite_Persistence(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "persist.db")

	// Write with one store instance.
	store1, err := storage.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore (1): %v", err)
	}
	_ = store1.Save(sampleResponse("persist-1"))
	store1.Close()

	// Read with a new store instance to verify persistence.
	store2, err := storage.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore (2): %v", err)
	}
	defer store2.Close()

	got, err := store2.Get("persist-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got == nil {
		t.Fatal("data did not persist across store instances")
	}
	if got.ID != "persist-1" {
		t.Errorf("ID mismatch: got %s, want persist-1", got.ID)
	}
}

// ---- Factory Tests ----

func TestFactory_Memory(t *testing.T) {
	store, err := storage.NewStore("memory", "")
	if err != nil {
		t.Fatalf("NewStore(memory): %v", err)
	}
	if store == nil {
		t.Fatal("NewStore returned nil")
	}
}

func TestFactory_SQLite(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.NewStore("sqlite", filepath.Join(dir, "factory.db"))
	if err != nil {
		t.Fatalf("NewStore(sqlite): %v", err)
	}
	if store == nil {
		t.Fatal("NewStore returned nil")
	}
}

func TestFactory_Unknown(t *testing.T) {
	_, err := storage.NewStore("redis", "")
	if err == nil {
		t.Fatal("expected error for unknown backend")
	}
}

func TestFactory_SQLiteEmptyPath(t *testing.T) {
	_, err := storage.NewStore("sqlite", "")
	if err == nil {
		t.Fatal("expected error for empty sqlite path")
	}
}

func TestSQLite_FileCreated(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "created.db")

	store, err := storage.NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer store.Close()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("database file was not created")
	}
}
