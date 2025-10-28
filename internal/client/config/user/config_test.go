package user

import (
	"os"
	"path/filepath"
	"testing"
)

// helper: returns a ConfigFile using a temp directory.
func newTempClient(t *testing.T) *ConfigFile {
	t.Helper()
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.yml")
	return &ConfigFile{Path: path}
}

func TestSetAndListAliases(t *testing.T) {
	c := newTempClient(t)

	// Initially empty
	aliases, err := c.ListAliases()
	if err != nil {
		t.Fatalf("ListAliases() error = %v", err)
	}
	if len(aliases) != 0 {
		t.Errorf("expected 0 aliases, got %d", len(aliases))
	}

	// Add an alias
	err = c.SetAlias("ls", "list --all")
	if err != nil {
		t.Fatalf("SetAlias() error = %v", err)
	}

	// Verify it was saved
	aliases, err = c.ListAliases()
	if err != nil {
		t.Fatalf("ListAliases() error = %v", err)
	}
	got, ok := aliases["ls"]
	if !ok || got != "list --all" {
		t.Errorf("alias mismatch: got=%q ok=%v, want %q", got, ok, "list --all")
	}

	// Update alias value
	err = c.SetAlias("ls", "list --verbose")
	if err != nil {
		t.Fatalf("SetAlias(update) error = %v", err)
	}
	aliases, _ = c.ListAliases()
	if aliases["ls"] != "list --verbose" {
		t.Errorf("alias not updated: got %q, want %q", aliases["ls"], "list --verbose")
	}
}

func TestRemoveAlias(t *testing.T) {
	c := newTempClient(t)

	// Create two aliases
	_ = c.SetAlias("a", "alpha")
	_ = c.SetAlias("b", "beta")

	err := c.RemoveAlias("a")
	if err != nil {
		t.Fatalf("RemoveAlias() error = %v", err)
	}

	aliases, err := c.ListAliases()
	if err != nil {
		t.Fatalf("ListAliases() error = %v", err)
	}

	if _, ok := aliases["a"]; ok {
		t.Errorf("expected alias 'a' removed, still present")
	}
	if _, ok := aliases["b"]; !ok {
		t.Errorf("expected alias 'b' still present")
	}
}

func TestLoadFeatureFlags_Defaults(t *testing.T) {
	c := newTempClient(t)

	flags, err := c.LoadFeatureFlags()
	if err != nil {
		t.Fatalf("LoadFeatureFlags() error = %v", err)
	}
	if flags.Trace {
		t.Errorf("expected Trace = false by default, got true")
	}
}

func TestPersistenceAcrossLoad(t *testing.T) {
	c := newTempClient(t)

	err := c.SetAlias("x", "execute")
	if err != nil {
		t.Fatalf("SetAlias() error = %v", err)
	}

	// Simulate a new client loading same file
	c2 := &ConfigFile{Path: c.Path}

	aliases, err := c2.ListAliases()
	if err != nil {
		t.Fatalf("ListAliases() error = %v", err)
	}

	if aliases["x"] != "execute" {
		t.Errorf("alias not persisted: got %q, want %q", aliases["x"], "execute")
	}

	// File should exist
	if _, err := os.Stat(c.Path); err != nil {
		t.Fatalf("expected config file to exist, got error: %v", err)
	}
}
