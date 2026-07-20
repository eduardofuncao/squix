package spinner

import (
	"fmt"
	"os"
	"time"

	isatty "github.com/mattn/go-isatty"
	"github.com/eduardofuncao/squix/internal/styles"
)

// Interactive reports whether stdout is a terminal. When false (e.g. piped to
// an editor or CI), spinner frames and cursor-control escapes would land as
// literal garbage in the captured output, so animation and erasure are skipped.
func Interactive() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

func Wait(done chan struct{}) {
	if !Interactive() {
		<-done
		return
	}
	spinnerStages := []string{"▉", "▊", "▋", "▌", "▍", "▎", "▏", "▎", "▍", "▌", "▋", "▊", "▉"}
	var passed time.Duration = 0
	for {
		for _, s := range spinnerStages {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r%s %.2fs", s, passed.Seconds())
				passed += 100 * time.Millisecond
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func CircleWait(done chan struct{}) {
	if !Interactive() {
		<-done
		return
	}
	// Custom pulsing animation
	stages := []string{" ", ".", "o", "O", "@", "*"}
	for {
		for _, s := range stages {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r%s Checking...", styles.Success.Render(s))
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func CircleWaitWithTimer(done chan struct{}) {
	if !Interactive() {
		<-done
		return
	}
	// Custom pulsing animation with timer
	stages := []string{" ", ".", "o", "O", "@", "*"}
	var passed time.Duration = 0
	for {
		for _, s := range stages {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r%s %.2fs", styles.Success.Render(s), passed.Seconds())
				passed += 100 * time.Millisecond
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
