package main

import (
	"fmt"

	"github.com/eduardofuncao/squix/internal/db"
	"github.com/eduardofuncao/squix/internal/styles"
)

func (a *App) handleDrivers() {
	fmt.Println(styles.Title.Render("Drivers in this build"))
	for _, t := range db.KnownDBTypes() {
		if db.IsDriverBuilt(t) {
			fmt.Println("  " + styles.Success.Render("✓") + " " + t)
		} else {
			msg := "✗ " + t + " — not included (excluded via -tags " + db.DriverBuildTag(t) + ")"
			fmt.Println("  " + styles.Faint.Render(msg))
		}
	}
	fmt.Println()
	fmt.Println(styles.Faint.Render("Rebuild without the tag, or install the full squix build, to enable excluded drivers."))
}
