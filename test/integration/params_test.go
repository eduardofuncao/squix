package integration

import (
	"strings"
	"testing"
)

func TestRunWithNamedParam(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	_, _, exitCode := env.RunSquix("add", "user_by_id", "SELECT * FROM users WHERE id = :user_id")
	if exitCode != 0 {
		t.Fatalf("add failed, exit %d", exitCode)
	}

	stdout, _, exitCode := env.RunSquix("run", "user_by_id", "--user_id", "1", "-f", "csv")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "Alice") {
		t.Errorf("expected Alice:\n%s", stdout)
	}
}

func TestRunWithPositionalParam(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	_, _, exitCode := env.RunSquix("add", "user_by_id", "SELECT * FROM users WHERE id = :user_id")
	if exitCode != 0 {
		t.Fatalf("add failed, exit %d", exitCode)
	}

	stdout, _, exitCode := env.RunSquix("run", "user_by_id", "2", "-f", "csv")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "Bob") {
		t.Errorf("expected Bob:\n%s", stdout)
	}
}

func TestMissingRequiredParam(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	_, _, exitCode := env.RunSquix("add", "user_by_id", "SELECT * FROM users WHERE id = :user_id")
	if exitCode != 0 {
		t.Fatalf("add failed, exit %d", exitCode)
	}

	_, stderr, exitCode := env.RunSquix("run", "user_by_id", "-f", "csv")
	if exitCode == 0 {
		t.Fatal("expected non-zero exit code")
	}
	if !strings.Contains(stderr, "missing required parameters") {
		t.Errorf("expected missing params error, got %q", stderr)
	}
}
