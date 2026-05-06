package parser

import (
	"os"
	"strings"
	"testing"

	"github.com/eduardofuncao/squix/internal/styles"
)

func TestMain(m *testing.M) {
	styles.InitScheme("default", nil)
	os.Exit(m.Run())
}

func TestFormatSQLWithLineBreaks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		excludes []string
	}{
		{
			name:     "empty string",
			input:    "",
			contains: nil,
		},
		{
			name:     "simple SELECT",
			input:    "SELECT * FROM users WHERE id = 1",
			contains: []string{"SELECT", "FROM", "WHERE"},
		},
		{
			name:  "multi-keyword",
			input: "SELECT a FROM t JOIN u ON t.id = u.id WHERE x = 1 GROUP BY a ORDER BY a LIMIT 10",
			contains: []string{"SELECT", "FROM", "JOIN", "ON", "WHERE", "GROUP BY", "ORDER BY", "LIMIT"},
		},
		{
			name:     "INSERT INTO",
			input:    "INSERT INTO users VALUES (1)",
			contains: []string{"INSERT INTO", "VALUES"},
		},
		{
			name:     "UPDATE SET WHERE",
			input:    "UPDATE users SET name = 'x' WHERE id = 1",
			contains: []string{"UPDATE", "SET", "WHERE"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatSQLWithLineBreaks(tt.input)
			if tt.input == "" {
				if got != "" {
					t.Errorf("expected empty, got %q", got)
				}
				return
			}
			for _, kw := range tt.contains {
				if !strings.Contains(got, kw) {
					t.Errorf("expected keyword %q in output: %q", kw, got)
				}
			}
			for _, ex := range tt.excludes {
				if strings.Contains(got, ex) {
					t.Errorf("should not contain %q in output: %q", ex, got)
				}
			}
			lines := strings.Split(got, "\n")
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed != "" && strings.HasPrefix(trimmed, "\n") {
					t.Errorf("line starts with newline: %q", line)
				}
			}
		})
	}
}

func TestHighlightSQL(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		got := HighlightSQL("")
		if got != "" {
			t.Errorf("expected empty, got %q", got)
		}
	})

	t.Run("single keyword", func(t *testing.T) {
		got := HighlightSQL("SELECT")
		if got == "" {
			t.Error("expected non-empty output")
		}
	})

	t.Run("compound keyword", func(t *testing.T) {
		got := HighlightSQL("SELECT * FROM t LEFT JOIN u ON t.id = u.id")
		if got == "" {
			t.Error("expected non-empty output")
		}
	})

	t.Run("string literal", func(t *testing.T) {
		got := HighlightSQL("SELECT * FROM t WHERE name = 'hello'")
		if got == "" {
			t.Error("expected non-empty output")
		}
	})

	t.Run("no keywords", func(t *testing.T) {
		got := HighlightSQL("users.id")
		if got != "users.id" {
			t.Errorf("expected unchanged, got %q", got)
		}
	})
}

func TestCountLines(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 1},
		{"hello", 1},
		{"hello\nworld", 2},
		{"a\nb\nc", 3},
		{"hello\n", 2},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := countLines(tt.input)
			if got != tt.expected {
				t.Errorf("got %d, want %d", got, tt.expected)
			}
		})
	}
}
