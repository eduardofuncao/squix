package table

import (
	"reflect"
	"testing"
)

func TestWrapText(t *testing.T) {
	tests := []struct {
		name    string
		content string
		width   int
		want    []string
	}{
		{"empty", "", 10, []string{""}},
		{"short line", "hello", 10, []string{"hello"}},
		{
			"word wrap drops boundary space",
			"The quick brown fox", 10,
			[]string{"The quick", "brown fox"},
		},
		{
			"hard-break long token",
			"abcdefghij", 4,
			[]string{"abcd", "efgh", "ij"},
		},
		{
			"preserve newlines",
			"a\nb", 10,
			[]string{"a", "b"},
		},
		{
			"blank lines preserved",
			"a\n\nb", 10,
			[]string{"a", "", "b"},
		},
		{
			"CRLF normalized to LF",
			"a\r\nb", 10,
			[]string{"a", "b"},
		},
		{
			"lone CR treated as line break",
			"a\rb", 10,
			[]string{"a", "b"},
		},
		{"width guard clamps to 1", "ab", 0, []string{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapText(tt.content, tt.width)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wrapText(%q, %d) = %#v, want %#v", tt.content, tt.width, got, tt.want)
			}
		})
	}
}

func TestWrapLineNoContentLoss(t *testing.T) {
	// A long base64-like blob must hard-break with every rune preserved.
	line := "AAABBBCCCDDDEEEFFFGGGHHHIIIJJJKKKLLLMMMNNNOOOPPP"
	wrapped := wrapLine(line, 11)
	joined := ""
	for _, w := range wrapped {
		joined += w
	}
	if joined != line {
		t.Errorf("content lost after wrap: got %q, want %q", joined, line)
	}
	for _, w := range wrapped {
		if len([]rune(w)) > 11 {
			t.Errorf("segment exceeds width: %q (%d runes)", w, len([]rune(w)))
		}
	}
}
