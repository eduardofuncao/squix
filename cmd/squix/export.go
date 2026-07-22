package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/eduardofuncao/squix/internal/backup"
	"github.com/eduardofuncao/squix/internal/styles"
)

func (a *App) handleExport() {
	var path string
	var doBackup bool
	var format string

	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "--backup":
			doBackup = true
		case "--format":
			if i+1 < len(os.Args) {
				i++
				format = os.Args[i]
			}
		default:
			if path == "" && arg[0] != '-' {
				path = arg
			}
		}
	}

	if !doBackup {
		printError("only --backup is supported currently")
	}

	if a.config.CurrentConnection == "" {
		printError("No active connection. Use 'squix switch <connection>' or 'squix init' first")
	}

	conn := a.config.Connections[a.config.CurrentConnection]
	dbType := conn.DBType
	connString := os.ExpandEnv(conn.ConnString)

	dumper, err := backup.CreateDumper(dbType)
	if err != nil {
		printError("%v", err)
	}

	fileExt := filepath.Ext(path)
	formatName, err := backup.ResolveFormat(fileExt, format, dumper)
	if err != nil {
		printError("%v", err)
	}

	spec := dumper.Formats()[formatName]
	outPath := backup.ResolvePath(path, a.config.CurrentConnection, spec.Ext)

	if err := dumper.Dump(connString, formatName, outPath); err != nil {
		printError("%v", err)
	}

	fmt.Println(styles.Success.Render(fmt.Sprintf("✓ Backup written to %s", outPath)))
}
