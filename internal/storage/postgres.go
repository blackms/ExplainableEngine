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
	"github.com/blackms/ExplainableEngine/migrations"
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

	// Run database migrations.
	if err := Migrate(db, migrations.FS, MigrateUp); err != nil {
		db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return &PostgresStore{db: db}, nil
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

// pgQueryBuilder builds PostgreSQL queries with numbered $N parameters.
type pgQueryBuilder struct {
	where string
	args  []any
	n     int // next parameter number
}

func newPgQueryBuilder() *pgQueryBuilder {
	return &pgQueryBuilder{where: " WHERE 1=1", n: 1}
}

func (b *pgQueryBuilder) addFilter(clause string, arg any) {
	b.where += fmt.Sprintf(" AND "+clause, b.n)
	b.args = append(b.args, arg)
	b.n++
}

func (b *pgQueryBuilder) addParam(arg any) string {
	p := fmt.Sprintf("$%d", b.n)
	b.args = append(b.args, arg)
	b.n++
	return p
}

func pgApplyFilters(b *pgQueryBuilder, opts ListOptions) {
	if opts.Target != "" {
		b.addFilter("data->>'target' LIKE $%d", "%"+opts.Target+"%")
	}
	if opts.MinConfidence > 0 {
		b.addFilter("(data->>'confidence')::float >= $%d", opts.MinConfidence)
	}
	if opts.MaxConfidence > 0 {
		b.addFilter("(data->>'confidence')::float <= $%d", opts.MaxConfidence)
	}
	if opts.FromTime != "" {
		b.addFilter("created_at >= $%d", opts.FromTime)
	}
	if opts.ToTime != "" {
		b.addFilter("created_at <= $%d", opts.ToTime)
	}
}

// List returns a paginated, filtered list of explanations sorted by created_at descending.
func (s *PostgresStore) List(opts ListOptions) (*ListResult, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Count query.
	cb := newPgQueryBuilder()
	pgApplyFilters(cb, opts)
	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM explanations"+cb.where, cb.args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count explanations: %w", err)
	}

	// Data query (separate parameter numbering).
	qb := newPgQueryBuilder()
	pgApplyFilters(qb, opts)
	if opts.Cursor != "" {
		qb.addFilter("(created_at, id) < (SELECT created_at, id FROM explanations WHERE id = $%d)", opts.Cursor)
	}
	limitParam := qb.addParam(limit + 1)
	query := "SELECT data FROM explanations" + qb.where + " ORDER BY created_at DESC, id DESC LIMIT " + limitParam

	rows, err := s.db.Query(query, qb.args...)
	if err != nil {
		return nil, fmt.Errorf("list explanations: %w", err)
	}
	defer rows.Close()

	var items []*models.ExplainResponse
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("scan explanation: %w", err)
		}
		var resp models.ExplainResponse
		if err := json.Unmarshal(data, &resp); err != nil {
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
func (s *PostgresStore) Count() (int, error) {
	var count int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM explanations").Scan(&count); err != nil {
		return 0, fmt.Errorf("count explanations: %w", err)
	}
	return count, nil
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
