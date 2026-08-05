package table

import "testing"

func TestFormatCell(t *testing.T) {
	tests := []struct {
		name    string
		content string
		width   int
		want    string
	}{
		{"empty", "", 5, "     "},
		{"short pads right", "ab", 5, "ab   "},
		{"exact width", "abcde", 5, "abcde"},
		{"truncates with ellipsis", "abcdef", 5, "abcd…"},
		{"lf collapses to first line with marker", "Michael...\nrest", 15, "Michael...…    "},
		{"crlf collapses to first line with marker", "a\r\nb", 5, "a…   "},
		{"cr collapses to first line with marker", "a\rb", 5, "a…   "},
		{"leading newline becomes ellipsis", "\nhidden", 5, "…    "},
		{"no content after newline", "abc\n", 5, "abc… "},
		{"truncation still applies after collapse", "abcdef\nghij", 5, "abcd…"},
		{"vt collapses with marker", "a\vb", 5, "a…   "},
		{"ff collapses with marker", "a\fb", 5, "a…   "},
		{"tab becomes space", "a\tb", 5, "a b  "},
		{"leading tab preserved as space", "\tfoo", 5, " foo "},
		{"multiple tabs collapse to spaces", "a\t\tb", 5, "a  b "},
		{"tab then width truncate", "Tab\tSized Position", 15, "Tab Sized Posi…"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatCell(tt.content, tt.width); got != tt.want {
				t.Errorf("formatCell(%q,%d)=%q want %q", tt.content, tt.width, got, tt.want)
			}
		})
	}
}
