package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/eduardofuncao/squix/internal/config"
	"github.com/eduardofuncao/squix/internal/db"
	"github.com/eduardofuncao/squix/internal/editor"
	"github.com/eduardofuncao/squix/internal/params"
	"github.com/eduardofuncao/squix/internal/run"
)

func (a *App) handleRun() {
	if a.config.CurrentConnection == "" {
		printError("No active connection.   Use 'squix switch <connection>' or 'squix init' first")
	}

	conn := config.FromConnectionYaml(a.config.Connections[a.config.CurrentConnection])
	args := os.Args[2:]

	// Check for empty args → open editor for new query (CLI-only feature)
	if len(args) == 0 {
		newQuery := a.createNewQueryOrEdit()
		if err := a.executeQueryWithParams(newQuery, conn, nil, nil); err != nil {
			printError("%v", err)
		}
		return
	}

	if err := a.runFromArgs(args, conn); err != nil {
		printError("%v", err)
	}
}

func (a *App) runFromArgs(args []string, conn db.DatabaseConnection) error {
	flags := parseRunFlagsFrom(args)

	resolved, err := run.ResolveQuery(flags, a.config, a.config.CurrentConnection, conn)
	if err != nil {
		return err
	}

	// Check if we need to create a new query via editor
	if run.ShouldCreateNewQuery(resolved) {
		return fmt.Errorf("no query specified. Use a query name, inline SQL, or --last")
	}

	if flags.EditMode && !flags.LastQuery {
		resolved.Query = a.editQueryOrExit(resolved.Query)
	}

	a.saveIfNeeded(resolved)

	// Parse parameter flags and positional args
	paramFlags := parseParameterFlagsFrom(args)
	positionalArgsSlice := parsePositionalArgsFrom(args, flags.Selector)

	positionalArgs := params.MapPositionalArgs(resolved.Query.SQL, positionalArgsSlice)

	return a.executeQueryWithParams(resolved.Query, conn, paramFlags, positionalArgs)
}

func parseRunFlagsFrom(args []string) run.Flags {
	flags := run.Flags{}

	for i, arg := range args {
		// Skip parameter flags and their values
		if strings.HasPrefix(arg, "--") && arg != "--edit" && arg != "-e" && arg != "--last" && arg != "-l" {
			// This is a parameter flag, skip it and its value
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				continue
			}
		}

		switch arg {
		case "--edit", "-e":
			flags.EditMode = true
		case "--last", "-l":
			flags.LastQuery = true
		default:
			if !strings.HasPrefix(arg, "--") && flags.Selector == "" {
				flags.Selector = arg
			}
		}
	}
	return flags
}

func parseParameterFlagsFrom(args []string) map[string]string {
	paramValues := make(map[string]string)

	i := 0
	for i < len(args) {
		arg := args[i]

		// Skip known flags
		if arg == "--edit" || arg == "-e" || arg == "--last" || arg == "-l" {
			i++
			continue
		}

		// Check if it's a parameter flag (--param value)
		if strings.HasPrefix(arg, "--") && len(arg) > 2 {
			paramName := strings.TrimPrefix(arg, "--")

			// If next arg exists and doesn't start with --, it's the value
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				paramValues[paramName] = args[i+1]
				i += 2 // Skip next arg as we've consumed it
			} else {
				// Flag without value - set to empty string (will use default)
				paramValues[paramName] = ""
				i++
			}
		} else {
			// Skip non-flag args (query selector and positional params handled elsewhere)
			i++
		}
	}

	return paramValues
}

// parsePositionalArgsFrom extracts positional parameter values
// Returns slice of values in order (after query selector)
func parsePositionalArgsFrom(args []string, selector string) []string {
	var positionals []string
	selectorSeen := false

	i := 0
	for i < len(args) {
		arg := args[i]

		// Skip flags and their values
		if strings.HasPrefix(arg, "--") {
			// Skip the flag itself
			i++
			// Skip the value if it exists and isn't another flag
			if i < len(args) && !strings.HasPrefix(args[i], "--") {
				i++
			}
			continue
		}
		if arg == "-e" || arg == "-l" {
			i++
			continue
		}

		// First non-flag arg is the query selector
		if !selectorSeen {
			selectorSeen = true
			i++
			continue
		}

		// Remaining non-flag args are positional parameter values
		positionals = append(positionals, arg)
		i++
	}

	return positionals
}

func (a *App) createNewQueryOrEdit() db.Query {
	instructions := `-- Enter your SQL run below
-- Save and exit to execute, or exit without saving to cancel
--
`
	editedSQL, err := editor.EditTempFileWithTemplate(instructions, "squix-run-")
	if err != nil {
		printError("Error opening editor: %v", err)
	}
	if editedSQL == "" {
		printError("Empty SQL, cancelled")
	}
	return db.Query{Name: "<runtime>", SQL: editedSQL, Id: -1}
}

func (a *App) editQueryOrExit(query db.Query) db.Query {
	editedSQL, err := editor.EditTempFile(query.SQL, "squix-run-")
	if err != nil {
		printError("Error opening editor: %v", err)
	}
	query.SQL = editedSQL
	return query
}

func (a *App) saveIfNeeded(resolved run.ResolvedQuery) {
	if !resolved.Saveable {
		return
	}

	// Save the query and update last query
	if err := a.config.SaveQueryAndLast(a.config.CurrentConnection, resolved.Query, true); err != nil {
		printError("Failed to save query: %v", err)
	}
}

