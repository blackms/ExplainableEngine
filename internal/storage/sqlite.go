package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/blackms/ExplainableEngine/internal/models"

	_ "modernc.org/sqlite"
)

// SQLiteStore persists explanations in a SQLite database.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore opens (or creates) a SQLite database at path and ensures the
// explanations table exists.
func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	// Enable WAL mode for better concurrent read performance.
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("set WAL mode: %w", err)
	}

	createSQL := `CREATE TABLE IF NOT EXISTS explanations (
		id         TEXT PRIMARY KEY,
		data       TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	if _, err := db.Exec(createSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("create table: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Save(resp *models.ExplainResponse) error {
	if resp == nil {
		return fmt.Errorf("cannot save nil response")
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}
	_, err = s.db.Exec(
		"INSERT OR REPLACE INTO explanations (id, data) VALUES (?, ?)",
		resp.ID, string(data),
	)
	if err != nil {
		return fmt.Errorf("insert explanation: %w", err)
	}
	return nil
}

func (s *SQLiteStore) Get(id string) (*models.ExplainResponse, error) {
	var data string
	err := s.db.QueryRow("SELECT data FROM explanations WHERE id = ?", id).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query explanation: %w", err)
	}

	var resp models.ExplainResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		return nil, fmt.Errorf("unmarshal explanation: %w", err)
	}
	return &resp, nil
}

func (s *SQLiteStore) Exists(id string) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(1) FROM explanations WHERE id = ?", id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check existence: %w", err)
	}
	return count > 0, nil
}

// List returns a paginated, filtered list of explanations sorted by created_at descending.
func (s *SQLiteStore) List(opts ListOptions) (*ListResult, error) {
	// Count total matching items (before pagination).
	countQuery := "SELECT COUNT(*) FROM explanations WHERE 1=1"
	dataQuery := "SELECT data FROM explanations WHERE 1=1"
	var args []any
	var countArgs []any

	if opts.Target != "" {
		clause := " AND json_extract(data, '$.target') LIKE ?"
		countQuery += clause
		dataQuery += clause
		v := "%" + opts.Target + "%"
		args = append(args, v)
		countArgs = append(countArgs, v)
	}
	if opts.MinConfidence > 0 {
		clause := " AND CAST(json_extract(data, '$.confidence') AS REAL) >= ?"
		countQuery += clause
		dataQuery += clause
		args = append(args, opts.MinConfidence)
		countArgs = append(countArgs, opts.MinConfidence)
	}
	if opts.MaxConfidence > 0 {
		clause := " AND CAST(json_extract(data, '$.confidence') AS REAL) <= ?"
		countQuery += clause
		dataQuery += clause
		args = append(args, opts.MaxConfidence)
		countArgs = append(countArgs, opts.MaxConfidence)
	}
	if opts.FromTime != "" {
		clause := " AND created_at >= ?"
		countQuery += clause
		dataQuery += clause
		args = append(args, opts.FromTime)
		countArgs = append(countArgs, opts.FromTime)
	}
	if opts.ToTime != "" {
		clause := " AND created_at <= ?"
		countQuery += clause
		dataQuery += clause
		args = append(args, opts.ToTime)
		countArgs = append(countArgs, opts.ToTime)
	}

	var total int
	if err := s.db.QueryRow(countQuery, countArgs...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count explanations: %w", err)
	}

	if opts.Cursor != "" {
		dataQuery += " AND (created_at, rowid) < (SELECT created_at, rowid FROM explanations WHERE id = ?)"
		args = append(args, opts.Cursor)
	}

	dataQuery += " ORDER BY created_at DESC, rowid DESC"

	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	dataQuery += " LIMIT ?"
	args = append(args, limit+1) // fetch one extra to detect next page

	rows, err := s.db.Query(dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("list explanations: %w", err)
	}
	defer rows.Close()

	var items []*models.ExplainResponse
	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("scan explanation: %w", err)
		}
		var resp models.ExplainResponse
		if err := json.Unmarshal([]byte(data), &resp); err != nil {
			return nil, fmt.Errorf("unmarshal explanation: %w", err)
		}
		items = append(items, &resp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate explanations: %w", err)
	}

	var nextCursor string
	if len(items) > limit {
		nextCursor = items[limit-1].ID
		items = items[:limit]
	}

	if items == nil {
		items = []*models.ExplainResponse{}
	}

	return &ListResult{
		Items:      items,
		NextCursor: nextCursor,
		Total:      total,
	}, nil
}

// Count returns the total number of stored explanations.
func (s *SQLiteStore) Count() (int, error) {
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM explanations").Scan(&count); err != nil {
		return 0, fmt.Errorf("count explanations: %w", err)
	}
	return count, nil
}

// Close closes the underlying database connection.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
