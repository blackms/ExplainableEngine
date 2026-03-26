package storage_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/blackms/ExplainableEngine/internal/storage"
)

func TestSQLite_Count(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer store.Close()

	count, err := store.Count()
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}

	now := time.Now()
	_ = store.Save(sampleResponseWithOpts("sc-1", "a", 0.8, now))
	_ = store.Save(sampleResponseWithOpts("sc-2", "b", 0.9, now.Add(time.Second)))

	count, err = store.Count()
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}

func TestSQLite_ListBasic(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer store.Close()

	now := time.Now()
	_ = store.Save(sampleResponseWithOpts("sl-1", "score", 0.8, now.Add(-2*time.Second)))
	_ = store.Save(sampleResponseWithOpts("sl-2", "rank", 0.9, now.Add(-1*time.Second)))
	_ = store.Save(sampleResponseWithOpts("sl-3", "score", 0.7, now))

	result, err := store.List(storage.ListOptions{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if result.Total != 3 {
		t.Errorf("total: expected 3, got %d", result.Total)
	}
	if len(result.Items) != 3 {
		t.Errorf("items: expected 3, got %d", len(result.Items))
	}
}

func TestSQLite_ListTargetFilter(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer store.Close()

	now := time.Now()
	_ = store.Save(sampleResponseWithOpts("stf-1", "user_score", 0.8, now.Add(-2*time.Second)))
	_ = store.Save(sampleResponseWithOpts("stf-2", "rank_score", 0.9, now.Add(-1*time.Second)))
	_ = store.Save(sampleResponseWithOpts("stf-3", "category", 0.7, now))

	result, err := store.List(storage.ListOptions{Target: "score"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(result.Items) != 2 {
		t.Errorf("items: expected 2, got %d", len(result.Items))
	}
}

func TestSQLite_ListConfidenceFilter(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer store.Close()

	now := time.Now()
	_ = store.Save(sampleResponseWithOpts("scf-1", "a", 0.5, now.Add(-2*time.Second)))
	_ = store.Save(sampleResponseWithOpts("scf-2", "b", 0.8, now.Add(-1*time.Second)))
	_ = store.Save(sampleResponseWithOpts("scf-3", "c", 0.95, now))

	result, err := store.List(storage.ListOptions{MinConfidence: 0.7})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(result.Items) != 2 {
		t.Errorf("items: expected 2 (>= 0.7), got %d", len(result.Items))
	}
}

func TestSQLite_ListPagination(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer store.Close()

	now := time.Now()
	for i := 0; i < 5; i++ {
		id := "sp-" + string(rune('a'+i))
		_ = store.Save(sampleResponseWithOpts(id, "score", 0.8, now.Add(time.Duration(i)*time.Second)))
	}

	page1, err := store.List(storage.ListOptions{Limit: 2})
	if err != nil {
		t.Fatalf("List page1: %v", err)
	}
	if len(page1.Items) != 2 {
		t.Fatalf("page1 items: expected 2, got %d", len(page1.Items))
	}
	if page1.NextCursor == "" {
		t.Fatal("page1 should have next_cursor")
	}

	page2, err := store.List(storage.ListOptions{Limit: 2, Cursor: page1.NextCursor})
	if err != nil {
		t.Fatalf("List page2: %v", err)
	}
	if len(page2.Items) != 2 {
		t.Fatalf("page2 items: expected 2, got %d", len(page2.Items))
	}

	// No overlap.
	for _, p1 := range page1.Items {
		for _, p2 := range page2.Items {
			if p1.ID == p2.ID {
				t.Errorf("overlap between pages: %s", p1.ID)
			}
		}
	}
}

func TestSQLite_ListEmpty(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	defer store.Close()

	result, err := store.List(storage.ListOptions{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(result.Items) != 0 {
		t.Errorf("items: expected 0, got %d", len(result.Items))
	}
	if result.Total != 0 {
		t.Errorf("total: expected 0, got %d", result.Total)
	}
}
