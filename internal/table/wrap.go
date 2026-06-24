package table

import "strings"

// wrapText splits content into display lines, word-wrapping each source line to
// at most width runes. Words longer than width are hard-broken so no content is
// lost. Newlines in the source are preserved as line breaks. Empty input yields
// a single empty line.
func wrapText(content string, width int) []string {
	width = max(width, 1)

	// Normalize CRLF and lone CR to LF so carriage returns don't survive into
	// the output as visible runes that throw off wrapping and alignment.
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	var out []string
	for line := range strings.SplitSeq(content, "\n") {
		out = append(out, wrapLine(line, width)...)
	}

	if len(out) == 0 {
		out = []string{""}
	}
	return out
}

// wrapLine word-wraps a single line (no embedded newlines) to width runes.
func wrapLine(line string, width int) []string {
	width = max(width, 1)

	runes := []rune(line)
	if len(runes) <= width {
		return []string{line}
	}

	var lines []string
	for len(runes) > 0 {
		if len(runes) <= width {
			lines = append(lines, string(runes))
			break
		}

		// Break at the last whitespace within the width window; this keeps words
		// intact and fills each line greedily.
		breakAt := -1
		for i := width - 1; i >= 0; i-- {
			if isSpace(runes[i]) {
				breakAt = i
				break
			}
		}

		if breakAt <= 0 {
			// No usable whitespace in the window: hard-break at width.
			lines = append(lines, string(runes[:width]))
			runes = runes[width:]
			continue
		}

		// Emit up to the break and resume after the break space. Only the single
		// boundary space is dropped; non-whitespace content is never lost.
		lines = append(lines, string(runes[:breakAt]))
		runes = runes[breakAt+1:]
	}

	if len(lines) == 0 {
		lines = []string{""}
	}
	return lines
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// wrapWidth returns the rune width available for detail-view content.
func (m Model) wrapWidth() int {
	return max(m.width-4, 10)
}

// detailViewportHeight is the number of content lines visible in the detail view
// (terminal height minus header/footer chrome). Shared by render and scroll.
func (m Model) detailViewportHeight() int {
	return max(m.height-10, 5)
}

// wrappedDetailLines returns the detail-view content word-wrapped to the current
// viewport width.
func (m Model) wrappedDetailLines() []string {
	return wrapText(m.detailViewContent, m.wrapWidth())
}
