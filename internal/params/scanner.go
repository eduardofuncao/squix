package params

import "strings"

type paramMatch struct {
	start      int
	end        int
	name       string
	defaultVal string
	hasDefault bool
}

// findSafeParamMatches scans SQL and returns :param or :param|default
// matches that are outside string literals, comments, and :: type casts
func findSafeParamMatches(sql string) []paramMatch {
	var result []paramMatch
	n := len(sql)
	i := 0

	for i < n {
		ch := sql[i]

		if ch == '\'' {
			i = skipString(sql, i)
			continue
		}

		if ch == '-' && i+1 < n && sql[i+1] == '-' {
			for i < n && sql[i] != '\n' {
				i++
			}
			continue
		}

		if ch == '/' && i+1 < n && sql[i+1] == '*' {
			i += 2
			for i+1 < n {
				if sql[i] == '*' && sql[i+1] == '/' {
					i += 2
					break
				}
				i++
			}
			continue
		}

		// :: postgres cast — skip both colons and the type name
		if ch == ':' && i+1 < n && sql[i+1] == ':' {
			i += 2
			for i < n && isIdentChar(sql[i]) {
				i++
			}
			continue
		}

		if ch == ':' && i+1 < n && isIdentStart(sql[i+1]) {
			start := i
			i++
			nameStart := i
			for i < n && isIdentChar(sql[i]) {
				i++
			}
			name := sql[nameStart:i]

			m := paramMatch{start: start, name: name}

			if i < n && sql[i] == '|' {
				m.hasDefault = true
				i++ // skip |
				if i < n && sql[i] == '\'' {
					// quoted default
					i++
					var def strings.Builder
					for i < n {
						if sql[i] == '\\' && i+1 < n {
							def.WriteByte(sql[i+1])
							i += 2
							continue
						}
						if sql[i] == '\'' {
							// check for '' escape
							if i+1 < n && sql[i+1] == '\'' {
								def.WriteByte('\'')
								i += 2
								continue
							}
							i++
							break
						}
						def.WriteByte(sql[i])
						i++
					}
					m.defaultVal = def.String()
				} else {
					// unquoted default
					defStart := i
					for i < n && sql[i] != ' ' && sql[i] != '\t' && sql[i] != '\n' && sql[i] != '\r' && sql[i] != ')' && sql[i] != ',' {
						i++
					}
					m.defaultVal = sql[defStart:i]
				}
			}

			m.end = i
			result = append(result, m)
			continue
		}

		i++
	}

	return result
}

func skipString(sql string, i int) int {
	n := len(sql)
	i++ // skip opening quote
	for i < n {
		if sql[i] == '\'' {
			if i+1 < n && sql[i+1] == '\'' {
				i += 2
				continue
			}
			return i + 1
		}
		i++
	}
	return i
}

func isIdentStart(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func isIdentChar(b byte) bool {
	return isIdentStart(b) || (b >= '0' && b <= '9')
}
