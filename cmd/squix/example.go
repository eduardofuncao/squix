package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eduardofuncao/squix/internal/db"
	"github.com/eduardofuncao/squix/internal/styles"
)

//go:embed example.sql
var exampleSQL []byte

func (a *App) handleExample() {
	path := "example.db"
	force := false

	for _, arg := range os.Args[2:] {
		switch {
		case arg == "--force":
			force = true
		case !strings.HasPrefix(arg, "-"):
			path = arg
		default:
			printError("unknown flag: %s", arg)
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		printError("resolve path: %v", err)
	}

	if _, err := os.Stat(absPath); err == nil {
		if !force {
			printError("%s already exists (use --force to overwrite)", path)
		}
		if err := os.Remove(absPath); err != nil {
			printError("remove existing file: %v", err)
		}
	} else if !os.IsNotExist(err) {
		printError("stat %s: %v", path, err)
	}

	conn, err := db.CreateConnection("example", "sqlite", absPath)
	if err != nil {
		printError("create connection: %v", err)
	}
	if err := conn.Open(); err != nil {
		printError("open database: %v", err)
	}
	defer conn.Close()

	if _, err := conn.Exec(string(exampleSQL)); err != nil {
		printError("build example database: %v", err)
	}

	fmt.Printf("%s created %s\n", styles.Success.Render("●"), styles.Title.Render(path))
	fmt.Println(
		styles.Faint.Render(
			"an office-themed SQLite database is ready to explore (employees, departments, and timesheets tables)",
		),
	)
	fmt.Println()
	fmt.Println(styles.Faint.Render("try it with:"))
	fmt.Printf(
		"  %s\n",
		styles.Title.Render(
			fmt.Sprintf(
				"squix init --name example --type sqlite --conn %s",
				path,
			),
		),
	)
	fmt.Printf("  %s\n", styles.Title.Render("squix explore"))
	fmt.Printf("  %s\n", styles.Title.Render("squix run \"select * from employees\""))
}
