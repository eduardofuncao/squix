package editor

import (
	"strings"
)

// isCommentLine reports whether a line is a full-line SQL comment (its trimmed
// form starts with "--"). Such lines are treated as editor instructions/header
// text and stripped from the SQL the user typed.
func isCommentLine(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "--")
}

// HasInstructions reports whether content contains any full-line "--" comment,
// i.e. there is something for StripInstructions to remove.
func HasInstructions(content string) bool {
	for line := range strings.SplitSeq(content, "\n") {
		if isCommentLine(line) {
			return true
		}
	}
	return false
}

// StripInstructions removes every full-line "--" comment and returns the
// remaining (non-comment) lines joined with newlines.
func StripInstructions(content string) string {
	var sql strings.Builder
	for line := range strings.SplitSeq(content, "\n") {
		if isCommentLine(line) {
			continue
		}
		sql.WriteString(line)
		sql.WriteByte('\n')
	}
	return strings.TrimSpace(sql.String())
}
