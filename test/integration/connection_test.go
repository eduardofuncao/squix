package integration

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestInitConnection(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	db2Path := filepath.Join(t.TempDir(), "test2.db")
	stdout, _, exitCode := env.RunSquix("init", "mydb2", db2Path)
	if exitCode != 0 {
		t.Fatalf("exit %d, stdout: %s", exitCode, stdout)
	}
	if !strings.Contains(stdout, "Connection created") {
		t.Errorf("expected creation confirmation, got %q", stdout)
	}

	stdout, _, _ = env.RunSquix("list", "connections")
	if !strings.Contains(stdout, "testdb") || !strings.Contains(stdout, "mydb2") {
		t.Errorf("expected both connections:\n%s", stdout)
	}
}

func TestSwitchConnection(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	db2Path := filepath.Join(t.TempDir(), "test2.db")
	env.RunSquix("init", "mydb2", db2Path)

	stdout, _, exitCode := env.RunSquix("switch", "testdb")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "Switched to") {
		t.Errorf("expected switch confirmation, got %q", stdout)
	}

	stdout, _, _ = env.RunSquix("list", "connections")
	if !strings.Contains(stdout, "testdb") {
		t.Errorf("expected testdb active:\n%s", stdout)
	}
}

func TestRemoveConnection(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	db2Path := filepath.Join(t.TempDir(), "test2.db")
	env.RunSquix("init", "mydb2", db2Path)

	stdout, _, exitCode := env.RunSquixWithStdin("y\n", "remove", "-c", "mydb2")
	if exitCode != 0 {
		t.Fatalf("exit %d, stdout: %s", exitCode, stdout)
	}
	if !strings.Contains(stdout, "Removed connection") {
		t.Errorf("expected removal confirmation, got %q", stdout)
	}

	stdout, _, _ = env.RunSquix("list", "connections")
	if strings.Contains(stdout, "mydb2") {
		t.Errorf("mydb2 still in list after remove:\n%s", stdout)
	}
}
