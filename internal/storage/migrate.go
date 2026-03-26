package storage

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"
)

// MigrationDirection indicates whether to apply or roll back migrations.
type MigrationDirection string

const (
	// MigrateUp applies all pending migrations.
	MigrateUp MigrationDirection = "up"
	// MigrateDown rolls back applied migrations in reverse order.
	MigrateDown MigrationDirection = "down"
)

// Migrate runs all pending migrations in the given direction using the
// provided embedded filesystem. It creates a schema_migrations tracking
// table automatically.
func Migrate(db *sql.DB, migrations fs.FS, direction MigrationDirection) error {
	if err := ensureMigrationsTable(db); err != nil {
		return err
	}

	currentVersion, err := currentMigrationVersion(db)
	if err != nil {
		return err
	}

	files, err := collectMigrationFiles(migrations, direction)
	if err != nil {
		return err
	}

	if direction == MigrateUp {
		return applyUp(db, migrations, files, currentVersion)
	}
	return applyDown(db, migrations, files, currentVersion)
}

// RunMigrations opens a database connection, runs migrations, and closes it.
// Useful for standalone migration tooling.
func RunMigrations(dsn string, migrations fs.FS, direction MigrationDirection) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("opening postgres for migration: %w", err)
	}
	defer db.Close()

	return Migrate(db, migrations, direction)
}

// ensureMigrationsTable creates the schema_migrations tracking table if it
// does not already exist.
func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    INTEGER PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("creating schema_migrations table: %w", err)
	}
	return nil
}

// currentMigrationVersion returns the highest applied migration version, or 0
// if no migrations have been applied.
func currentMigrationVersion(db *sql.DB) (int, error) {
	var v int
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&v)
	if err != nil {
		return 0, fmt.Errorf("reading current migration version: %w", err)
	}
	return v, nil
}

// collectMigrationFiles returns sorted filenames matching the given direction.
func collectMigrationFiles(migrations fs.FS, direction MigrationDirection) ([]string, error) {
	suffix := "." + string(direction) + ".sql"

	entries, err := fs.ReadDir(migrations, ".")
	if err != nil {
		return nil, fmt.Errorf("reading migration files: %w", err)
	}

	var files []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), suffix) {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	return files, nil
}

// applyUp runs each migration whose version exceeds currentVersion.
func applyUp(db *sql.DB, migrations fs.FS, files []string, currentVersion int) error {
	for _, f := range files {
		version := ExtractVersion(f)
		if version <= currentVersion {
			continue
		}

		content, err := fs.ReadFile(migrations, f)
		if err != nil {
			return fmt.Errorf("reading migration file %s: %w", f, err)
		}

		log.Printf("Applying migration %03d (up): %s", version, f)
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("migration %d failed: %w", version, err)
		}

		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			return fmt.Errorf("recording migration %d: %w", version, err)
		}
	}
	return nil
}

// applyDown rolls back each migration whose version is <= currentVersion, in
// reverse order.
func applyDown(db *sql.DB, migrations fs.FS, files []string, currentVersion int) error {
	for i := len(files) - 1; i >= 0; i-- {
		f := files[i]
		version := ExtractVersion(f)
		if version > currentVersion {
			continue
		}

		content, err := fs.ReadFile(migrations, f)
		if err != nil {
			return fmt.Errorf("reading migration file %s: %w", f, err)
		}

		log.Printf("Applying migration %03d (down): %s", version, f)
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("migration %d rollback failed: %w", version, err)
		}

		if _, err := db.Exec("DELETE FROM schema_migrations WHERE version = $1", version); err != nil {
			return fmt.Errorf("removing migration %d record: %w", version, err)
		}
	}
	return nil
}

// ExtractVersion parses the leading numeric prefix from a migration filename.
// For example, "001_create_explanations.up.sql" returns 1.
func ExtractVersion(filename string) int {
	parts := strings.SplitN(filename, "_", 2)
	if len(parts) == 0 {
		return 0
	}
	var v int
	fmt.Sscanf(parts[0], "%d", &v)
	return v
}
