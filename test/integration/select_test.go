package integration

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSelectJSON(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("run", "SELECT id, name FROM users ORDER BY id", "-f", "json")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}

	var results []map[string]string
	if err := json.Unmarshal([]byte(stdout), &results); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, stdout)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(results))
	}
	if results[0]["name"] != "Alice" {
		t.Errorf("first row = %v", results[0])
	}
}

func TestSelectCSV(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("run", "SELECT id, name FROM users ORDER BY id", "-f", "csv")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}

	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines (header + 3 rows), got %d", len(lines))
	}
	if lines[0] != "id,name" {
		t.Errorf("header = %q", lines[0])
	}
	if !strings.Contains(lines[1], "Alice") {
		t.Errorf("first row = %q", lines[1])
	}
}

func TestSelectTSV(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("run", "SELECT id, name FROM users ORDER BY id", "-f", "tsv")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}

	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "\t") {
		t.Error("expected tabs in header")
	}
}

func TestSelectMarkdown(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("run", "SELECT id, name FROM users ORDER BY id", "-f", "markdown")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}

	if !strings.Contains(stdout, "|") || !strings.Contains(stdout, "---") {
		t.Errorf("expected markdown table:\n%s", stdout)
	}
}

func TestSelectHTML(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("run", "SELECT id, name FROM users ORDER BY id", "-f", "html")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}

	if !strings.Contains(stdout, "<table>") || !strings.Contains(stdout, "<th>id</th>") || !strings.Contains(stdout, "<td>Alice</td>") {
		t.Errorf("expected HTML table:\n%s", stdout)
	}
}

func TestWhereClause(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("run", "SELECT name FROM users WHERE age > 28 ORDER BY name", "-f", "csv")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}

	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(lines), lines)
	}
	if !strings.Contains(lines[1], "Alice") || !strings.Contains(lines[2], "Charlie") {
		t.Errorf("expected Alice and Charlie: %v", lines)
	}
}

func TestJoinQuery(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("run", "SELECT u.name, o.product FROM users u JOIN orders o ON u.id = o.user_id ORDER BY u.name, o.product", "-f", "csv")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}

	if !strings.Contains(stdout, "Alice") || !strings.Contains(stdout, "Bob") {
		t.Errorf("expected join results:\n%s", stdout)
	}
}

func TestNoResults(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, stderr, exitCode := env.RunSquix("run", "SELECT * FROM users WHERE id = 999", "-f", "json")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if stdout != "" {
		t.Errorf("expected empty stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "No results found") {
		t.Errorf("expected 'No results found' on stderr, got %q", stderr)
	}
}

func TestInvalidSQL(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	_, stderr, exitCode := env.RunSquix("run", "SELECTT * FROM users", "-f", "json")
	if exitCode == 0 {
		t.Fatal("expected non-zero exit code")
	}
	if !strings.Contains(stderr, "Error") {
		t.Errorf("expected error on stderr, got %q", stderr)
	}
}

func TestRunLast(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	env.RunSquix("add", "last_test", "SELECT 1 AS val")

	_, _, exitCode := env.RunSquix("run", "last_test", "-f", "csv")
	if exitCode != 0 {
		t.Fatalf("first run exit %d", exitCode)
	}

	last, _, exitCode := env.RunSquix("run", "--last", "-f", "csv")
	if exitCode != 0 {
		t.Fatalf("--last exit %d", exitCode)
	}

	if !strings.Contains(last, "1") {
		t.Errorf("expected result from last query, got %q", last)
	}
}
