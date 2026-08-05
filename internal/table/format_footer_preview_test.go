package table

import "testing"

func TestFormatFooterPreview(t *testing.T) {
	tests := []struct {
		name  string
		value string
		max   int
		want  string
	}{
		{"empty", "", 10, ""},
		{"short passthrough", "abc", 10, "abc"},
		{"exact width", "abcde", 5, "abcde"},
		{"truncates with ellipsis", "abcdef", 5, "abcd…"},
		{"max zero returns full", "abcdef", 0, "abcdef"},
		{"max negative returns full", "abcdef", -1, "abcdef"},
		{"lf collapses to first line with marker", "Michael...\nrest", 20, "Michael...…"},
		{"crlf collapses to first line with marker", "a\r\nb", 10, "a…"},
		{"cr collapses to first line with marker", "a\rb", 10, "a…"},
		{"leading newline becomes ellipsis", "\nhidden", 10, "…"},
		{"truncate after collapse", "abcdef\nghij", 4, "abc…"},
		{"vt collapses with marker", "a\vb", 10, "a…"},
		{"ff collapses with marker", "a\fb", 10, "a…"},
		{"tab becomes space", "a\tb", 10, "a b"},
		{"leading tab preserved as space", "\tfoo", 10, " foo"},
		{"tab then width truncate", "Tab\tSized Position", 8, "Tab Siz…"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatFooterPreview(tt.value, tt.max); got != tt.want {
				t.Errorf("formatFooterPreview(%q,%d)=%q want %q", tt.value, tt.max, got, tt.want)
			}
		})
	}
}
