package config

import (
	"os"
	"path/filepath"
	"testing"
)

// withTempCfg swaps CfgPath to a temp dir for the test and restores it after.
func withTempCfg(t *testing.T) {
	t.Helper()
	orig := CfgPath
	CfgPath = t.TempDir()
	t.Cleanup(func() { CfgPath = orig })
}

func TestLoadLastQueryMissingFile(t *testing.T) {
	withTempCfg(t)
	got, err := LoadLastQuery()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	withTempCfg(t)
	want := "SELECT 1 FROM users"
	if err := SaveLastQuery(want); err != nil {
		t.Fatalf("SaveLastQuery: %v", err)
	}
	got, err := LoadLastQuery()
	if err != nil {
		t.Fatalf("LoadLastQuery: %v", err)
	}
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestSaveLastQueryTrims(t *testing.T) {
	withTempCfg(t)
	if err := SaveLastQuery("  SELECT 1  "); err != nil {
		t.Fatalf("SaveLastQuery: %v", err)
	}
	got, err := LoadLastQuery()
	if err != nil {
		t.Fatalf("LoadLastQuery: %v", err)
	}
	if got != "SELECT 1" {
		t.Fatalf("expected trimmed %q, got %q", "SELECT 1", got)
	}
}

func TestSaveLastQueryEmptyNoOp(t *testing.T) {
	withTempCfg(t)
	// First write something
	if err := SaveLastQuery("SELECT 1"); err != nil {
		t.Fatalf("SaveLastQuery: %v", err)
	}
	// Empty save should not clobber
	if err := SaveLastQuery(""); err != nil {
		t.Fatalf("SaveLastQuery empty: %v", err)
	}
	got, _ := LoadLastQuery()
	if got != "SELECT 1" {
		t.Fatalf("expected prior value preserved, got %q", got)
	}
	// Whitespace-only is also a no-op: prior value preserved.
	if err := SaveLastQuery("   \n  "); err != nil {
		t.Fatalf("SaveLastQuery whitespace: %v", err)
	}
	got, _ = LoadLastQuery()
	if got != "SELECT 1" {
		t.Fatalf("expected prior value preserved after whitespace save, got %q", got)
	}
}

func TestClearLastQuery(t *testing.T) {
	withTempCfg(t)
	if err := SaveLastQuery("SELECT 1"); err != nil {
		t.Fatalf("SaveLastQuery: %v", err)
	}
	if err := ClearLastQuery(); err != nil {
		t.Fatalf("ClearLastQuery: %v", err)
	}
	if _, err := os.Stat(LastQueryPath()); !os.IsNotExist(err) {
		t.Fatalf("expected file removed, stat err=%v", err)
	}
}

func TestClearLastQueryMissingNoError(t *testing.T) {
	withTempCfg(t)
	if err := ClearLastQuery(); err != nil {
		t.Fatalf("ClearLastQuery on missing file: %v", err)
	}
}

func TestSaveLastQueryCreatesCfgDir(t *testing.T) {
	orig := CfgPath
	CfgPath = filepath.Join(t.TempDir(), "nested", "squix")
	t.Cleanup(func() { CfgPath = orig })

	if err := SaveLastQuery("SELECT 1"); err != nil {
		t.Fatalf("SaveLastQuery: %v", err)
	}
	got, _ := LoadLastQuery()
	if got != "SELECT 1" {
		t.Fatalf("expected round-trip, got %q", got)
	}
}

func TestSaveLastQueryPerms(t *testing.T) {
	withTempCfg(t)
	if err := SaveLastQuery("SELECT 1"); err != nil {
		t.Fatalf("SaveLastQuery: %v", err)
	}
	info, err := os.Stat(LastQueryPath())
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected 0600 perms, got %v", info.Mode().Perm())
	}
}
