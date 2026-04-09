package params

func ExtractParameters(sql string) map[string]string {
	params := make(map[string]string)
	matches := findSafeParamMatches(sql)

	seen := make(map[string]bool)
	for _, m := range matches {
		if !seen[m.name] {
			seen[m.name] = true
			params[m.name] = m.defaultVal
		}
	}

	return params
}
