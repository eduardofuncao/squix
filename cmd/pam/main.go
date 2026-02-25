package main

import (
	"log"

	"github.com/eduardofuncao/pam/internal/config"
	"github.com/eduardofuncao/pam/internal/styles"
)

func main() {
	cfg, err := config.LoadConfig(config.CfgFile)
	if err != nil {
		log.Fatal("Could not load config file", err)
	}

	// Initialize color scheme
	styles.InitScheme(cfg.ColorScheme, cfg.CustomColorScheme)

	app := NewApp(cfg)
	app.Run()
}
