package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/eduardofuncao/squix/internal/config"
	"github.com/eduardofuncao/squix/internal/db"
	"github.com/eduardofuncao/squix/internal/run"
	"github.com/eduardofuncao/squix/internal/styles"
	"github.com/eduardofuncao/squix/internal/table"
)

func (a *App) handleExplore() {
	if len(os.Args) < 3 {
		a.listTablesAndViews()
		return
	}

	// Check if argument is a supported file type
	if ext, ok := getSupportedExtension(os.Args[2]); ok {
		if _, err := os.Stat(os.Args[2]); err != nil {
			printError("File not found: %s", os.Args[2])
		}
		if err := a.handleExploreFile(os.Args[2], ext); err != nil {
			printError("%v", err)
		}
		return
	}

	tableName := os.Args[2]
	limit := a.config.DefaultRowLimit
	if limit == 0 {
		limit = 1000
	}

	// Parse optional -l/--limit flag
	for i := 3; i < len(os.Args); i++ {
		if os.Args[i] == "-l" || os.Args[i] == "--limit" {
			if i+1 < len(os.Args) {
				parsed, err := strconv.Atoi(os.Args[i+1])
				if err != nil {
					printError("Invalid limit value: %s", os.Args[i+1])
				}
				limit = parsed
			}
			break
		}
	}

	if a.config.CurrentConnection == "" {
		printError("No active connection. Use 'squix switch <connection>' or 'squix init' first")
	}

	conn := config.FromConnectionYaml(a.config.Connections[a.config.CurrentConnection])

	if err := conn.Open(); err != nil {
		printError("Could not open connection: %v", err)
	}
	defer conn.Close()

	sql := fmt.Sprintf("SELECT * FROM %s", tableName)
	sql = conn.ApplyRowLimit(sql, limit)

	var onRerun func(string) error
	onRerun = func(newSQL string) error {
		return run.ExecuteSelect(newSQL, tableName, run.ExecutionParams{
			Query:      db.Query{Name: tableName, SQL: newSQL},
			Connection: conn,
			Config:     a.config,
			OnRerun:    onRerun,
		})
	}
	if err := run.ExecuteSelect(sql, tableName, run.ExecutionParams{
		Query:      db.Query{Name: tableName, SQL: sql},
		Connection: conn,
		Config:     a.config,
		OnRerun:    onRerun,
	}); err != nil {
		printError("%v", err)
	}
}

func (a *App) listTablesAndViews() {
	if a.config.CurrentConnection == "" {
		printError(
			"No active connection. Use 'squix switch <connection>' or 'squix init' first",
		)
	}

	conn := config.FromConnectionYaml(
		a.config.Connections[a.config.CurrentConnection],
	)

	if err := conn.Open(); err != nil {
		printError(
			"Could not open connection to %s: %v",
			a.config.CurrentConnection,
			err,
		)
	}
	defer conn.Close()

	tables, err := conn.GetTables()
	if err != nil {
		printError("Could not list tables: %v", err)
	}

	views, err := conn.GetViews()
	if err != nil {
		printError("Could not list views: %v", err)
	}

	if len(tables) > 0 {
		fmt.Printf("%s tables %s\n", styles.Title.Render("◆"), styles.Faint.Render(fmt.Sprintf("(%d)", len(tables))))
		a.formatTableList(tables)
		fmt.Println()
	}

	if len(views) > 0 {
		fmt.Printf("%s views %s\n", styles.Title.Render("◆"), styles.Faint.Render(fmt.Sprintf("(%d)", len(views))))
		a.formatTableList(views)
		fmt.Println()
	}
}

func (a *App) formatTableList(items []string) {
	if len(items) == 0 {
		return
	}

	// Filter out any empty or whitespace-only items
	filtered := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			filtered = append(filtered, trimmed)
		}
	}
	items = filtered

	maxLen := 0
	for _, item := range items {
		if len(item) > maxLen {
			maxLen = len(item)
		}
	}

	// Add padding for spacing
	columnWidth := maxLen + 2

	// Simple terminal width detection (default to 120 if we can't detect)
	termWidth := 120

	numColumns := termWidth / columnWidth
	if numColumns < 1 {
		numColumns = 1
	}

	for i, item := range items {
		fmt.Printf("%-*s", columnWidth, item)

		// New line after filling a row
		if (i+1)%numColumns == 0 {
			fmt.Println()
		}
	}

	// New line if last row wasn't complete
	if len(items)%numColumns != 0 {
		fmt.Println()
	}
}

var supportedFileTypes = map[string]bool{
	".csv": true,
	// ".json": true, // Future support
}

func getSupportedExtension(arg string) (string, bool) {
	ext := strings.ToLower(filepath.Ext(arg))
	if supported, ok := supportedFileTypes[ext]; ok && supported {
		return ext, true
	}
	return "", false
}

func parseCSV(path string) (columns []string, data [][]string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	// Read header row
	columns, err = reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("could not read CSV header: %w", err)
	}

	// Read data rows
	data, err = reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("could not read CSV data: %w", err)
	}

	return columns, data, nil
}

func (a *App) handleExploreFile(path string, ext string) error {
	columns, data, err := parseCSV(path)
	if err != nil {
		return err
	}

	// Render table with nil connection (read-only mode)
	columnTypes := make([]string, len(columns))
	_, err = table.Render(
		columns,
		columnTypes,
		data,
		0, // elapsed time
		nil, // no database connection
		"", // no table name
		"", // no primary key
		db.Query{}, // empty query
		a.config.DefaultColumnWidth,
		a.config.UIVisibility,
		a.config.KeyMap, // use config keybindings if available
		nil, // no save callback
	)
	if err != nil {
		return fmt.Errorf("could not render table: %w", err)
	}

	return nil
}
