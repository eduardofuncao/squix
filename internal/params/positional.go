package params

func MapPositionalArgs(sql string, positionals []string) map[string]string {
	result := make(map[string]string)

	if len(positionals) == 0 {
		return result
	}

	paramDefs := ExtractParameters(sql)
	matches := findSafeParamMatches(sql)

	var paramNames []string
	seen := make(map[string]bool)

	for _, m := range matches {
		if _, exists := paramDefs[m.name]; exists && !seen[m.name] {
			paramNames = append(paramNames, m.name)
			seen[m.name] = true
		}
	}

	for i, value := range positionals {
		if i < len(paramNames) {
			result[paramNames[i]] = value
		}
	}

	return result
}
