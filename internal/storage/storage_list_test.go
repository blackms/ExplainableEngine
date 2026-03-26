package storage_test

import (
	"testing"
	"time"

	"github.com/blackms/ExplainableEngine/internal/models"
	"github.com/blackms/ExplainableEngine/internal/storage"
)

func sampleResponseWithOpts(id, target string, confidence float64, createdAt time.Time) *models.ExplainResponse {
	return &models.ExplainResponse{
		ID:         id,
		Target:     target,
		FinalValue: 0.72,
		Confidence: confidence,
		Breakdown: []models.BreakdownItem{
			{NodeID: "a", Label: "A", Value: 0.8, Weight: 0.4},
		},
		TopDrivers: []models.DriverItem{
			{Name: "A", Impact: 0.44, Rank: 1},
		},
		Metadata: models.ExplainMetadata{
			Version:   "0.1.0",
			CreatedAt: createdAt,
		},
	}
}

func TestInMemory_Count(t *testing.T) {
	store := storage.NewInMemoryStore()

	count, err := store.Count()
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}

	_ = store.Save(sampleResponse("c-1"))
	_ = store.Save(sampleResponse("c-2"))

	count, err = store.Count()
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}

func TestInMemory_ListBasic(t *testing.T) {
	store := storage.NewInMemoryStore()
	now := time.Now()

	_ = store.Save(sampleResponseWithOpts("l-1", "score", 0.8, now.Add(-3*time.Second)))
	_ = store.Save(sampleResponseWithOpts("l-2", "score", 0.9, now.Add(-2*time.Second)))
	_ = store.Save(sampleResponseWithOpts("l-3", "rank", 0.7, now.Add(-1*time.Second)))

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
	// Should be newest first.
	if result.Items[0].ID != "l-3" {
		t.Errorf("first item: expected l-3, got %s", result.Items[0].ID)
	}
}

func TestInMemory_ListWithLimit(t *testing.T) {
	store := storage.NewInMemoryStore()
	now := time.Now()

	for i := 0; i < 5; i++ {
		id := "lim-" + string(rune('a'+i))
		_ = store.Save(sampleResponseWithOpts(id, "score", 0.8, now.Add(time.Duration(i)*time.Second)))
	}

	result, err := store.List(storage.ListOptions{Limit: 2})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(result.Items) != 2 {
		t.Errorf("items: expected 2, got %d", len(result.Items))
	}
	if result.NextCursor == "" {
		t.Error("expected next_cursor to be set")
	}
	if result.Total != 5 {
		t.Errorf("total: expected 5, got %d", result.Total)
	}
}

func TestInMemory_ListCursorPagination(t *testing.T) {
	store := storage.NewInMemoryStore()
	now := time.Now()

	for i := 0; i < 5; i++ {
		id := "pg-" + string(rune('a'+i))
		_ = store.Save(sampleResponseWithOpts(id, "score", 0.8, now.Add(time.Duration(i)*time.Second)))
	}

	// Get page 1.
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

	// Get page 2 using cursor.
	page2, err := store.List(storage.ListOptions{Limit: 2, Cursor: page1.NextCursor})
	if err != nil {
		t.Fatalf("List page2: %v", err)
	}
	if len(page2.Items) != 2 {
		t.Fatalf("page2 items: expected 2, got %d", len(page2.Items))
	}

	// Ensure no overlap.
	for _, p1 := range page1.Items {
		for _, p2 := range page2.Items {
			if p1.ID == p2.ID {
				t.Errorf("overlap between pages: %s", p1.ID)
			}
		}
	}
}

func TestInMemory_ListTargetFilter(t *testing.T) {
	store := storage.NewInMemoryStore()
	now := time.Now()

	_ = store.Save(sampleResponseWithOpts("tf-1", "user_score", 0.8, now.Add(-2*time.Second)))
	_ = store.Save(sampleResponseWithOpts("tf-2", "rank_score", 0.9, now.Add(-1*time.Second)))
	_ = store.Save(sampleResponseWithOpts("tf-3", "category", 0.7, now))

	result, err := store.List(storage.ListOptions{Target: "score"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(result.Items) != 2 {
		t.Errorf("items: expected 2 (score matches), got %d", len(result.Items))
	}
}

func TestInMemory_ListConfidenceFilter(t *testing.T) {
	store := storage.NewInMemoryStore()
	now := time.Now()

	_ = store.Save(sampleResponseWithOpts("cf-1", "a", 0.5, now.Add(-2*time.Second)))
	_ = store.Save(sampleResponseWithOpts("cf-2", "b", 0.8, now.Add(-1*time.Second)))
	_ = store.Save(sampleResponseWithOpts("cf-3", "c", 0.95, now))

	result, err := store.List(storage.ListOptions{MinConfidence: 0.7})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(result.Items) != 2 {
		t.Errorf("items: expected 2 (>= 0.7), got %d", len(result.Items))
	}

	result, err = store.List(storage.ListOptions{MaxConfidence: 0.8})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(result.Items) != 2 {
		t.Errorf("items: expected 2 (<= 0.8), got %d", len(result.Items))
	}
}

func TestInMemory_ListEmpty(t *testing.T) {
	store := storage.NewInMemoryStore()

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
