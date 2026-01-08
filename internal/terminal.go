package internal

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

func TerminalSize() (width, height int, err error) {
	fd := int(os.Stdout.Fd())
	return term.GetSize(fd)
}

func CursorPos() (row, col int, err error) {
	fd := int(os.Stdin.Fd())

	// Save terminal state and set raw mode
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return 0, 0, err
	}
	defer term.Restore(fd, oldState)

	// Query cursor position
	fmt.Print("\x1b[6n")

	// Read response
	reader := bufio.NewReader(os.Stdin)
	_, err = fmt.Fscanf(reader, "\x1b[%d;%dR", &row, &col)
	return row, col, err
}

func ScrollUp(lines int) {
	if lines > 0 {
		fmt.Printf("\x1b[%dS", lines)
	}
}

func ScrollDown(lines int) {
	if lines > 0 {
		fmt.Printf("\x1b[%dT", lines)
	}
}
