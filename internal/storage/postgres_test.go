package storage

import (
	"os"
	"testing"
)

// TestPostgresStoreImplementsInterface is a compile-time check that
// PostgresStore satisfies the ExplanationStore interface.
func TestPostgresStoreImplementsInterface(t *testing.T) {
	var _ ExplanationStore = (*PostgresStore)(nil)
}

func TestBuildPostgresDSN_Defaults(t *testing.T) {
	// Clear any env vars that might interfere.
	envVars := []string{"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD", "DB_SSLMODE"}
	saved := make(map[string]string, len(envVars))
	for _, k := range envVars {
		saved[k] = os.Getenv(k)
		os.Unsetenv(k)
	}
	defer func() {
		for k, v := range saved {
			if v != "" {
				os.Setenv(k, v)
			}
		}
	}()

	dsn := BuildPostgresDSN()
	expected := "host=localhost port=5432 dbname=explainable_engine user=postgres password= sslmode=disable"
	if dsn != expected {
		t.Errorf("default DSN mismatch:\n got:  %s\n want: %s", dsn, expected)
	}
}

func TestBuildPostgresDSN_CustomEnv(t *testing.T) {
	envs := map[string]string{
		"DB_HOST":     "db.example.com",
		"DB_PORT":     "5433",
		"DB_NAME":     "mydb",
		"DB_USER":     "admin",
		"DB_PASSWORD": "secret",
		"DB_SSLMODE":  "require",
	}

	// Save and set.
	saved := make(map[string]string, len(envs))
	for k, v := range envs {
		saved[k] = os.Getenv(k)
		os.Setenv(k, v)
	}
	defer func() {
		for k, v := range saved {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}()

	dsn := BuildPostgresDSN()
	expected := "host=db.example.com port=5433 dbname=mydb user=admin password=secret sslmode=require"
	if dsn != expected {
		t.Errorf("custom DSN mismatch:\n got:  %s\n want: %s", dsn, expected)
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	const key = "TEST_EXPLAINABLE_ENGINE_ENV"
	os.Unsetenv(key)

	if got := getEnvOrDefault(key, "fallback"); got != "fallback" {
		t.Errorf("expected fallback, got %s", got)
	}

	os.Setenv(key, "custom")
	defer os.Unsetenv(key)

	if got := getEnvOrDefault(key, "fallback"); got != "custom" {
		t.Errorf("expected custom, got %s", got)
	}
}
