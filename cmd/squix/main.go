package main

import (
	"log"

	"github.com/eduardofuncao/squix/internal/config"
	"github.com/eduardofuncao/squix/internal/styles"
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
