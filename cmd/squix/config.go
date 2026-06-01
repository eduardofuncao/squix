package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/eduardofuncao/squix/internal/config"
	"github.com/eduardofuncao/squix/internal/editor"
	"github.com/eduardofuncao/squix/internal/styles"
)

func (a *App) handleConfig() {
	editorCmd, err := editor.CheckEditor()
	if err != nil {
		log.Fatal(err)
	}

	// Open config file in editor
	cmd := exec.Command(editorCmd, config.CfgFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to open editor: %v", err)
	}

	newConfig, err := config.LoadConfig(config.CfgFile)
	if err != nil {
		log.Printf("Warning: Could not reload config: %v", err)
	} else {
		a.config = newConfig
		fmt.Println(styles.Success.Render("✓ Configuration reloaded successfully"))
	}
}
