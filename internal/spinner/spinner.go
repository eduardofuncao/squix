package spinner

import (
	"fmt"
	"os"
	"time"

	"github.com/eduardofuncao/squix/internal/styles"
	isatty "github.com/mattn/go-isatty"
)

// Interactive reports whether stdout is a terminal. When false (e.g. piped to
// an editor or CI), spinner frames and cursor-control escapes would land as
// literal garbage in the captured output, so animation and erasure are skipped.
func Interactive() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

// Stages is the pulsing frame sequence shared by all squix spinners so the
// stdout and TUI variants look identical.
var Stages = []string{" ", ".", "o", "O", "@", "*"}

// TickInterval is how often spinner frames advance.
const TickInterval = 100 * time.Millisecond

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
	for {
		for _, s := range Stages {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r%s Checking...", styles.Success.Render(s))
				time.Sleep(TickInterval)
			}
		}
	}
}

func CircleWaitWithTimer(done chan struct{}) {
	if !Interactive() {
		<-done
		return
	}
	var passed time.Duration = 0
	for {
		for _, s := range Stages {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r%s %.2fs", styles.Success.Render(s), passed.Seconds())
				passed += TickInterval
				time.Sleep(TickInterval)
			}
		}
	}
}
