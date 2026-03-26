package storage

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/blackms/ExplainableEngine/internal/models"
)

const defaultMaxSize = 1000

// InMemoryStore is a thread-safe, LRU-bounded in-memory implementation of ExplanationStore.
type InMemoryStore struct {
	mu      sync.RWMutex
	data    map[string]*models.ExplainResponse
	order   []string // oldest first (index 0 = oldest)
	maxSize int
}

// InMemoryOption configures an InMemoryStore.
type InMemoryOption func(*InMemoryStore)

// WithMaxSize sets the maximum number of entries before LRU eviction kicks in.
func WithMaxSize(n int) InMemoryOption {
	return func(s *InMemoryStore) {
		if n > 0 {
			s.maxSize = n
		}
	}
}

// NewInMemoryStore creates a new InMemoryStore with the given options.
func NewInMemoryStore(opts ...InMemoryOption) *InMemoryStore {
	s := &InMemoryStore{
		data:    make(map[string]*models.ExplainResponse),
		order:   make([]string, 0),
		maxSize: defaultMaxSize,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *InMemoryStore) Save(resp *models.ExplainResponse) error {
	if resp == nil {
		return fmt.Errorf("cannot save nil response")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	// If key already exists, remove it from order so we can re-append (refresh).
	if _, exists := s.data[resp.ID]; exists {
		s.removeFromOrder(resp.ID)
	}

	// Evict oldest if at capacity.
	if len(s.data) >= s.maxSize {
		s.evictOldest()
	}

	s.data[resp.ID] = resp
	s.order = append(s.order, resp.ID)
	return nil
}

func (s *InMemoryStore) Get(id string) (*models.ExplainResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	resp, ok := s.data[id]
	if !ok {
		return nil, nil
	}
	return resp, nil
}

func (s *InMemoryStore) Exists(id string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[id]
	return ok, nil
}

// evictOldest removes the oldest entry. Caller must hold the write lock.
func (s *InMemoryStore) evictOldest() {
	if len(s.order) == 0 {
		return
	}
	oldest := s.order[0]
	s.order = s.order[1:]
	delete(s.data, oldest)
}

// removeFromOrder removes a key from the order slice. Caller must hold the write lock.
func (s *InMemoryStore) removeFromOrder(key string) {
	for i, k := range s.order {
		if k == key {
			s.order = append(s.order[:i], s.order[i+1:]...)
			return
		}
	}
}

// List returns a paginated, filtered list of explanations sorted by created_at descending.
func (s *InMemoryStore) List(opts ListOptions) (*ListResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Collect all items sorted by created_at descending.
	all := make([]*models.ExplainResponse, 0, len(s.data))
	for _, resp := range s.data {
		all = append(all, resp)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].Metadata.CreatedAt.After(all[j].Metadata.CreatedAt)
	})

	// Apply filters.
	var filtered []*models.ExplainResponse
	for _, resp := range all {
		if opts.Target != "" && !strings.Contains(resp.Target, opts.Target) {
			continue
		}
		if opts.MinConfidence > 0 && resp.Confidence < opts.MinConfidence {
			continue
		}
		if opts.MaxConfidence > 0 && resp.Confidence > opts.MaxConfidence {
			continue
		}
		if opts.FromTime != "" {
			fromT, err := time.Parse(time.RFC3339, opts.FromTime)
			if err == nil && resp.Metadata.CreatedAt.Before(fromT) {
				continue
			}
		}
		if opts.ToTime != "" {
			toT, err := time.Parse(time.RFC3339, opts.ToTime)
			if err == nil && resp.Metadata.CreatedAt.After(toT) {
				continue
			}
		}
		filtered = append(filtered, resp)
	}

	total := len(filtered)

	// Apply cursor: skip items until we pass the cursor ID.
	if opts.Cursor != "" {
		found := false
		for i, resp := range filtered {
			if resp.ID == opts.Cursor {
				filtered = filtered[i+1:]
				found = true
				break
			}
		}
		if !found {
			filtered = nil
		}
	}

	// Apply limit.
	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var nextCursor string
	if len(filtered) > limit {
		nextCursor = filtered[limit-1].ID
		filtered = filtered[:limit]
	}

	if filtered == nil {
		filtered = []*models.ExplainResponse{}
	}

	return &ListResult{
		Items:      filtered,
		NextCursor: nextCursor,
		Total:      total,
	}, nil
}

// Count returns the total number of stored explanations.
func (s *InMemoryStore) Count() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data), nil
}