type executorFunc func(run.ExecutionParams) error

func (a *App) executeQueryWithParamsInternal(query db.Query, conn db.DatabaseConnection, paramFlags, positionalArgs map[string]string, executor executorFunc) error {
	sql, args, displaySQL, err := a.processParameters(query.SQL, conn, paramFlags, positionalArgs)
	if err != nil {
		return err
	}

	processedQuery := db.Query{
		Name: query.Name,
		SQL:  sql,
		Id:   query.Id,
	}

	var onRerun func(string) error
	onRerun = func(editedSQL string) error {
		finalSQL := editedSQL
		finalArgs := []any{}
		finalDisplaySQL := ""

		if strings.Contains(editedSQL, ":") {
			var procErr error
			finalSQL, finalArgs, finalDisplaySQL, procErr = a.processParameters(editedSQL, conn, paramFlags, positionalArgs)
			if procErr != nil {
				return procErr
			}
		}
		if finalDisplaySQL == "" {
			finalDisplaySQL = finalSQL
		}

		processedQuery := db.Query{
			Name: query.Name,
			SQL:  finalSQL,
			Id:   query.Id,
		}

		return executor(run.ExecutionParams{
			Query:        processedQuery,
			Connection:   conn,
			Config:       a.config,
			SaveCallback: a.saveQueryFromTable,
			Args:         finalArgs,
			DisplaySQL:   finalDisplaySQL,
			OnRerun:      onRerun,
		})
	}

	return executor(run.ExecutionParams{
		Query:        processedQuery,
		Connection:   conn,
		Config:       a.config,
		SaveCallback: a.saveQueryFromTable,
		Args:         args,
		DisplaySQL:   displaySQL,
		OnRerun:      onRerun,
	})
}

func (a *App) executeQueryWithParams(query db.Query, conn db.DatabaseConnection, paramFlags, positionalArgs map[string]string) error {
	return a.executeQueryWithParamsInternal(query, conn, paramFlags, positionalArgs, run.Execute)
}

// processParameters handles parameter extraction, validation, and substitution
func (a *App) processParameters(sql string, conn db.DatabaseConnection, cliValues, positionals map[string]string) (string, []any, string, error) {
	// Extract parameter definitions from SQL
	paramDefs := params.ExtractParameters(sql)

	if len(paramDefs) == 0 {
		return sql, []any{}, "", nil
	}

	// Map positional args to parameter names
	for k, v := range positionals {
		cliValues[k] = v
	}

	// Validate CLI values
	if err := params.ValidateCLIValues(cliValues, paramDefs); err != nil {
		return "", nil, "", fmt.Errorf("parameter validation error: %w", err)
	}

	// Validate param names don't conflict with reserved flags
	if err := params.ValidateParamNames(paramDefs); err != nil {
		return "", nil, "", fmt.Errorf("parameter name conflict: %w", err)
	}

	// Resolve parameters (CLI > defaults)
	paramValues := params.ResolveParameters(paramDefs, cliValues)

	// Check for missing required parameters
	missing := params.GetMissingRequired(paramDefs, paramValues)
	if len(missing) > 0 {
		// Launch interactive TUI
		collectedValues, err := params.CollectParameters(sql, missing, paramDefs)
		if err != nil {
			return "", nil, "", fmt.Errorf("error collecting parameters: %w", err)
		}
		// Merge collected values
		for k, v := range collectedValues {
			paramValues[k] = v
		}
	}

	// Substitute parameters with DB-specific placeholders
	finalSQL, args, err := params.SubstituteParameters(sql, paramValues, conn)
	if err != nil {
		return "", nil, "", fmt.Errorf("error substituting parameters: %w", err)
	}

	// Generate display SQL with actual values for TUI
	displaySQL := params.GenerateDisplaySQL(sql, paramValues)

	// For Oracle, use literal substitution instead of prepared statements
	// This is a workaround for Oracle driver prepared statement issues
	if conn.GetDbType() == "oracle" && len(args) > 0 {
		finalSQL = substituteOracleLiterals(finalSQL, args)
		args = []any{}
	}

	return finalSQL, args, displaySQL, nil
}

// substituteOracleLiterals replaces :1, :2 placeholders with actual values for Oracle
func substituteOracleLiterals(sql string, args []any) string {
	result := sql
	for i, arg := range args {
		placeholder := fmt.Sprintf(":%d", i+1)
		var value string
		switch v := arg.(type) {
		case string:
			// Escape single quotes in strings
			escaped := strings.ReplaceAll(v, "'", "''")
			value = fmt.Sprintf("'%s'", escaped)
		case int, int32, int64:
			value = fmt.Sprintf("%d", v)
		case float32, float64:
			value = fmt.Sprintf("%f", v)
		default:
			// For other types, try to convert to string
			value = fmt.Sprintf("'%v'", v)
		}
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

func (a *App) saveQueryFromTable(query db.Query) (db.Query, error) {
	connName := a.config.CurrentConnection
	if connName == "" {
		return db.Query{}, fmt.Errorf("no active connection")
	}

	// Save query with auto-ID generation
	savedQuery, err := a.config.SaveQueryToConnection(connName, query)
	if err != nil {
		return db.Query{}, err
	}

	// Update last query
	if err := a.config.UpdateLastQuery(connName, savedQuery); err != nil {
		return db.Query{}, err
	}

	return savedQuery, nil
}
