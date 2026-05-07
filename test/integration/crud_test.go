package integration

import (
	"strings"
	"testing"
)

func TestInsertStatement(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	_, stderr, exitCode := env.RunSquix("run", "INSERT INTO users (name, email, age) VALUES ('Diana', 'diana@example.com', 28)", "-f", "json")
	if exitCode != 0 {
		t.Fatalf("exit %d, stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stderr, "successfully") {
		t.Errorf("expected success on stderr, got %q", stderr)
	}

	stdout, _, _ := env.RunSquix("run", "SELECT name FROM users WHERE email = 'diana@example.com'", "-f", "csv")
	if !strings.Contains(stdout, "Diana") {
		t.Errorf("Diana not found after insert:\n%s", stdout)
	}
}

func TestUpdateStatement(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	_, _, exitCode := env.RunSquix("run", "UPDATE users SET age = 31 WHERE name = 'Alice'", "-f", "json")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}

	stdout, _, _ := env.RunSquix("run", "SELECT age FROM users WHERE name = 'Alice'", "-f", "csv")
	if !strings.Contains(stdout, "31") {
		t.Errorf("age not updated:\n%s", stdout)
	}
}

func TestDeleteStatement(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	_, _, exitCode := env.RunSquix("run", "DELETE FROM orders WHERE user_id = 2", "-f", "json")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}

	stdout, _, _ := env.RunSquix("run", "SELECT COUNT(*) FROM orders WHERE user_id = 2", "-f", "csv")
	if !strings.Contains(stdout, "0") {
		t.Errorf("orders not deleted:\n%s", stdout)
	}
}

func TestNULLValues(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	_, _, exitCode := env.RunSquix("run", "INSERT INTO users (name, email) VALUES ('Eve', 'eve@example.com')", "-f", "json")
	if exitCode != 0 {
		t.Fatalf("insert exit %d", exitCode)
	}

	stdout, _, exitCode := env.RunSquix("run", "SELECT name, age FROM users WHERE name = 'Eve'", "-f", "json")
	if exitCode != 0 {
		t.Fatalf("select exit %d", exitCode)
	}
	if !strings.Contains(stdout, "Eve") {
		t.Errorf("Eve not found:\n%s", stdout)
	}
}

func TestSpecialCharacters(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	_, _, exitCode := env.RunSquix("run", "INSERT INTO users (name, email, age) VALUES ('O''Brien', 'ob@test.com', 40)", "-f", "json")
	if exitCode != 0 {
		t.Fatalf("insert exit %d", exitCode)
	}

	stdout, _, exitCode := env.RunSquix("run", "SELECT name, email FROM users WHERE email = 'ob@test.com'", "-f", "csv")
	if exitCode != 0 {
		t.Fatalf("select exit %d", exitCode)
	}
	if !strings.Contains(stdout, "O'Brien") {
		t.Errorf("O'Brien not found:\n%s", stdout)
	}
}
