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

// Close closes the underlying database connection.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
