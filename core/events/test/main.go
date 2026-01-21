package main

import (
	"context"
	"fmt"
	"os"

	"github.com/AnatoleLucet/loom-term/core/events"
	"golang.org/x/term"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Enable kitty keyboard
	fmt.Print("\x1b[>7u")
	defer fmt.Print("\x1b[<u")

	// Enable mouse tracking
	fmt.Printf("\x1b[?1000h\x1b[?1003h\x1b[?1006h")
	defer fmt.Printf("\x1b[?1000l\x1b[?1003l\x1b[?1006l")

	// Enable bracketed paste mode
	fmt.Printf("\x1b[?2004h")
	defer fmt.Printf("\x1b[?2004l")

	state, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), state)

	listenner := events.NewListener()

	evts, errs := listenner.Listen(ctx)
	fmt.Println("listening for events...")
	for {
		select {
		case event := <-evts:
			fmt.Printf("\r%+v\n", event)
		case err := <-errs:
			cancel()
			println("error:", err.Error())
		}
	}

}
