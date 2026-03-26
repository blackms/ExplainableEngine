package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/blackms/ExplainableEngine/internal/models"
)

// PostgresStore persists explanations in a PostgreSQL database using JSONB.
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore opens a connection to PostgreSQL at dsn, configures the
// connection pool, verifies connectivity, and ensures the explanations table exists.
func NewPostgresStore(dsn string) (*PostgresStore, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening postgres: %w", err)
	}

	// Connection pool settings.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("pinging postgres: %w", err)
	}

	// Auto-create table.
	if err := createPostgresTable(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("creating table: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

func createPostgresTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS explanations (
			id         TEXT PRIMARY KEY,
			data       JSONB NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	return err
}

// Save persists an ExplainResponse. Uses INSERT ... ON CONFLICT DO NOTHING
// for idempotent writes.
func (s *PostgresStore) Save(resp *models.ExplainResponse) error {
	if resp == nil {
		return fmt.Errorf("cannot save nil response")
	}
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}
	_, err = s.db.Exec(
		"INSERT INTO explanations (id, data) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING",
		resp.ID, data,
	)
	if err != nil {
		return fmt.Errorf("insert explanation: %w", err)
	}
	return nil
}

// Get retrieves an ExplainResponse by ID. Returns nil, nil if not found.
func (s *PostgresStore) Get(id string) (*models.ExplainResponse, error) {
	var data []byte
	err := s.db.QueryRow("SELECT data FROM explanations WHERE id = $1", id).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query explanation: %w", err)
	}

	var resp models.ExplainResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal explanation: %w", err)
	}
	return &resp, nil
}

// Exists checks whether an explanation with the given ID is stored.
func (s *PostgresStore) Exists(id string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM explanations WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check existence: %w", err)
	}
	return exists, nil
}

// Close closes the underlying database connection pool.
func (s *PostgresStore) Close() error {
	return s.db.Close()
}

// BuildPostgresDSN constructs a PostgreSQL connection string from environment
// variables, falling back to sensible defaults for local development.
func BuildPostgresDSN() string {
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	name := getEnvOrDefault("DB_NAME", "explainable_engine")
	user := getEnvOrDefault("DB_USER", "postgres")
	pass := getEnvOrDefault("DB_PASSWORD", "")
	sslmode := getEnvOrDefault("DB_SSLMODE", "disable")

	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		host, port, name, user, pass, sslmode)
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
