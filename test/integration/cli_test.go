package integration

import (
	"encoding/json"
	"strings"
	"testing"
)

// --- SELECT + Export Tests ---

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
	if len(lines) != 3 { // header + Alice + Charlie
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

// --- CRUD Tests ---

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

// --- Parameterized Query Tests ---

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

// --- Query Management Tests ---

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
	t.Logf("output: %q", stdout)
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

// --- Explore + Explain Tests ---

func TestExploreListTables(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("explore")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "users") || !strings.Contains(stdout, "orders") {
		t.Errorf("expected users and orders:\n%s", stdout)
	}
}

func TestExplainTable(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("explain", "users")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "users") {
		t.Errorf("expected users in explain output:\n%s", stdout)
	}
}

func TestExplainWithDepth(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("explain", "users", "-d", "2")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "orders") {
		t.Errorf("expected FK relationship to orders:\n%s", stdout)
	}
}

func TestExplainVerbose(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("explain", "users", "-v")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "users") {
		t.Errorf("expected verbose explain:\n%s", stdout)
	}
}
