package main

import (
	"fmt"
	"time"

	"github.com/eduardofuncao/squix/internal/config"
	"github.com/eduardofuncao/squix/internal/db"
	"github.com/eduardofuncao/squix/internal/spinner"
	"github.com/eduardofuncao/squix/internal/styles"
)

func (a *App) printConnStatus(conn db.DatabaseConnection) {
	currConn := a.config.Connections[a.config.CurrentConnection]
	connInfo := fmt.Sprintf("%s/%s", currConn.DBType, currConn.Name)
	if currConn.Schema != "" {
		connInfo += fmt.Sprintf(" (schema: %s)", currConn.Schema)
	}

	queryCount := 0
	if currConn.Queries != nil {
		queryCount = len(currConn.Queries)
	}

	fmt.Printf("Using %s\n", styles.Title.Render(connInfo))

	done := make(chan struct{})
	reachable := make(chan bool)

	go func() {
		err := conn.Ping()
		reachable <- (err == nil)
	}()

	go spinner.CircleWait(done)

	var isReachable bool
	select {
	case isReachable = <-reachable:
		close(done)
	case <-time.After(5 * time.Second):
		close(done)
		isReachable = false
	}

	// Clear spinner lines
	fmt.Print("\r\033[2K")
	fmt.Print("\033[1A")
	fmt.Print("\r\033[2K")

	circleIcon := "●"
	statusText := "reachable"
	if !isReachable {
		circleIcon = "○"
		statusText = "unreachable"
	}

	fmt.Printf("%s Using %s\n", styles.Success.Render(circleIcon), styles.Title.Render(connInfo))
	fmt.Printf("  %d saved queries, %s\n", queryCount, styles.Faint.Render(statusText))
}

func (a *App) handleStatus() {
	if a.config.CurrentConnection == "" {
		fmt.Println(styles.Faint.Render("No active connection"))
		return
	}

	conn := config.FromConnectionYaml(a.config.Connections[a.config.CurrentConnection])
	conn.Open()
	defer conn.Close()

	a.printConnStatus(conn)
}
