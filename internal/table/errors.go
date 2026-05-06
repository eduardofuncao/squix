package table

import (
	"fmt"
	"github.com/eduardofuncao/squix/internal/styles"
	"os"
)

func printError(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(os.Stderr, styles.Error.Render("✗ Error:"), msg)
	os.Exit(1)
}
