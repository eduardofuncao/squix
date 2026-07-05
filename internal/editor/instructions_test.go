package editor

import "testing"

func TestStripInstructions(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "bug #76: SQL typed right after first instruction line",
			in: "-- Enter your SQL run below\n" +
				"SELECT * FROM testtable;\n" +
				"-- Save and exit to execute, or exit without saving to cancel\n" +
				"--\n",
			want: "SELECT * FROM testtable;",
		},
		{
			name: "SQL typed below the -- separator (normal case)",
			in: "-- Enter your SQL run below\n" +
				"-- Save and exit to execute, or exit without saving to cancel\n" +
				"--\n" +
				"SELECT * FROM testtable;\n",
			want: "SELECT * FROM testtable;",
		},
		{
			name: "re-run template keeps only the existing query",
			in: "-- Enter your SQL run below\n" +
				"-- Last query (edit, clear, or replace). Save to re-run\n" +
				"--\n" +
				"SELECT 1;\n" +
				"SELECT 2;\n",
			want: "SELECT 1;\nSELECT 2;",
		},
		{
			name: "add.go 'Creating new run' template",
			in: "-- Creating new run:  myrun\n" +
				"-- Connection: dev (sqlite)\n" +
				"-- Write your SQL run below and save\n" +
				"\n" +
				"SELECT * FROM users;\n",
			want: "SELECT * FROM users;",
		},
		{
			name: "multiline SQL interleaved with comments",
			in: "-- header\n" +
				"WITH t AS (\n" +
				"  SELECT 1\n" +
				")\n" +
				"-- inline note\n" +
				"SELECT * FROM t;\n",
			want: "WITH t AS (\n  SELECT 1\n)\nSELECT * FROM t;",
		},
		{
			name: "trailing inline comment is preserved",
			in:   "SELECT 1 -- keep me\n",
			want: "SELECT 1 -- keep me",
		},
		{
			name: "plain SQL with no comments is unchanged",
			in:   "SELECT 1;\nSELECT 2;\n",
			want: "SELECT 1;\nSELECT 2;",
		},
		{
			name: "only comments yields empty string",
			in: "-- Enter your SQL run below\n" +
				"-- Save and exit to execute, or exit without saving to cancel\n" +
				"--\n",
			want: "",
		},
		{
			name: "empty input",
			in:   "",
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := StripInstructions(tc.in)
			if got != tc.want {
				t.Errorf("StripInstructions(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestHasInstructions(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"-- Enter your SQL run below\nSELECT 1;\n", true},
		{"SELECT 1;\n", false},
		{"", false},
		{"SELECT 1 -- inline\n", false}, // inline comment is not a full-line comment
		{"  -- indented comment\nSELECT 1;\n", true},
	}

	for _, tc := range tests {
		got := HasInstructions(tc.in)
		if got != tc.want {
			t.Errorf("HasInstructions(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
