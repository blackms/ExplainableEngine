package migrations

import (
	"io/fs"
	"testing"
)

func TestEmbeddedMigrations(t *testing.T) {
	entries, err := fs.ReadDir(FS, ".")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) < 4 {
		t.Errorf("expected at least 4 migration files, got %d", len(entries))
	}

	// Verify expected files are present.
	expected := map[string]bool{
		"001_create_explanations.up.sql":   false,
		"001_create_explanations.down.sql": false,
		"002_add_indexes.up.sql":           false,
		"002_add_indexes.down.sql":         false,
	}
	for _, e := range entries {
		if _, ok := expected[e.Name()]; ok {
			expected[e.Name()] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("expected migration file %q not found in embedded FS", name)
		}
	}
}
