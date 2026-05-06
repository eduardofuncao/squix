package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/eduardofuncao/squix/internal/db"
	"github.com/eduardofuncao/squix/internal/styles"
)

func (a *App) handleRemove() {
	if len(os.Args) < 3 {
		printError("Usage: squix remove <run-name> [--connection|-c <conn-name>]")
	}

	// Parse flags
	var connectionName string
	args := os.Args[2:]
	for i, arg := range args {
		if arg == "--connection" || arg == "-c" {
			if i+1 < len(args) {
				connectionName = args[i+1]
			}
			break
		}
	}

	// If --connection flag was used, remove connection
	if connectionName != "" {
		a.removeConnection(connectionName)
		return
	}

	// Otherwise, remove query (original behavior)
	conn := a.config.Connections[a.config.CurrentConnection]
	queries := conn.Queries

	query, exists := db.FindQueryWithSelector(queries, os.Args[2])
	if !exists {
		printError("Query '%s' could not be found", os.Args[2])
	}

	delete(conn.Queries, query.Name)

	err := a.config.Save()
	if err != nil {
		printError("Could not save configuration file: %v", err)
	}

	fmt.Println(styles.Success.Render(fmt.Sprintf("✓ Removed run '%s'", query.Name)))
}

func (a *App) removeConnection(connName string) {
	conn, exists := a.config.Connections[connName]
	if !exists {
		printError("Connection '%s' does not exist", connName)
		return
	}

	queryCount := len(conn.Queries)

	if !a.confirmDeletion(connName, queryCount) {
		fmt.Println(styles.Faint.Render("Aborted"))
		return
	}

	if a.config.CurrentConnection == connName {
		a.config.CurrentConnection = ""
	}

	delete(a.config.Connections, connName)

	err := a.config.Save()
	if err != nil {
		printError("Could not save configuration file: %v", err)
		return
	}

	fmt.Println(styles.Success.Render(fmt.Sprintf("✓ Removed connection '%s' and %d queries", connName, queryCount)))
}

func (a *App) confirmDeletion(connName string, queryCount int) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(styles.Error.Render(fmt.Sprintf("This will delete connection '%s' and its %d queries. Continue? [y/N]: ", connName, queryCount)))

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
