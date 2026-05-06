package table

import (
	"encoding/json"
	"strings"
	"testing"
)

var (
	testHeaders = []string{"id", "name", "email"}
	testRows    = [][]string{
		{"1", "Alice", "alice@example.com"},
		{"2", "Bob", "bob@example.com"},
	}
	emptyRows = [][]string{}
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"csv", "csv", false},
		{"CSV", "csv", false},
		{"json", "json", false},
		{"JSON", "json", false},
		{"tsv", "tsv", false},
		{"html", "html", false},
		{"sql", "sql", false},
		{"markdown", "markdown", false},
		{"md", "markdown", false},
		{"Md", "markdown", false},
		{"", "", true},
		{"xml", "", true},
		{"csv2", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestFormatCSV(t *testing.T) {
	t.Run("standard output", func(t *testing.T) {
		got, err := FormatCSV(testHeaders, testRows)
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(strings.TrimSpace(got), "\n")
		if len(lines) != 3 {
			t.Fatalf("expected 3 lines, got %d", len(lines))
		}
		if lines[0] != "id,name,email" {
			t.Errorf("header = %q", lines[0])
		}
	})

	t.Run("empty rows", func(t *testing.T) {
		got, err := FormatCSV(testHeaders, emptyRows)
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(strings.TrimSpace(got), "\n")
		if len(lines) != 1 {
			t.Fatalf("expected 1 line, got %d", len(lines))
		}
	})

	t.Run("field with comma", func(t *testing.T) {
		got, err := FormatCSV([]string{"name"}, [][]string{{"hello,world"}})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(got, `"hello,world"`) {
			t.Errorf("expected quoted comma field, got %q", got)
		}
	})

	t.Run("field with quotes", func(t *testing.T) {
		got, err := FormatCSV([]string{"name"}, [][]string{{`he said "hi"`}})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(got, `"he said ""hi"""`) {
			t.Errorf("expected escaped quotes, got %q", got)
		}
	})
}

func TestFormatJSON(t *testing.T) {
	t.Run("standard output", func(t *testing.T) {
		got, err := FormatJSON(testHeaders, testRows)
		if err != nil {
			t.Fatal(err)
		}

		var results []map[string]string
		if err := json.Unmarshal([]byte(got), &results); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		if len(results) != 2 {
			t.Fatalf("expected 2 objects, got %d", len(results))
		}
		if results[0]["name"] != "Alice" {
			t.Errorf("results[0][name] = %q", results[0]["name"])
		}
	})

	t.Run("empty rows", func(t *testing.T) {
		got, err := FormatJSON(testHeaders, emptyRows)
		if err != nil {
			t.Fatal(err)
		}
		if got != "[]" {
			t.Errorf("expected '[]', got %q", got)
		}
	})

	t.Run("row with fewer columns", func(t *testing.T) {
		got, err := FormatJSON([]string{"a", "b", "c"}, [][]string{{"1"}})
		if err != nil {
			t.Fatal(err)
		}
		var results []map[string]string
		json.Unmarshal([]byte(got), &results)
		if results[0]["a"] != "1" {
			t.Errorf("expected a=1, got %q", results[0]["a"])
		}
		if _, exists := results[0]["b"]; exists {
			t.Error("expected b to be absent")
		}
	})
}

func TestFormatTSV(t *testing.T) {
	t.Run("standard output", func(t *testing.T) {
		got, err := FormatTSV(testHeaders, testRows)
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(strings.TrimSpace(got), "\n")
		if len(lines) != 3 {
			t.Fatalf("expected 3 lines, got %d", len(lines))
		}
		if !strings.Contains(lines[0], "\t") {
			t.Error("expected tabs in header")
		}
	})

	t.Run("empty rows", func(t *testing.T) {
		got, err := FormatTSV(testHeaders, emptyRows)
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(strings.TrimSpace(got), "\n")
		if len(lines) != 1 {
			t.Fatalf("expected 1 line, got %d", len(lines))
		}
	})
}

func TestFormatHTML(t *testing.T) {
	t.Run("standard output", func(t *testing.T) {
		got, err := FormatHTML(testHeaders, testRows, FormatOptions{QueryName: "test"})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(got, "<!DOCTYPE html>") {
			t.Error("missing DOCTYPE")
		}
		if !strings.Contains(got, "<th>id</th>") {
			t.Error("missing table header")
		}
		if !strings.Contains(got, "<td>Alice</td>") {
			t.Error("missing table data")
		}
		if !strings.Contains(got, "<h3>") {
			t.Error("missing title")
		}
	})

	t.Run("alternating row classes", func(t *testing.T) {
		got, err := FormatHTML(testHeaders, testRows, FormatOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(got, `class="odd"`) {
			t.Error("expected odd row class on second row")
		}
	})

	t.Run("HTML escaping", func(t *testing.T) {
		got, err := FormatHTML([]string{"val"}, [][]string{{`a <b> & "c" 'd'`}}, FormatOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(got, "a &lt;b&gt; &amp; &quot;c&quot; &#39;d&#39;") {
			t.Errorf("unexpected escaping: %s", got)
		}
	})

	t.Run("no title when QueryName empty", func(t *testing.T) {
		got, err := FormatHTML(testHeaders, testRows, FormatOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(got, "<h3>") {
			t.Error("expected no title")
		}
	})

	t.Run("empty rows", func(t *testing.T) {
		got, err := FormatHTML(testHeaders, emptyRows, FormatOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(got, "<thead>") {
			t.Error("missing thead")
		}
		if !strings.Contains(got, "<tbody>") {
			t.Error("missing tbody")
		}
	})
}

func TestFormatSQL(t *testing.T) {
	t.Run("standard output", func(t *testing.T) {
		got, err := FormatSQL(testHeaders, testRows, FormatOptions{TableName: "users"})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(got, "INSERT INTO users") {
			t.Error("missing INSERT INTO")
		}
		if !strings.Contains(got, `"id"`) {
			t.Error("missing quoted column name")
		}
	})

	t.Run("empty table name returns error", func(t *testing.T) {
		_, err := FormatSQL(testHeaders, testRows, FormatOptions{})
		if err == nil {
			t.Error("expected error for empty table name")
		}
	})

	t.Run("NULL handling", func(t *testing.T) {
		got, err := FormatSQL([]string{"a"}, [][]string{{""}, {"NULL"}}, FormatOptions{TableName: "t"})
		if err != nil {
			t.Fatal(err)
		}
		if strings.Count(got, "NULL") != 2 {
			t.Errorf("expected 2 NULLs, got: %s", got)
		}
	})

	t.Run("single quote escaping", func(t *testing.T) {
		got, err := FormatSQL([]string{"name"}, [][]string{{"O'Brien"}}, FormatOptions{TableName: "t"})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(got, `'O''Brien'`) {
			t.Errorf("expected escaped quotes, got: %s", got)
		}
	})

	t.Run("multiple rows", func(t *testing.T) {
		got, err := FormatSQL(testHeaders, testRows, FormatOptions{TableName: "users"})
		if err != nil {
			t.Fatal(err)
		}
		if strings.Count(got, "INSERT INTO") != 2 {
			t.Errorf("expected 2 INSERT statements, got: %s", got)
		}
	})
}

func TestFormatMarkdown(t *testing.T) {
	t.Run("standard output", func(t *testing.T) {
		got, err := FormatMarkdown(testHeaders, testRows)
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(got, "\n")
		if len(lines) != 5 {
			t.Fatalf("expected 5 lines, got %d", len(lines))
		}
		if !strings.HasPrefix(lines[0], "|") {
			t.Error("header should start with |")
		}
		if !strings.Contains(lines[1], "---") {
			t.Error("separator should contain ---")
		}
	})

	t.Run("empty rows", func(t *testing.T) {
		got, err := FormatMarkdown(testHeaders, emptyRows)
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(got, "\n")
		if len(lines) != 3 {
			t.Fatalf("expected 3 lines (header + separator + trailing newline), got %d", len(lines))
		}
	})
}

func TestFormatExport(t *testing.T) {
	t.Run("dispatches to correct format", func(t *testing.T) {
		formats := []string{"csv", "json", "tsv", "html", "markdown"}
		for _, f := range formats {
			got, err := FormatExport(testHeaders, testRows, f, FormatOptions{})
			if err != nil {
				t.Errorf("format %q: %v", f, err)
			}
			if got == "" {
				t.Errorf("format %q: empty output", f)
			}
		}
	})

	t.Run("sql format needs table name", func(t *testing.T) {
		_, err := FormatExport(testHeaders, testRows, "sql", FormatOptions{})
		if err == nil {
			t.Error("expected error for sql without table name")
		}
	})

	t.Run("unknown format defaults to CSV", func(t *testing.T) {
		got, err := FormatExport(testHeaders, testRows, "unknown", FormatOptions{})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(got, "id,name,email") {
			t.Errorf("expected CSV output, got: %s", got)
		}
	})
}
