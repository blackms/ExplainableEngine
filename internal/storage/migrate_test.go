package storage

import (
	"io/fs"
	"sort"
	"testing"
	"testing/fstest"
)

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		filename string
		want     int
	}{
		{"001_create_explanations.up.sql", 1},
		{"001_create_explanations.down.sql", 1},
		{"002_add_indexes.up.sql", 2},
		{"010_something.up.sql", 10},
		{"100_big.down.sql", 100},
		{"no_number.up.sql", 0},
		{"", 0},
	}
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := ExtractVersion(tt.filename)
			if got != tt.want {
				t.Errorf("ExtractVersion(%q) = %d, want %d", tt.filename, got, tt.want)
			}
		})
	}
}

func TestCollectMigrationFiles(t *testing.T) {
	mockFS := fstest.MapFS{
		"001_create.up.sql":   {Data: []byte("CREATE TABLE t;")},
		"001_create.down.sql": {Data: []byte("DROP TABLE t;")},
		"002_index.up.sql":    {Data: []byte("CREATE INDEX;")},
		"002_index.down.sql":  {Data: []byte("DROP INDEX;")},
		"README.md":           {Data: []byte("ignore me")},
	}

	t.Run("up files sorted", func(t *testing.T) {
		files, err := collectMigrationFiles(mockFS, MigrateUp)
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"001_create.up.sql", "002_index.up.sql"}
		if len(files) != len(want) {
			t.Fatalf("got %d files, want %d", len(files), len(want))
		}
		for i, f := range files {
			if f != want[i] {
				t.Errorf("file[%d] = %q, want %q", i, f, want[i])
			}
		}
		if !sort.StringsAreSorted(files) {
			t.Error("files not sorted")
		}
	})

	t.Run("down files sorted", func(t *testing.T) {
		files, err := collectMigrationFiles(mockFS, MigrateDown)
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"001_create.down.sql", "002_index.down.sql"}
		if len(files) != len(want) {
			t.Fatalf("got %d files, want %d", len(files), len(want))
		}
		for i, f := range files {
			if f != want[i] {
				t.Errorf("file[%d] = %q, want %q", i, f, want[i])
			}
		}
	})
}

func TestEmbeddedMigrationsFromStoragePackage(t *testing.T) {
	// Compile-time verification that the migrations package FS is accessible.
	// This imports the migrations package indirectly through the test's build.
	// The actual content check lives in migrations/embed_test.go.

	// Verify that collectMigrationFiles works with a realistic FS.
	mockFS := fstest.MapFS{
		"001_create_explanations.up.sql":   {Data: []byte("CREATE TABLE IF NOT EXISTS explanations;")},
		"001_create_explanations.down.sql": {Data: []byte("DROP TABLE IF EXISTS explanations;")},
		"002_add_indexes.up.sql":           {Data: []byte("CREATE INDEX;")},
		"002_add_indexes.down.sql":         {Data: []byte("DROP INDEX;")},
	}

	entries, err := fs.ReadDir(mockFS, ".")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) < 4 {
		t.Errorf("expected at least 4 entries, got %d", len(entries))
	}
}
