package integration

import (
	"strings"
	"testing"
)

func TestAddQuery(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("add", "list_users", "SELECT * FROM users")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "Added query") {
		t.Errorf("expected add confirmation, got %q", stdout)
	}

	stdout, _, _ = env.RunSquix("list", "queries")
	if !strings.Contains(stdout, "list_users") {
		t.Errorf("query not in list:\n%s", stdout)
	}
}

func TestListQueries(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	env.RunSquix("add", "q1", "SELECT 1")
	env.RunSquix("add", "q2", "SELECT 2")

	stdout, _, exitCode := env.RunSquix("list", "queries")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "q1") || !strings.Contains(stdout, "q2") {
		t.Errorf("expected both queries:\n%s", stdout)
	}
}

func TestListConnections(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("list", "connections")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "testdb") {
		t.Errorf("expected testdb connection:\n%s", stdout)
	}
}

func TestRemoveQuery(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	env.RunSquix("add", "temp_q", "SELECT 1")

	stdout, _, exitCode := env.RunSquix("remove", "temp_q")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "Removed") {
		t.Errorf("expected removal confirmation, got %q", stdout)
	}

	stdout, _, _ = env.RunSquix("list", "queries")
	if strings.Contains(stdout, "temp_q") {
		t.Errorf("query still in list after remove:\n%s", stdout)
	}
}
