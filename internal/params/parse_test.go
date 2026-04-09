package params

import (
	"testing"
)

func TestExtractParameters_Basic(t *testing.T) {
	sql := "SELECT * FROM users WHERE name = :name AND age > :age|25"
	params := ExtractParameters(sql)

	if len(params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(params))
	}
	if params["name"] != "" {
		t.Errorf("expected empty default for name, got %q", params["name"])
	}
	if params["age"] != "25" {
		t.Errorf("expected default 25 for age, got %q", params["age"])
	}
}

func TestExtractParameters_QuotedDefault(t *testing.T) {
	sql := "SELECT * FROM users WHERE name = :name|'Alice'"
	params := ExtractParameters(sql)

	if len(params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(params))
	}
	if params["name"] != "Alice" {
		t.Errorf("expected default Alice, got %q", params["name"])
	}
}

func TestExtractParameters_StringLiteral(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want int
	}{
		{"time in string", "SELECT 'Time: 12:30'", 0},
		{"colon in string", "SELECT * FROM t WHERE x = '10:30'", 0},
		{"param outside string", "SELECT 'hello' FROM t WHERE x = :val", 1},
		{"mixed", "SELECT 'a:fake' FROM t WHERE x = :real|5", 1},
		{"escaped quotes", "SELECT 'it''s :not a param'", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := ExtractParameters(tt.sql)
			if len(params) != tt.want {
				t.Errorf("sql=%q: expected %d params, got %d (%v)", tt.sql, tt.want, len(params), params)
			}
		})
	}
}

func TestExtractParameters_PostgresCast(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want int
	}{
		{"simple cast", "SELECT col::text FROM t", 0},
		{"cast with param", "SELECT col::text FROM t WHERE id = :id", 1},
		{"double cast", "SELECT a::int, b::varchar FROM t", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := ExtractParameters(tt.sql)
			if len(params) != tt.want {
				t.Errorf("sql=%q: expected %d params, got %d (%v)", tt.sql, tt.want, len(params), params)
			}
		})
	}
}

func TestExtractParameters_Comments(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want int
	}{
		{"line comment", "SELECT * FROM t -- WHERE x = :fake", 0},
		{"block comment", "SELECT * FROM t /* :fake */ WHERE x = :real", 1},
		{"param before comment", "SELECT * FROM t WHERE x = :val -- comment", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := ExtractParameters(tt.sql)
			if len(params) != tt.want {
				t.Errorf("sql=%q: expected %d params, got %d (%v)", tt.sql, tt.want, len(params), params)
			}
		})
	}
}

func TestExtractParameters_Duplicates(t *testing.T) {
	sql := "SELECT * FROM t WHERE x = :id AND y = :id"
	params := ExtractParameters(sql)
	if len(params) != 1 {
		t.Errorf("expected 1 unique param, got %d", len(params))
	}
}
