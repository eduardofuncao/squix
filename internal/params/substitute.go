package params

import (
	"fmt"
	"strings"

	"github.com/eduardofuncao/squix/internal/db"
)

func SubstituteParameters(sql string, paramValues map[string]string, conn db.DatabaseConnection) (string, []any, error) {
	if len(paramValues) == 0 {
		return sql, []any{}, nil
	}

	matches := findSafeParamMatches(sql)

	if len(matches) == 0 {
		return sql, []any{}, nil
	}

	var orderedValues []any
	paramIndex := make(map[string]int)
	currentIndex := 1

	for _, m := range matches {
		if _, exists := paramIndex[m.name]; !exists {
			if value, ok := paramValues[m.name]; ok {
				paramIndex[m.name] = currentIndex
				orderedValues = append(orderedValues, value)
				currentIndex++
			} else {
				return "", nil, fmt.Errorf("missing value for parameter: %s", m.name)
			}
		}
	}

	result := replaceSafeMatches(sql, matches, func(m paramMatch) string {
		if idx, ok := paramIndex[m.name]; ok {
			return conn.GetPlaceholder(idx)
		}
		return sql[m.start:m.end]
	})

	return result, orderedValues, nil
}

func GenerateDisplaySQL(sql string, paramValues map[string]string) string {
	matches := findSafeParamMatches(sql)

	return replaceSafeMatches(sql, matches, func(m paramMatch) string {
		if value, ok := paramValues[m.name]; ok {
			if isNumeric(value) {
				return value
			}
			return "'" + strings.ReplaceAll(value, "'", "''") + "'"
		}
		return sql[m.start:m.end]
	})
}

// replaceSafeMatches rebuilds SQL by replacing each paramMatch using fn
func replaceSafeMatches(sql string, matches []paramMatch, fn func(paramMatch) string) string {
	if len(matches) == 0 {
		return sql
	}

	var buf strings.Builder
	prev := 0
	for _, m := range matches {
		buf.WriteString(sql[prev:m.start])
		buf.WriteString(fn(m))
		prev = m.end
	}
	buf.WriteString(sql[prev:])
	return buf.String()
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}

	hasDigits := false
	hasDot := false

	for i, r := range s {
		if r >= '0' && r <= '9' {
			hasDigits = true
		} else if r == '.' && !hasDot && i > 0 && i < len(s)-1 {
			hasDot = true
		} else if (r == '-' || r == '+') && i == 0 {
			// leading sign
		} else {
			return false
		}
	}

	return hasDigits
}
