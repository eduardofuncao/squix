package main

import (
	"fmt"
	"os"

	"github.com/eduardofuncao/squix/internal/config"
	"github.com/eduardofuncao/squix/internal/styles"
)

func (a *App) handleSwitch() {
	if len(os.Args) < 3 {
		printError("Usage: squix switch/use <db-name>")
	}

	connName := os.Args[2]
	conn, ok := a.config.Connections[connName]
	if !ok {
		printError("Connection '%s' does not exist", connName)
	}
	a.config.CurrentConnection = connName

	err := a.config.Save()
	if err != nil {
		printError("Could not save configuration file: %v", err)
	}

	// The last-query recovery is scoped to a connection; switching resets it.
	if err := config.ClearLastQuery(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not clear last-query: %v\n", err)
	}

	fmt.Println(styles.Success.Render("⇄ Switched to: "), styles.Title.Render(fmt.Sprintf("%s/%s", conn.DBType, connName)))
}
